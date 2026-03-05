"""Abstract LLM moderation provider interface."""

from abc import ABC, abstractmethod

from moderation.analysis.domain import ModerationResult


class ModerationProvider(ABC):
    """Contract for any LLM-based moderation implementation."""

    @abstractmethod
    async def analyze(self, title: str, content: str) -> ModerationResult:
        """Analyze poem content and return a moderation result."""
        ...

    @abstractmethod
    def provider_name(self) -> str:
        """Human-readable name of the LLM provider."""
        ...

    @abstractmethod
    def model_name(self) -> str:
        """Model identifier being used."""
        ...
