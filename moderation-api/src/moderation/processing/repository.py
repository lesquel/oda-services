"""Database repository for moderation — updates poems and creates logs."""

import uuid
from datetime import datetime, timezone

from sqlalchemy import update
from sqlalchemy.ext.asyncio import AsyncSession

from moderation.analysis.domain import ModerationResult
from moderation.shared.models import ModerationLog, Poem


class ModerationRepository:
    """Async SQLAlchemy repository for moderation persistence."""

    def __init__(self, session: AsyncSession) -> None:
        self._session = session

    async def update_poem_moderation(
        self, poem_id: str, result: ModerationResult, *, provider: str, model: str
    ) -> None:
        """Update poem moderation fields and create a log entry."""
        now = datetime.now(timezone.utc)
        status = "rejected" if result.is_flagged else "approved"

        # Update poem
        stmt = (
            update(Poem)
            .where(Poem.id == poem_id)
            .values(
                moderation_status=status,
                moderation_score=result.score,
                moderation_reason=result.reason,
                moderated_at=now,
                moderated_by=f"{provider}/{model}",
                # If approved, set status back to published; if rejected, set to rejected
                status="published" if not result.is_flagged else "rejected",
            )
        )
        await self._session.execute(stmt)

        # Create log entry
        log = ModerationLog(
            id=uuid.uuid4(),
            poem_id=poem_id,
            status=status,
            score=result.score,
            reason=result.reason,
            provider=provider,
            model=model,
            categories=",".join(result.categories),
        )
        self._session.add(log)
        await self._session.commit()
