# Etapa de build com Poetry
FROM python:3.11-slim AS builder

# Instalar Poetry
RUN pip install poetry

# Configurar Poetry para não criar venv (já estamos no container)
ENV POETRY_NO_INTERACTION=1 \
    POETRY_VENV_IN_PROJECT=1 \
    POETRY_CACHE_DIR=/tmp/poetry_cache

# Definir diretório de trabalho
WORKDIR /app

# Copiar arquivos de configuração do Poetry
COPY pyproject.toml poetry.lock ./

# Configurar Poetry para criar venv no projeto
RUN poetry config virtualenvs.in-project true

# Instalar dependências
RUN poetry install --only=main && rm -rf $POETRY_CACHE_DIR

# Etapa de produção
FROM python:3.11-slim AS runtime

# Criar usuário não-root para segurança
RUN groupadd -r appuser && useradd -r -g appuser appuser

# Instalar dependências do sistema necessárias
RUN apt-get update && apt-get install -y \
    && rm -rf /var/lib/apt/lists/*

# Definir diretório de trabalho
WORKDIR /app

# Copiar o ambiente virtual do builder
COPY --from=builder /app/.venv /app/.venv

# Certificar que usamos o venv
ENV PATH="/app/.venv/bin:$PATH"

# Copiar código da aplicação
COPY src ./src

# Criar diretório para modelos se não existir
RUN mkdir -p src/models && chown -R appuser:appuser /app

# Usar usuário não-root
USER appuser

# Expor porta
EXPOSE 8000

# Comando para executar a aplicação
CMD ["uvicorn", "src.main:app", "--host", "0.0.0.0", "--port", "8000"]
