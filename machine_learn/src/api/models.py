"""
Modelos Pydantic para a API de predição de categorias.
"""

from typing import List

from pydantic import BaseModel, Field


class PredictExpenseCategoryRequest(BaseModel):
    """Modelo de requisição para predição de categoria."""

    name: str = Field(..., description="Nome da despesa", example="Farmácia")
    amount_cents: float = Field(..., description="Valor em centavos", example=5000.0)

    class Config:
        json_schema_extra = {"example": {"name": "Farmácia", "amount_cents": 5000.0}}


class PredictExpenseCategoryBatchRequest(BaseModel):
    """Modelo de requisição para predição em lote."""

    expenses: List[PredictExpenseCategoryRequest] = Field(
        ..., description="Lista de despesas"
    )


class PredictionResponse(BaseModel):
    """Modelo de resposta para predição única."""

    category_id: int = Field(..., description="ID da categoria predita")
    name: str = Field(..., description="Nome da despesa")
    amount_cents: float = Field(..., description="Valor em centavos")


class BatchPredictionResponse(BaseModel):
    """Modelo de resposta para predição em lote."""

    predictions: List[PredictionResponse] = Field(..., description="Lista de predições")
