"""Tests for the analysis domain models."""

from moderation.analysis.domain import ModerationResult, PoemPayload


def test_moderation_result_safe():
    result = ModerationResult(is_flagged=False, score=0.1, categories=[], reason="Safe content")
    assert not result.is_flagged
    assert result.score == 0.1
    assert result.categories == []


def test_moderation_result_flagged():
    result = ModerationResult(
        is_flagged=True,
        score=0.9,
        categories=["hate_speech", "violence"],
        reason="Contains hate speech",
    )
    assert result.is_flagged
    assert len(result.categories) == 2


def test_poem_payload():
    payload = PoemPayload(poem_id="abc-123", title="My poem", content="Some content")
    assert payload.poem_id == "abc-123"
    assert payload.title == "My poem"
