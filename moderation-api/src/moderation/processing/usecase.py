"""Moderation use case — orchestrates analysis + persistence."""

import logging

from moderation.analysis.domain import PoemPayload
from moderation.analysis.langchain_provider import LangChainModerationProvider
from moderation.config import settings
from moderation.processing.repository import ModerationRepository
from moderation.shared.database import async_session

logger = logging.getLogger(__name__)


class ModerationUseCase:
    """Process a single poem through the moderation pipeline."""

    def __init__(self) -> None:
        self._provider = LangChainModerationProvider()

    async def moderate(self, payload: PoemPayload) -> None:
        """Analyze a poem and persist the result."""
        logger.info("Moderating poem %s: '%s'", payload.poem_id, payload.title[:50])

        result = await self._provider.analyze(title=payload.title, content=payload.content)

        logger.info(
            "Poem %s — flagged=%s score=%.2f categories=%s",
            payload.poem_id,
            result.is_flagged,
            result.score,
            result.categories,
        )

        # Apply threshold override: if score < threshold, force-approve even if LLM flagged
        if result.score < settings.moderation_threshold and result.is_flagged:
            logger.info(
                "Score %.2f < threshold %.2f — overriding to approved",
                result.score,
                settings.moderation_threshold,
            )
            result.is_flagged = False
            result.reason = f"Auto-approved (score {result.score:.2f} below threshold)"

        async with async_session() as session:
            repo = ModerationRepository(session)
            await repo.update_poem_moderation(
                poem_id=payload.poem_id,
                result=result,
                provider=self._provider.provider_name(),
                model=self._provider.model_name(),
            )

        logger.info("Poem %s moderation complete: %s", payload.poem_id, "rejected" if result.is_flagged else "approved")
