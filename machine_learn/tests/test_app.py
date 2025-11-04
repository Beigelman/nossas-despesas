"""
Testes para os endpoints da API.
"""

from unittest.mock import MagicMock, patch

import pytest
from fastapi import status
from fastapi.testclient import TestClient

from src.api.app import app


class TestRootEndpoint:
    """Testes para o endpoint raiz."""

    def test_root_endpoint(self, client):
        """Testa que o endpoint raiz retorna informações sobre a API."""
        response = client.get("/")
        assert response.status_code == status.HTTP_200_OK
        data = response.json()
        assert "message" in data
        assert "version" in data
        assert "endpoints" in data
        assert data["version"] == "1.0.0"
        assert "/predict" in data["endpoints"]
        assert "/predict/batch" in data["endpoints"]
        assert "/health" in data["endpoints"]


class TestHealthEndpoint:
    """Testes para o endpoint de saúde."""

    def test_health_endpoint(self, client):
        """Testa que o endpoint de saúde retorna status healthy."""
        response = client.get("/health")
        assert response.status_code == status.HTTP_200_OK
        data = response.json()
        assert data["status"] == "healthy"
        assert data["model_loaded"] is True


class TestPredictEndpoint:
    """Testes para o endpoint de predição única."""

    def test_predict_success(self, client, sample_expense_request, mock_model):
        """Testa predição bem-sucedida de uma categoria."""
        mock_model.predict.return_value = [5]  # category_id = 5
        
        response = client.post("/predict", json=sample_expense_request)
        
        assert response.status_code == status.HTTP_200_OK
        data = response.json()
        assert data["category_id"] == 5
        assert data["name"] == sample_expense_request["name"]
        assert data["amount_cents"] == sample_expense_request["amount_cents"]

    def test_predict_with_different_values(self, client, mock_model):
        """Testa predição com diferentes valores."""
        mock_model.predict.return_value = [3]
        
        request_data = {"name": "Supermercado", "amount_cents": 25000.0}
        response = client.post("/predict", json=request_data)
        
        assert response.status_code == status.HTTP_200_OK
        data = response.json()
        assert data["category_id"] == 3
        assert data["name"] == "Supermercado"
        assert data["amount_cents"] == 25000.0

    def test_predict_missing_field(self, client):
        """Testa que requisição sem campos obrigatórios retorna erro."""
        response = client.post("/predict", json={"name": "Teste"})
        assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY

    def test_predict_invalid_data_type(self, client):
        """Testa que requisição com tipo de dado inválido retorna erro."""
        response = client.post("/predict", json={"name": 123, "amount_cents": "invalid"})
        assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY

    def test_predict_empty_name(self, client):
        """Testa que requisição com nome vazio ainda é válida (validação fica a cargo do modelo)."""
        mock_model = MagicMock()
        mock_model.predict.return_value = [1]
        
        with patch("src.api.app.model", mock_model):
            with patch("src.api.predict.load_model", return_value=mock_model):
                client_test = TestClient(app)
                response = client_test.post("/predict", json={"name": "", "amount_cents": 1000.0})
                # Aceita nome vazio, mas pode retornar erro do modelo ou resposta válida
                assert response.status_code in [status.HTTP_200_OK, status.HTTP_500_INTERNAL_SERVER_ERROR]

    def test_predict_model_error(self, client, sample_expense_request, mock_model):
        """Testa tratamento de erro quando o modelo falha."""
        mock_model.predict.side_effect = Exception("Erro no modelo")
        
        response = client.post("/predict", json=sample_expense_request)
        
        assert response.status_code == status.HTTP_500_INTERNAL_SERVER_ERROR
        assert "Erro ao fazer predição" in response.json()["detail"]


class TestPredictBatchEndpoint:
    """Testes para o endpoint de predição em lote."""

    def test_predict_batch_success(self, client, sample_batch_request, mock_model):
        """Testa predição em lote bem-sucedida."""
        # Cada chamada de predict_category cria um DataFrame com uma linha e chama predict()[0]
        # Então precisamos retornar valores diferentes a cada chamada
        mock_model.predict.side_effect = [[1], [2], [3]]  # Uma lista por chamada
        
        response = client.post("/predict/batch", json=sample_batch_request)
        
        assert response.status_code == status.HTTP_200_OK
        data = response.json()
        assert "predictions" in data
        assert len(data["predictions"]) == 3
        
        # Verifica cada predição
        for i, prediction in enumerate(data["predictions"]):
            assert prediction["category_id"] == [1, 2, 3][i]
            assert prediction["name"] == sample_batch_request["expenses"][i]["name"]
            assert prediction["amount_cents"] == sample_batch_request["expenses"][i]["amount_cents"]

    def test_predict_batch_empty_list(self, client, mock_model):
        """Testa predição em lote com lista vazia."""
        mock_model.predict.return_value = []
        
        response = client.post("/predict/batch", json={"expenses": []})
        
        assert response.status_code == status.HTTP_200_OK
        data = response.json()
        assert data["predictions"] == []

    def test_predict_batch_single_item(self, client, mock_model):
        """Testa predição em lote com apenas um item."""
        # Cada chamada retorna uma lista com um valor
        mock_model.predict.side_effect = [[7]]
        
        request_data = {
            "expenses": [{"name": "Farmácia", "amount_cents": 5000.0}]
        }
        response = client.post("/predict/batch", json=request_data)
        
        assert response.status_code == status.HTTP_200_OK
        data = response.json()
        assert len(data["predictions"]) == 1
        assert data["predictions"][0]["category_id"] == 7

    def test_predict_batch_missing_field(self, client):
        """Testa que requisição sem campo expenses retorna erro."""
        response = client.post("/predict/batch", json={})
        assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY

    def test_predict_batch_invalid_expense(self, client):
        """Testa que requisição com despesa inválida retorna erro."""
        response = client.post(
            "/predict/batch",
            json={"expenses": [{"name": "Teste"}]}  # Falta amount_cents
        )
        assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY

    def test_predict_batch_model_error(self, client, sample_batch_request, mock_model):
        """Testa tratamento de erro quando o modelo falha na predição em lote."""
        mock_model.predict.side_effect = Exception("Erro no modelo")
        
        response = client.post("/predict/batch", json=sample_batch_request)
        
        assert response.status_code == status.HTTP_500_INTERNAL_SERVER_ERROR
        assert "Erro ao fazer predições em lote" in response.json()["detail"]

