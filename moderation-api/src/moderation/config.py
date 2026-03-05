"""ODA Moderation API — configuration via pydantic-settings."""

from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    database_url: str = "postgresql+asyncpg://postgres:postgres@localhost:5432/oda"
    nats_url: str = "nats://localhost:4222"

    llm_provider: str = "openai"  # "openai" | "anthropic"
    openai_api_key: str = ""
    anthropic_api_key: str = ""

    moderation_threshold: float = 0.7
    port: int = 8084

    # NATS JetStream config
    nats_stream: str = "POEMS"
    nats_subject: str = "POEMS.moderate"
    nats_consumer: str = "moderation-worker"

    model_config = {"env_file": ".env", "extra": "ignore"}


settings = Settings()
