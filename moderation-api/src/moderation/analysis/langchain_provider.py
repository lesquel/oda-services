"""LangChain LCEL-based moderation provider supporting OpenAI and Anthropic."""

import json
import logging

from langchain_core.language_models import BaseChatModel
from langchain_core.output_parsers import StrOutputParser

from moderation.analysis.domain import ModerationResult
from moderation.analysis.prompt import moderation_prompt
from moderation.analysis.provider import ModerationProvider
from moderation.config import settings

logger = logging.getLogger(__name__)


def _build_llm() -> BaseChatModel:
    """Build the LLM instance based on config."""
    if settings.llm_provider == "anthropic":
        from langchain_anthropic import ChatAnthropic

        return ChatAnthropic(
            model="claude-3-5-haiku-20241022",
            api_key=settings.anthropic_api_key,
            temperature=0,
            max_tokens=512,
        )
    # Default: OpenAI
    from langchain_openai import ChatOpenAI

    return ChatOpenAI(
        model="gpt-4o-mini",
        api_key=settings.openai_api_key,
        temperature=0,
        max_tokens=512,
    )


class LangChainModerationProvider(ModerationProvider):
    """LCEL chain: prompt → LLM → parse JSON → ModerationResult."""

    def __init__(self) -> None:
        self._llm = _build_llm()
        self._chain = moderation_prompt | self._llm | StrOutputParser()

    async def analyze(self, title: str, content: str) -> ModerationResult:
        raw = await self._chain.ainvoke({"title": title, "content": content})

        # Strip markdown fences if present
        text = raw.strip()
        if text.startswith("```"):
            text = text.split("\n", 1)[1] if "\n" in text else text[3:]
            if text.endswith("```"):
                text = text[:-3]
            text = text.strip()

        try:
            data = json.loads(text)
            return ModerationResult(**data)
        except (json.JSONDecodeError, ValueError):
            logger.error("Failed to parse LLM response: %s", text[:200])
            # Fail-safe: flag for manual review
            return ModerationResult(
                is_flagged=True,
                score=1.0,
                categories=["parse_error"],
                reason=f"Could not parse LLM response. Raw: {text[:100]}",
            )

    def provider_name(self) -> str:
        return settings.llm_provider

    def model_name(self) -> str:
        return self._llm.model_name if hasattr(self._llm, "model_name") else str(self._llm.model)
