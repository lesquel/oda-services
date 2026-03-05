FROM python:3.12-slim AS base

WORKDIR /app

# Install uv for fast dependency resolution
RUN pip install --no-cache-dir uv

# Copy dependency files first for caching
COPY pyproject.toml .

# Install dependencies
RUN uv pip install --system -e ".[dev]"

# Copy application code
COPY src/ src/

EXPOSE 8084

CMD ["uvicorn", "moderation.main:app", "--host", "0.0.0.0", "--port", "8084", "--log-level", "info"]
