"""
Serviço FastAPI para predição de categorias de despesas.
"""

from fastapi import FastAPI, HTTPException

from src.api.models import (
    BatchPredictionResponse,
    PredictExpenseCategoryBatchRequest,
    PredictExpenseCategoryRequest,
    PredictionResponse,
)
from src.api.predict import load_model, predict_category

# Carrega o modelo uma vez na inicialização do serviço
model = load_model()

app = FastAPI(
    title="API de Predição de Categorias",
    description="API para predizer a categoria de despesas baseado no nome e valor",
    version="1.0.0",
)


@app.get("/")
async def root():
    """Endpoint raiz com informações sobre a API."""
    return {
        "message": "API de Predição de Categorias de Despesas",
        "version": "1.0.0",
        "endpoints": {
            "/predict": "POST - Predição única de categoria",
            "/predict/batch": "POST - Predição em lote de categorias",
            "/health": "GET - Status de saúde da API",
            "/docs": "GET - Documentação interativa da API",
        },
    }


@app.get("/health")
async def health():
    """Endpoint de saúde da API."""
    return {"status": "healthy", "model_loaded": model is not None}


@app.post("/predict", response_model=PredictionResponse)
async def predict(request: PredictExpenseCategoryRequest):
    """
    Endpoint para predizer a categoria de uma despesa.

    Recebe o nome da despesa e o valor em centavos,
    retorna o category_id predito.
    """
    try:
        category_id = predict_category(
            name=request.name, amount_cents=request.amount_cents, model=model
        )

        return PredictionResponse(
            category_id=int(category_id),
            name=request.name,
            amount_cents=request.amount_cents,
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Erro ao fazer predição: {str(e)}")


@app.post("/predict/batch", response_model=BatchPredictionResponse)
async def predict_batch(request: PredictExpenseCategoryBatchRequest):
    """
    Endpoint para predizer categorias de múltiplas despesas em lote.

    Recebe uma lista de despesas com nome e valor em centavos,
    retorna uma lista de predições com category_id.
    """
    try:
        predictions = []
        for expense in request.expenses:
            category_id = predict_category(
                name=expense.name, amount_cents=expense.amount_cents, model=model
            )
            predictions.append(
                PredictionResponse(
                    category_id=int(category_id),
                    name=expense.name,
                    amount_cents=expense.amount_cents,
                )
            )

        return BatchPredictionResponse(predictions=predictions)
    except Exception as e:
        raise HTTPException(
            status_code=500, detail=f"Erro ao fazer predições em lote: {str(e)}"
        )
