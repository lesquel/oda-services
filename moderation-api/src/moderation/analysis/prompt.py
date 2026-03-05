"""Prompt templates for poem moderation analysis."""

from langchain_core.prompts import ChatPromptTemplate

MODERATION_SYSTEM = """You are a content moderation specialist for a poetry platform.
Analyze the poem below and determine if it violates community guidelines.

Categories of violations:
- hate_speech: Content promoting hatred against protected groups
- violence: Graphic violence, threats, or glorification of harm
- sexual_content: Explicit sexual material
- self_harm: Content promoting self-harm or suicide
- harassment: Targeted harassment or bullying
- illegal_activity: Content promoting illegal activities
- spam: Non-poetry content, advertisements, or spam

IMPORTANT: Poetry often uses metaphor, dark imagery, and strong emotions.
A poem about death is NOT necessarily promoting self-harm.
A poem about war is NOT necessarily promoting violence.
Use artistic context to judge intent.

Respond ONLY with valid JSON matching this exact schema:
{{
  "is_flagged": boolean,
  "score": float (0-1, where 0=completely safe, 1=severe violation),
  "categories": [list of violated category names, empty if safe],
  "reason": "brief explanation"
}}"""

MODERATION_HUMAN = """Title: {title}

Content:
{content}"""

moderation_prompt = ChatPromptTemplate.from_messages([
    ("system", MODERATION_SYSTEM),
    ("human", MODERATION_HUMAN),
])
