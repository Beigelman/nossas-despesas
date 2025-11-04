# Etapa de build com Poetry
FROM python:3.11-slim AS builder

# Instalar Poetry com versão específica para reprodutibilidade
RUN pip install --no-cache-dir poetry==1.7.1

# Configurar Poetry para não criar venv (já estamos no container)
ENV POETRY_NO_INTERACTION=1 \
    POETRY_VENV_IN_PROJECT=1 \
    POETRY_CACHE_DIR=/tmp/poetry_cache

# Definir diretório de trabalho
WORKDIR /app

# Copiar arquivos de configuração do Poetry
COPY machine_learn/pyproject.toml machine_learn/poetry.lock ./

# Instalar dependências
RUN poetry install --only=main && rm -rf $POETRY_CACHE_DIR

# Etapa de produção
FROM python:3.11-slim AS runtime

# Criar usuário não-root para segurança
RUN groupadd -r appuser && useradd -r -g appuser appuser

# Definir diretório de trabalho
WORKDIR /app

# Copiar o ambiente virtual do builder
COPY --from=builder /app/.venv /app/.venv

# Certificar que usamos o venv
ENV PATH="/app/.venv/bin:$PATH" \
    PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1

# Copiar código da aplicação (inclui modelos em src/models)
COPY machine_learn/src ./src

# Criar diretório para modelos e ajustar permissões ANTES de mudar usuário
RUN mkdir -p src/models && chown -R appuser:appuser /app

# Usar usuário não-root
USER appuser

# Expor porta
EXPOSE 8000

# Healthcheck para monitoramento
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:8000/health')" || exit 1

# Comando para executar a aplicação
CMD ["uvicorn", "src.main:app", "--host", "0.0.0.0", "--port", "8000"]
