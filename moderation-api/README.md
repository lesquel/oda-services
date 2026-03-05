# ODA Moderation API

Python service that uses LangChain + LLMs (OpenAI / Anthropic) to moderate poem content via NATS JetStream.

## Architecture

```
NATS ──► consumer.py ──► usecase.py ──► provider.py (LangChain) ──► repository.py ──► PostgreSQL
```

### Features

- **analysis/**: LLM analysis domain, prompt templates, LangChain provider
- **processing/**: NATS consumer, moderation use case, DB repository
- **health/**: FastAPI health endpoint
- **shared/**: Database engine, SQLAlchemy models, config

## Local development

```bash
uv pip install -e ".[dev]"
uvicorn moderation.main:app --reload --port 8084
```

## Environment variables

| Variable               | Default                  | Description                  |
|------------------------|--------------------------|------------------------------|
| `DATABASE_URL`         | `postgresql+asyncpg://…` | Async PostgreSQL DSN         |
| `NATS_URL`             | `nats://localhost:4222`  | NATS server URL              |
| `LLM_PROVIDER`         | `openai`                 | `openai` or `anthropic`      |
| `OPENAI_API_KEY`       | —                        | OpenAI API key               |
| `ANTHROPIC_API_KEY`    | —                        | Anthropic API key            |
| `MODERATION_THRESHOLD` | `0.7`                    | Score threshold for flagging |
| `PORT`                 | `8084`                   | HTTP port                    |
