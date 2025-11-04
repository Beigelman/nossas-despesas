# Pipeline de Treinamento - Classificação de Categorias de Despesas

Este projeto contém um pipeline de machine learning para prever a categoria (`category_id`) de uma despesa baseado no nome da despesa e no valor em centavos.

## Estrutura

```
training/
├── data/
│   └── training_data.csv    # Dados de treinamento
├── models/                   # Modelos treinados (criado após treinamento)
├── train_model.py          # Script de treinamento
├── predict.py              # Script para fazer predições
├── pyproject.toml          # Configuração do Poetry
├── poetry.lock             # Lock file do Poetry (gerado automaticamente)
├── requirements.txt        # Dependências Python (legado)
└── README.md              # Este arquivo
```

## Instalação

Este projeto usa [Poetry](https://python-poetry.org/) para gerenciamento de dependências.

1. Instale o Poetry (se ainda não tiver instalado):

```bash
curl -sSL https://install.python-poetry.org | python3 -
```

Ou no Windows (PowerShell):
```powershell
(Invoke-WebRequest -Uri https://install.python-poetry.org -UseBasicParsing).Content | python -
```

2. Instale as dependências:

```bash
# Instalação básica (sem XGBoost e LightGBM)
poetry install

# Instalação completa (com XGBoost e LightGBM para melhor performance)
poetry install --extras full
```

3. Ative o ambiente virtual do Poetry:

```bash
poetry shell
```

Ou execute comandos diretamente com Poetry:

```bash
poetry run python train_model.py
```

## Uso

### Treinamento do Modelo

Execute o script de treinamento:

```bash
# Com Poetry shell ativado
python train_model.py

# Ou usando Poetry diretamente
poetry run python train_model.py
```

O script irá:
1. Carregar os dados de `data/training_data.csv`
2. Preprocessar os dados (limpar valores nulos, normalizar tipos)
3. Dividir os dados em conjunto de treino (80%) e validação (20%)
4. Criar um pipeline que:
   - Processa o texto do nome usando TF-IDF
   - Normaliza o valor em centavos
   - Treina um classificador Random Forest
5. Avaliar o modelo e mostrar métricas
6. Salvar o modelo treinado em `models/category_classifier.pkl`

### Fazer Predições

#### Via linha de comando:

```bash
# Com Poetry shell ativado
python predict.py "Farmácia" 5000

# Ou usando Poetry diretamente
poetry run python predict.py "Farmácia" 5000
```

#### Via código Python:

```python
from predict import predict_category

# Predição única
category_id = predict_category("Farmácia", 5000)
print(f"Categoria predita: {category_id}")

# Predição em lote
import pandas as pd
from predict import predict_batch

df = pd.DataFrame({
    'name': ['Farmácia', 'Gasolina', 'Aluguel'],
    'amount_cents': [5000, 25000, 350000]
})

df_with_predictions = predict_batch(df)
print(df_with_predictions)
```

## Testes

O projeto inclui uma suíte completa de testes para a API. Os testes são executados usando pytest.

### Executando os Testes

```bash
# Executar todos os testes
poetry run pytest

# Executar com verbose (mostra mais detalhes)
poetry run pytest -v

# Executar com cobertura de código
poetry run pytest --cov=src --cov-report=html

# Executar um arquivo de teste específico
poetry run pytest tests/test_app.py

# Executar uma classe de teste específica
poetry run pytest tests/test_app.py::TestPredictEndpoint
```

### Estrutura dos Testes

Os testes estão organizados em:

- `tests/test_app.py` - Testes para os endpoints da API (GET /, GET /health, POST /predict, POST /predict/batch)
- `tests/test_models.py` - Testes para validação dos modelos Pydantic (request/response)
- `tests/test_predict.py` - Testes para as funções de predição (load_model, predict_category, predict_batch)
- `tests/conftest.py` - Fixtures compartilhadas e configurações de teste

Os testes usam mocks para evitar a necessidade de ter o modelo treinado disponível durante a execução dos testes.

## Pipeline de Machine Learning

O pipeline implementado:

1. **Feature Engineering:**
   - **Texto (name)**: Usa TF-IDF com n-grams (1-2), máximo de 500 features
   - **Numérico (amount_cents)**: Normalização usando StandardScaler

2. **Classificador:**
   - Random Forest com 200 árvores
   - Profundidade máxima de 20
   - Otimizado para classificação multiclasse

## Métricas de Avaliação

O script de treinamento mostra:
- Acurácia no conjunto de validação
- Relatório de classificação completo (precision, recall, F1-score)
- Matriz de confusão

## Melhorias Futuras

Possíveis melhorias que podem ser implementadas:

1. **Feature Engineering:**
   - Extrair palavras-chave específicas do nome
   - Criar features binárias para valores comuns
   - Normalização mais sofisticada do texto (remover acentos, normalizar variações)

2. **Modelos Alternativos:**
   - XGBoost ou LightGBM para melhor performance
   - Ensemble de múltiplos modelos
   - Modelos de deep learning (LSTM/Transformer) para processamento de texto

3. **Otimização:**
   - Busca de hiperparâmetros usando GridSearchCV ou RandomizedSearchCV
   - Validação cruzada para melhor avaliação
   - Tratamento de classes desbalanceadas

4. **Deploy:**
   - API REST usando Flask/FastAPI
   - Integração com cloud functions
   - Versionamento de modelos

