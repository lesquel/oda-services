"""FastAPI application entrypoint with NATS consumer lifecycle."""

import asyncio
import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI

from moderation.config import settings
from moderation.health.router import router as health_router
from moderation.processing.consumer import NATSConsumer
from moderation.shared.database import engine

logger = logging.getLogger("moderation")
logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(name)s: %(message)s")


@asynccontextmanager
async def lifespan(_app: FastAPI):
    """Start NATS consumer on startup, gracefully stop on shutdown."""
    consumer = NATSConsumer()
    task = asyncio.create_task(consumer.start())
    logger.info("Moderation consumer started")
    yield
    await consumer.stop()
    task.cancel()
    await engine.dispose()
    logger.info("Moderation service shut down cleanly")


app = FastAPI(
    title="ODA Moderation API",
    version="0.1.0",
    lifespan=lifespan,
)

app.include_router(health_router)

if __name__ == "__main__":
    import uvicorn

    uvicorn.run("moderation.main:app", host="0.0.0.0", port=settings.port, reload=True)
