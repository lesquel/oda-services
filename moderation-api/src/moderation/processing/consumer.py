"""NATS JetStream consumer — subscribes to POEMS.moderate and dispatches to use case."""

import asyncio
import json
import logging

import nats
from nats.js.api import ConsumerConfig, DeliverPolicy, StreamConfig

from moderation.analysis.domain import PoemPayload
from moderation.config import settings
from moderation.processing.usecase import ModerationUseCase

logger = logging.getLogger(__name__)


class NATSConsumer:
    """Long-running NATS JetStream pull consumer."""

    def __init__(self) -> None:
        self._nc: nats.NATS | None = None
        self._sub = None
        self._running = False
        self._usecase = ModerationUseCase()

    async def start(self) -> None:
        """Connect to NATS and start consuming messages."""
        self._running = True

        # Retry connection with backoff
        for attempt in range(30):
            try:
                self._nc = await nats.connect(settings.nats_url)
                logger.info("Connected to NATS at %s", settings.nats_url)
                break
            except Exception:
                wait = min(2**attempt, 30)
                logger.warning("NATS connection attempt %d failed, retrying in %ds…", attempt + 1, wait)
                await asyncio.sleep(wait)
        else:
            logger.error("Failed to connect to NATS after 30 attempts")
            return

        js = self._nc.jetstream()

        # Ensure stream exists
        try:
            await js.find_stream_info_by_subject(settings.nats_subject)
        except nats.js.errors.NotFoundError:
            await js.add_stream(
                StreamConfig(
                    name=settings.nats_stream,
                    subjects=[f"{settings.nats_stream}.*"],
                    retention="workqueue",
                    max_age=86400 * 1_000_000_000,  # 24h in nanoseconds
                )
            )
            logger.info("Created JetStream stream: %s", settings.nats_stream)

        # Create durable push subscription
        self._sub = await js.subscribe(
            settings.nats_subject,
            durable=settings.nats_consumer,
            config=ConsumerConfig(
                deliver_policy=DeliverPolicy.ALL,
                ack_wait=120,  # 2min for LLM calls
            ),
        )
        logger.info("Subscribed to %s (durable: %s)", settings.nats_subject, settings.nats_consumer)

        # Message loop
        while self._running:
            try:
                msgs = await self._sub.fetch(batch=1, timeout=5)
                for msg in msgs:
                    await self._handle_message(msg)
            except nats.errors.TimeoutError:
                continue
            except Exception:
                logger.exception("Error in consumer loop")
                await asyncio.sleep(1)

    async def _handle_message(self, msg) -> None:
        """Process a single NATS message."""
        try:
            data = json.loads(msg.data.decode())
            payload = PoemPayload(**data)
            await self._usecase.moderate(payload)
            await msg.ack()
        except Exception:
            logger.exception("Failed to process message: %s", msg.data[:200])
            # NAK with delay for retry
            await msg.nak(delay=10)

    async def stop(self) -> None:
        """Gracefully stop the consumer."""
        self._running = False
        if self._sub:
            await self._sub.unsubscribe()
        if self._nc and self._nc.is_connected:
            await self._nc.drain()
            logger.info("NATS consumer stopped")
