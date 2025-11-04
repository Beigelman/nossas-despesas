"""
Configurações e fixtures compartilhadas para os testes.
"""

from unittest.mock import MagicMock, patch

import pytest
from fastapi.testclient import TestClient

# Faz o patch do load_model antes de qualquer importação do app
# Isso evita que o modelo real seja carregado durante os testes
_mock_model = MagicMock()
_mock_model.predict.return_value = [1]
_patcher = patch("src.api.predict.load_model", return_value=_mock_model)
_patcher.start()


@pytest.fixture(scope="session", autouse=True)
def mock_load_model():
    """Fixture que mantém o patch do load_model durante toda a sessão de testes."""
    yield _mock_model


@pytest.fixture
def mock_model(mock_load_model):
    """Fixture que retorna um modelo mockado."""
    return mock_load_model


@pytest.fixture
def client(mock_model):
    """Fixture que retorna um cliente de teste do FastAPI com modelo mockado."""
    # Importa o módulo app - como load_model já está com patch, o modelo será mockado
    import importlib
    import src.api.app as app_module
    
    # Recarrega o módulo para garantir que usa o modelo mockado
    importlib.reload(app_module)
    
    from src.api.app import app
    
    # Garante que o modelo no módulo é o mockado
    app_module.model = mock_model
    
    # Cria o cliente de teste
    test_client = TestClient(app)
    
    yield test_client


@pytest.fixture
def sample_expense_request():
    """Fixture com dados de exemplo para uma despesa."""
    return {
        "name": "Farmácia",
        "amount_cents": 5000.0,
    }


@pytest.fixture
def sample_batch_request():
    """Fixture com dados de exemplo para predição em lote."""
    return {
        "expenses": [
            {"name": "Farmácia", "amount_cents": 5000.0},
            {"name": "Supermercado", "amount_cents": 15000.0},
            {"name": "Restaurante", "amount_cents": 8000.0},
        ]
    }

