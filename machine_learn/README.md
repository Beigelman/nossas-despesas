# ML Service - Nossas Despesas

The Machine Learning service provides automatic expense category classification using trained models. It's built with FastAPI and Python, offering both real-time and batch prediction capabilities.

## Overview

This service predicts the `category_id` of an expense based on:
- **Expense name** (text): The description or name of the expense
- **Amount in cents** (numeric): The expense value in cents

The service uses a trained machine learning model (currently SVM with RBF kernel) to classify expenses into predefined categories.

## Technology Stack

- **Framework**: FastAPI
- **Language**: Python 3.9+
- **ML Libraries**: scikit-learn, pandas, numpy
- **Model Persistence**: joblib
- **Optional**: XGBoost, LightGBM (for improved performance)
- **Package Management**: Poetry
- **Testing**: pytest, pytest-asyncio, httpx

## Project Structure

```
machine_learn/
├── src/
│   ├── api/
│   │   ├── app.py              # FastAPI application
│   │   ├── models.py            # Pydantic request/response models
│   │   └── predict.py           # Prediction functions
│   ├── models/                  # Trained model files
│   │   └── svm_rbf.pkl         # Current production model
│   └── main.py                  # Application entry point
├── training/
│   ├── train_model.py          # Model training script
│   └── data/                    # Training data (if available)
│       └── training_data.csv
├── tests/
│   ├── conftest.py             # Pytest configuration and fixtures
│   ├── test_app.py             # API endpoint tests
│   └── test_predict.py          # Prediction function tests
├── pyproject.toml               # Poetry dependencies
├── pytest.ini                   # Pytest configuration
└── Dockerfile                   # Container configuration
```

## Getting Started

### Prerequisites

- Python 3.9 or higher
- Poetry (for dependency management)

### Installation

1. **Install Poetry** (if not already installed):
   ```bash
   curl -sSL https://install.python-poetry.org | python3 -
   ```

2. **Install dependencies**:
   ```bash
   cd machine_learn
   
   # Basic installation (without XGBoost and LightGBM)
   poetry install
   
   # Full installation (with XGBoost and LightGBM for better performance)
   poetry install --extras full
   ```

3. **Activate Poetry shell** (optional):
   ```bash
   poetry shell
   ```

### Running the Service

Start the FastAPI server:

```bash
# With Poetry shell activated
uvicorn src.main:app --host 0.0.0.0 --port 8000 --reload

# Or directly with Poetry
poetry run uvicorn src.main:app --host 0.0.0.0 --port 8000 --reload
```

The API will be available at:
- **API**: `http://localhost:8000`
- **Interactive Docs**: `http://localhost:8000/docs`
- **Alternative Docs**: `http://localhost:8000/redoc`

## API Endpoints

### `GET /`
Root endpoint with API information and available endpoints.

**Response**:
```json
{
  "message": "API de Predição de Categorias de Despesas",
  "version": "1.0.0",
  "endpoints": {
    "/predict": "POST - Predição única de categoria",
    "/predict/batch": "POST - Predição em lote de categorias",
    "/health": "GET - Status de saúde da API",
    "/docs": "GET - Documentação interativa da API"
  }
}
```

### `GET /health`
Health check endpoint.

**Response**:
```json
{
  "status": "healthy",
  "model_loaded": true
}
```

### `POST /predict`
Predict category for a single expense.

**Request**:
```json
{
  "name": "Farmácia",
  "amount_cents": 5000
}
```

**Response**:
```json
{
  "category_id": 5,
  "name": "Farmácia",
  "amount_cents": 5000
}
```

### `POST /predict/batch`
Predict categories for multiple expenses in batch.

**Request**:
```json
{
  "expenses": [
    {"name": "Farmácia", "amount_cents": 5000},
    {"name": "Gasolina", "amount_cents": 25000},
    {"name": "Aluguel", "amount_cents": 350000}
  ]
}
```

**Response**:
```json
{
  "predictions": [
    {"category_id": 5, "name": "Farmácia", "amount_cents": 5000},
    {"category_id": 3, "name": "Gasolina", "amount_cents": 25000},
    {"category_id": 1, "name": "Aluguel", "amount_cents": 350000}
  ]
}
```

## Model Training

### Training Pipeline

The training script (`training/train_model.py`) implements a comprehensive ML pipeline:

1. **Data Loading**: Loads training data from CSV
2. **Preprocessing**:
   - Removes null values
   - Normalizes data types
   - Handles class imbalance
3. **Feature Engineering**:
   - **Text (name)**: TF-IDF vectorization with n-grams (1-2), max 500 features
   - **Numeric (amount_cents)**: StandardScaler normalization
