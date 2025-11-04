"""
Script para fazer predições usando o modelo treinado.
"""

from pathlib import Path

import joblib
import pandas as pd


def load_model(model_path="../models/svm_rbf.pkl"):
    """Carrega o modelo treinado."""
    model_path = Path(__file__).parent / model_path
    if not model_path.exists():
        raise FileNotFoundError(
            f"Modelo não encontrado em {model_path}. Execute train_model.py primeiro."
        )

    return joblib.load(model_path)


def predict_category(name, amount_cents, model=None):
    """
    Prediz a categoria de uma despesa.

    Args:
        name: Nome da despesa (string)
        amount_cents: Valor em centavos (int ou float)
        model: Modelo carregado (opcional, será carregado automaticamente se None)

    Returns:
        category_id: ID da categoria predita
    """
    if model is None:
        model = load_model()

    # Prepara dados
    data = pd.DataFrame([{"name": str(name), "amount_cents": float(amount_cents)}])

    # Faz predição
    prediction = model.predict(data)[0]

    return prediction


def predict_batch(df, model=None):
    """
    Prediz categorias para múltiplas despesas.

    Args:
        df: DataFrame com colunas 'name' e 'amount_cents'
        model: Modelo carregado (opcional)

    Returns:
        DataFrame com coluna 'predicted_category_id' adicionada
    """
    if model is None:
        model = load_model()

    # Faz predições
    predictions = model.predict(df[["name", "amount_cents"]])

    # Adiciona predições ao DataFrame
    df = df.copy()
    df["predicted_category_id"] = predictions

    return df
