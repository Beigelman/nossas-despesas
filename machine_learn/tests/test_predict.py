"""
Testes para as funções de predição.
"""

from unittest.mock import MagicMock, patch

import pandas as pd
import pytest
from src.api.predict import predict_batch, predict_category


class TestLoadModel:
    """Testes para a função load_model."""

    def test_load_model_success(self):
        """Testa carregamento bem-sucedido do modelo."""
        # Para este teste, precisamos parar temporariamente o patch global
        from tests.conftest import _patcher

        _patcher.stop()
        try:
            with patch("src.api.predict.joblib.load") as mock_joblib_load:
                with patch("src.api.predict.Path.exists") as mock_exists:
                    mock_exists.return_value = True
                    mock_model = MagicMock()
                    mock_joblib_load.return_value = mock_model

                    # Importa a função diretamente para evitar cache
                    import importlib

                    import src.api.predict as predict_module

                    importlib.reload(predict_module)
                    from src.api.predict import load_model as original_load_model

                    model = original_load_model()

                    assert model == mock_model
                    mock_exists.assert_called_once()
                    mock_joblib_load.assert_called_once()
        finally:
            _patcher.start()

    def test_load_model_not_found(self):
        """Testa erro quando modelo não é encontrado."""
        from tests.conftest import _patcher

        _patcher.stop()
        try:
            with patch("src.api.predict.Path.exists") as mock_exists:
                mock_exists.return_value = False

                # Importa a função diretamente para evitar cache
                import importlib

                import src.api.predict as predict_module

                importlib.reload(predict_module)
                from src.api.predict import load_model as original_load_model

                with pytest.raises(FileNotFoundError) as exc_info:
                    original_load_model()

                assert "Modelo não encontrado" in str(exc_info.value)
        finally:
            _patcher.start()

    def test_load_model_custom_path(self):
        """Testa carregamento com caminho customizado."""
        from tests.conftest import _patcher

        _patcher.stop()
        try:
            with patch("src.api.predict.joblib.load") as mock_joblib_load:
                with patch("src.api.predict.Path.exists") as mock_exists:
                    mock_exists.return_value = True
                    mock_model = MagicMock()
                    mock_joblib_load.return_value = mock_model

                    # Importa a função diretamente para evitar cache
                    import importlib

                    import src.api.predict as predict_module

                    importlib.reload(predict_module)
                    from src.api.predict import load_model as original_load_model

                    model = original_load_model(model_path="custom/path/model.pkl")

                    assert model == mock_model
        finally:
            _patcher.start()


class TestPredictCategory:
    """Testes para a função predict_category."""

    def test_predict_category_with_model(self):
        """Testa predição com modelo fornecido."""
        mock_model = MagicMock()
        mock_model.predict.return_value = [5]

        result = predict_category(
            name="Farmácia", amount_cents=5000.0, model=mock_model
        )

        assert result == 5
        mock_model.predict.assert_called_once()

        # Verifica que o DataFrame passado ao modelo está correto
        call_args = mock_model.predict.call_args[0][0]
        assert isinstance(call_args, pd.DataFrame)
        assert call_args.iloc[0]["name"] == "Farmácia"
        assert call_args.iloc[0]["amount_cents"] == 5000.0

    def test_predict_category_with_int_amount(self):
        """Testa predição com amount_cents como int."""
        mock_model = MagicMock()
        mock_model.predict.return_value = [3]

        result = predict_category(
            name="Supermercado", amount_cents=15000, model=mock_model
        )

        assert result == 3

    def test_predict_category_without_model(self):
        """Testa predição sem modelo fornecido (carrega automaticamente)."""
        mock_model = MagicMock()
        mock_model.predict.return_value = [7]

        with patch("src.api.predict.load_model", return_value=mock_model):
            result = predict_category(name="Restaurante", amount_cents=8000.0)

        assert result == 7

    def test_predict_category_empty_name(self):
        """Testa predição com nome vazio."""
        mock_model = MagicMock()
        mock_model.predict.return_value = [1]

        result = predict_category(name="", amount_cents=1000.0, model=mock_model)

        assert result == 1
        call_args = mock_model.predict.call_args[0][0]
        assert call_args.iloc[0]["name"] == ""

    def test_predict_category_zero_amount(self):
        """Testa predição com valor zero."""
        mock_model = MagicMock()
        mock_model.predict.return_value = [2]

        result = predict_category(name="Teste", amount_cents=0.0, model=mock_model)

        assert result == 2
        call_args = mock_model.predict.call_args[0][0]
        assert call_args.iloc[0]["amount_cents"] == 0.0

    def test_predict_category_type_conversion(self):
        """Testa que converte tipos corretamente."""
        mock_model = MagicMock()
        mock_model.predict.return_value = [4]

        # Passa name como número e amount_cents como string
        result = predict_category(name=123, amount_cents="5000", model=mock_model)

        assert result == 4
        call_args = mock_model.predict.call_args[0][0]
        assert call_args.iloc[0]["name"] == "123"  # Convertido para string
        assert call_args.iloc[0]["amount_cents"] == 5000.0  # Convertido para float


class TestPredictBatch:
    """Testes para a função predict_batch."""

    def test_predict_batch_with_model(self):
        """Testa predição em lote com modelo fornecido."""
        mock_model = MagicMock()
        mock_model.predict.return_value = [1, 2, 3]

        df = pd.DataFrame(
            {
                "name": ["Farmácia", "Supermercado", "Restaurante"],
                "amount_cents": [5000.0, 15000.0, 8000.0],
            }
        )

        result = predict_batch(df, model=mock_model)

        assert "predicted_category_id" in result.columns
        assert list(result["predicted_category_id"]) == [1, 2, 3]
        assert len(result) == 3
        mock_model.predict.assert_called_once()

    def test_predict_batch_without_model(self):
        """Testa predição em lote sem modelo fornecido."""
        mock_model = MagicMock()
        mock_model.predict.return_value = [5, 6]

        df = pd.DataFrame(
            {
                "name": ["Teste1", "Teste2"],
                "amount_cents": [1000.0, 2000.0],
            }
        )

        with patch("src.api.predict.load_model", return_value=mock_model):
            result = predict_batch(df)

        assert "predicted_category_id" in result.columns
        assert list(result["predicted_category_id"]) == [5, 6]

    def test_predict_batch_empty_dataframe(self):
        """Testa predição em lote com DataFrame vazio."""
        mock_model = MagicMock()
        mock_model.predict.return_value = []

        df = pd.DataFrame(columns=["name", "amount_cents"])

        result = predict_batch(df, model=mock_model)

        assert "predicted_category_id" in result.columns
        assert len(result) == 0

    def test_predict_batch_single_row(self):
        """Testa predição em lote com apenas uma linha."""
        mock_model = MagicMock()
        mock_model.predict.return_value = [7]

        df = pd.DataFrame(
            {
                "name": ["Único"],
                "amount_cents": [3000.0],
            }
        )

        result = predict_batch(df, model=mock_model)

        assert len(result) == 1
        assert result.iloc[0]["predicted_category_id"] == 7

    def test_predict_batch_does_not_modify_original(self):
        """Testa que a função não modifica o DataFrame original."""
        mock_model = MagicMock()
        mock_model.predict.return_value = [1]

        df = pd.DataFrame(
            {
                "name": ["Teste"],
                "amount_cents": [1000.0],
            }
        )
        result = predict_batch(df, model=mock_model)

        # DataFrame original não deve ter a coluna predicted_category_id
        assert "predicted_category_id" not in df.columns
        assert "predicted_category_id" in result.columns
        # Verifica que o resultado é uma cópia
        assert result is not df