4. **Model Training**: Tests multiple algorithms:
   - Random Forest
   - Gradient Boosting
   - Logistic Regression
   - SVM (RBF) - Currently best performing
   - K-Nearest Neighbors
   - XGBoost (if available)
   - LightGBM (if available)
5. **Evaluation**: Compares models using accuracy, F1-score (macro and weighted)
6. **Model Persistence**: Saves the best model to `src/models/`

### Running Training

1. **Prepare training data**:
   - Place CSV file at `training/data/training_data.csv`
   - Required columns: `name`, `amount_cents`, `category_id`

2. **Run training**:
   ```bash
   cd training
   poetry run python train_model.py
   ```

3. **Model output**:
   - Best model saved to `src/models/svm_rbf.pkl` (or best performing)
   - Default model saved to `src/models/category_classifier.pkl`

### Training Data Format

The training CSV should have the following structure:

```csv
name,amount_cents,category_id
Farmácia,5000,5
Gasolina,25000,3
Aluguel,350000,1
...
```

## Model Architecture

### Current Production Model

**SVM (RBF Kernel)**:
- **Text Processing**: TF-IDF with 1-2 n-grams, 500 max features
- **Numeric Processing**: StandardScaler
- **Classifier**: Support Vector Machine with RBF kernel
- **Probability**: Enabled for confidence scores

### Feature Engineering

1. **Text Features (name)**:
   - TF-IDF vectorization
   - N-gram range: (1, 2)
   - Maximum features: 500
   - Lowercase conversion
   - Accent stripping (Unicode)

2. **Numeric Features (amount_cents)**:
   - StandardScaler normalization
   - Handles various expense amounts

### Model Selection

The training script automatically:
- Tests multiple algorithms
- Compares performance metrics
- Selects the best model based on accuracy and F1-score
- Saves the best model for production use

## Testing

### Running Tests

```bash
# Run all tests
poetry run pytest

# Run with verbose output
poetry run pytest -v

# Run with code coverage
poetry run pytest --cov=src --cov-report=html

# Run specific test file
poetry run pytest tests/test_app.py

# Run specific test class
poetry run pytest tests/test_app.py::TestPredictEndpoint
```

### Test Structure

- **`tests/test_app.py`**: Tests for API endpoints (GET /, GET /health, POST /predict, POST /predict/batch)
- **`tests/test_predict.py`**: Tests for prediction functions (load_model, predict_category, predict_batch)
- **`tests/conftest.py`**: Shared fixtures and test configuration

Tests use mocks to avoid requiring a trained model during test execution.

## Development

### Code Organization

- **`src/api/app.py`**: FastAPI application setup and route definitions
- **`src/api/models.py`**: Pydantic models for request/response validation
- **`src/api/predict.py`**: Core prediction logic and model loading
- **`training/train_model.py`**: Model training pipeline

### Adding New Models

To add a new model to the training pipeline:

1. Import the model class in `training/train_model.py`
2. Add it to the `create_models()` function
3. The training script will automatically test it and compare performance

### Model Versioning

Models are saved with descriptive names based on the algorithm:
- `svm_rbf.pkl` - SVM with RBF kernel
- `random_forest.pkl` - Random Forest
- `category_classifier.pkl` - Default/best model

## Deployment

The service is deployed to **Google Cloud Run** as a containerized application.

### Docker Build

```bash
docker build -t nossas-despesas-ml .
docker run -p 8000:8000 nossas-despesas-ml
```

### Environment Variables

The service doesn't require environment variables for basic operation. The model path is hardcoded to `src/models/svm_rbf.pkl`.

### Production Considerations

- **Model Loading**: Model is loaded once at startup for performance
- **Error Handling**: All endpoints include proper error handling
- **Health Checks**: `/health` endpoint for monitoring
- **Logging**: FastAPI's built-in logging for request tracking

## Future Improvements

### Model Enhancements
- Hyperparameter tuning with GridSearchCV/RandomizedSearchCV
- Cross-validation for better evaluation
- Handling class imbalance with SMOTE or class weights
- Deep learning models (LSTM/Transformer) for text processing

### Feature Engineering
- Extract specific keywords from expense names
- Create binary features for common values
- More sophisticated text normalization
- Feature selection to reduce dimensionality

### Infrastructure
- Model versioning system
- A/B testing for model comparison
- Model retraining pipeline
- Monitoring and metrics collection

## Contributing

1. Follow Python best practices (PEP 8)
2. Write tests for new features
3. Update documentation as needed
4. Run `poetry run pytest` before committing
5. Ensure code passes linting (flake8, mypy)
