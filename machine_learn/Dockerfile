# Build stage with Poetry
FROM python:3.11-slim AS builder

# Install Poetry with a specific version for reproducibility
# Using a version that supports package-mode (1.5.0+)
RUN pip install --no-cache-dir poetry==2.2.1

# Configure Poetry to not create a venv (we are already in the container)
ENV POETRY_NO_INTERACTION=1 \
    POETRY_VENV_IN_PROJECT=1 \
    POETRY_CACHE_DIR=/tmp/poetry_cache

# Set working directory
WORKDIR /app

# Copy Poetry configuration files
COPY machine_learn/pyproject.toml machine_learn/poetry.lock ./

# Configure Poetry to create venv in project
RUN poetry config virtualenvs.in-project true && \
    poetry config virtualenvs.create true

# Install dependencies
RUN poetry install --only=main --no-root && rm -rf $POETRY_CACHE_DIR

# Production stage
FROM python:3.11-slim AS runtime

# Create non-root user for security
RUN groupadd -r appuser && useradd -r -g appuser appuser

# Set working directory
WORKDIR /app

# Copy the virtual environment from the builder
COPY --from=builder /app/.venv /app/.venv

# Ensure we use the venv
ENV PATH="/app/.venv/bin:$PATH" \
    PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1

# Copy application code (includes models in src/models)
COPY machine_learn/src ./src

# Create directory for models and adjust permissions BEFORE switching user
RUN mkdir -p src/models && chown -R appuser:appuser /app

# Use non-root user
USER appuser

# Expose port
EXPOSE 8000

# Command to run the application
CMD ["uvicorn", "src.main:app", "--host", "0.0.0.0", "--port", "8000"]
