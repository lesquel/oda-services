"""Analysis domain models."""

from pydantic import BaseModel, Field


class ModerationResult(BaseModel):
    """Structured output from the LLM moderation analysis."""

    is_flagged: bool = Field(description="Whether the content should be flagged for review")
    score: float = Field(ge=0.0, le=1.0, description="Severity score from 0 (safe) to 1 (severe)")
    categories: list[str] = Field(
        default_factory=list,
        description="List of violated category names",
    )
    reason: str = Field(description="Human-readable explanation of the moderation decision")


class PoemPayload(BaseModel):
    """Message payload received from NATS for moderation."""

    poem_id: str
    title: str
    content: str
