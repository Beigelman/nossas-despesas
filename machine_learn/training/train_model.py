"""
Pipeline de treinamento para predição de category_id de despesas.
Usa o nome da despesa e o valor em centavos para prever a categoria.
"""

import os
import time
from pathlib import Path

import joblib
import numpy as np
import pandas as pd
from sklearn.compose import ColumnTransformer
from sklearn.ensemble import (
    GradientBoostingClassifier,
    RandomForestClassifier,
)
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.linear_model import LogisticRegression
from sklearn.metrics import accuracy_score, classification_report, f1_score
from sklearn.model_selection import train_test_split
from sklearn.neighbors import KNeighborsClassifier
from sklearn.pipeline import Pipeline
from sklearn.preprocessing import LabelEncoder, StandardScaler
from sklearn.svm import SVC

# Tentar importar XGBoost e LightGBM (opcionais)
try:
    import xgboost as xgb

    XGBOOST_AVAILABLE = True
except ImportError:
    XGBOOST_AVAILABLE = False

try:
    import lightgbm as lgb

    LIGHTGBM_AVAILABLE = True
except ImportError:
    LIGHTGBM_AVAILABLE = False


def load_data(data_path):
    """Carrega os dados de treinamento."""
    print(f"Carregando dados de {data_path}...")
    df = pd.read_csv(data_path)
    print(f"Dados carregados: {len(df)} registros")
    print(f"Número de categorias únicas: {df['category_id'].nunique()}")
    return df


def preprocess_data(df):
    """Preprocessa os dados."""
    # Remove linhas com valores nulos
    df = df.dropna()

    # Garante que name é string
    df["name"] = df["name"].astype(str)

    # Garante que amount_cents é numérico
    df["amount_cents"] = pd.to_numeric(df["amount_cents"], errors="coerce")

    # Remove linhas com valores inválidos
    df = df.dropna()

    print(f"Dados após pré-processamento: {len(df)} registros")
    return df


def create_preprocessor():
    """Cria o pré-processador comum para todos os modelos."""
    # Pipeline para processar texto (name)
    text_transformer = Pipeline(
        [
            (
                "tfidf",
                TfidfVectorizer(
                    max_features=500,
                    ngram_range=(1, 2),
                    min_df=2,
                    lowercase=True,
                    strip_accents="unicode",
                ),
            )
        ]
    )

    # Pipeline para processar features numéricas (amount_cents)
    numeric_transformer = Pipeline([("scaler", StandardScaler())])

    # Combina os transformers
    preprocessor = ColumnTransformer(
        transformers=[
            ("text", text_transformer, "name"),
            ("numeric", numeric_transformer, ["amount_cents"]),
        ],
        remainder="drop",
    )

    return preprocessor


def create_models():
    """Cria uma lista de modelos para testar."""
    preprocessor = create_preprocessor()

    models = {
        "Random Forest": Pipeline(
            [
                ("preprocessor", preprocessor),
                (
                    "classifier",
                    RandomForestClassifier(
                        n_estimators=200,
                        max_depth=20,
                        min_samples_split=5,
                        min_samples_leaf=2,
                        random_state=42,
                        n_jobs=-1,
                    ),
                ),
            ]
        ),
        "Gradient Boosting": Pipeline(
            [
                ("preprocessor", preprocessor),
                (
                    "classifier",
                    GradientBoostingClassifier(
                        n_estimators=100,
                        max_depth=5,
                        learning_rate=0.1,
                        random_state=42,
                    ),
                ),
            ]
        ),
        "Logistic Regression": Pipeline(
            [
                ("preprocessor", preprocessor),
                (
                    "classifier",
                    LogisticRegression(
                        max_iter=1000,
                        random_state=42,
                        n_jobs=-1,
                    ),
                ),
            ]
        ),
        "SVM (RBF)": Pipeline(
            [
                ("preprocessor", preprocessor),
                (
                    "classifier",
                    SVC(kernel="rbf", random_state=42, probability=True),
                ),
            ]
        ),
        "K-Nearest Neighbors": Pipeline(
            [
                ("preprocessor", preprocessor),
                ("classifier", KNeighborsClassifier(n_neighbors=5)),
            ]
        ),
    }

    # Adiciona XGBoost se disponível
    if XGBOOST_AVAILABLE:
        models["XGBoost"] = Pipeline(
            [
                ("preprocessor", preprocessor),
                (
                    "classifier",
                    xgb.XGBClassifier(
                        n_estimators=100,
                        max_depth=6,
                        learning_rate=0.1,
                        random_state=42,
                        n_jobs=-1,
                        eval_metric="mlogloss",
                    ),
                ),
            ]
        )

    # Adiciona LightGBM se disponível
    if LIGHTGBM_AVAILABLE:
        models["LightGBM"] = Pipeline(
            [
                ("preprocessor", preprocessor),
                (
                    "classifier",
                    lgb.LGBMClassifier(
                        n_estimators=100,
                        max_depth=6,
                        learning_rate=0.1,
                        random_state=42,
                        n_jobs=-1,
                        verbose=-1,
                    ),
                ),
            ]
        )

    return models


def train_and_evaluate_model(model, name, X_train, y_train, X_val, y_val):
    """Treina e avalia um modelo."""
    print(f"\n{'=' * 60}")
    print(f"Treinando {name}...")
    print(f"{'=' * 60}")

    start_time = time.time()

    # Para XGBoost e LightGBM, precisa codificar labels para começar em 0
    use_label_encoder = name in ["XGBoost", "LightGBM"]
    label_encoder = None

    if use_label_encoder:
        label_encoder = LabelEncoder()
        y_train_encoded = label_encoder.fit_transform(y_train)
    else:
        y_train_encoded = y_train

    # Treina o modelo
    model.fit(X_train, y_train_encoded)

    training_time = time.time() - start_time

    # Faz predições
    y_pred_encoded = model.predict(X_val)

    # Decodifica predições se necessário
    if use_label_encoder:
        y_pred = label_encoder.inverse_transform(y_pred_encoded)
    else:
        y_pred = y_pred_encoded

    # Calcula métricas usando labels originais
    accuracy = accuracy_score(y_val, y_pred)
    f1_macro = f1_score(y_val, y_pred, average="macro")
    f1_weighted = f1_score(y_val, y_pred, average="weighted")

    print(f"\nTempo de treinamento: {training_time:.2f} segundos")
    print(f"Acurácia: {accuracy:.4f}")
    print(f"F1-score (macro): {f1_macro:.4f}")
    print(f"F1-score (weighted): {f1_weighted:.4f}")

    # Armazena o encoder se usado
    result = {
        "name": name,
        "model": model,
        "accuracy": accuracy,
        "f1_macro": f1_macro,
        "f1_weighted": f1_weighted,
        "training_time": training_time,
        "y_pred": y_pred,
    }

    if use_label_encoder:
        result["label_encoder"] = label_encoder

    return result


def save_model(pipeline, model_name, model_dir="models"):
    """Salva o modelo treinado."""
    os.makedirs(model_dir, exist_ok=True)
    # Sanitiza o nome do modelo para nome de arquivo
    safe_name = model_name.lower().replace(" ", "_").replace("(", "").replace(")", "")
    model_path = os.path.join(model_dir, f"{safe_name}.pkl")
    joblib.dump(pipeline, model_path)
    print(f"Modelo salvo em: {model_path}")
    return model_path


def print_comparison(results):
    """Imprime uma tabela comparativa dos resultados."""
    print("\n" + "=" * 80)
    print("COMPARAÇÃO DE MODELOS")
    print("=" * 80)

    # Ordena por acurácia
    sorted_results = sorted(results, key=lambda x: x["accuracy"], reverse=True)

    print(
        f"\n{'Modelo':<25} {'Acurácia':<12} {'F1 (macro)':<12} {'F1 (weighted)':<15} {'Tempo (s)':<10}"
    )
    print("-" * 80)

    for result in sorted_results:
        print(
            f"{result['name']:<25} "
            f"{result['accuracy']:<12.4f} "
            f"{result['f1_macro']:<12.4f} "
            f"{result['f1_weighted']:<15.4f} "
            f"{result['training_time']:<10.2f}"
        )

    print("\n" + "=" * 80)
    best_model = sorted_results[0]
    print(f"\n🏆 MELHOR MODELO: {best_model['name']}")
    print(f"   Acurácia: {best_model['accuracy']:.4f}")
    print(f"   F1-score (macro): {best_model['f1_macro']:.4f}")
    print(f"   F1-score (weighted): {best_model['f1_weighted']:.4f}")
    print("=" * 80)


def main():
    """Função principal."""
    # Configurações
    data_path = Path(__file__).parent / "data" / "training_data.csv"
    # Salva os modelos na pasta models dentro de src
    model_dir = Path(__file__).parent.parent / "src" / "models"

    # Carrega dados
    df = load_data(data_path)

    # Preprocessa dados
    df = preprocess_data(df)

    # Separa features e target
    X = df[["name", "amount_cents"]]
    y = df["category_id"]

    # Verifica distribuição de classes
    class_counts = y.value_counts()
    print("\nDistribuição de classes:")
    print(f"Classes com apenas 1 amostra: {(class_counts == 1).sum()}")
    print(f"Classes com 2+ amostras: {(class_counts >= 2).sum()}")

    # Filtra classes com pelo menos 2 amostras para stratify funcionar
    # Classes com apenas 1 amostra vão apenas para treino
    valid_classes = class_counts[class_counts >= 2].index
    mask = y.isin(valid_classes)

    if mask.sum() < len(y):
        print(
            f"\nUsando {mask.sum()} amostras de classes com 2+ amostras para split estratificado"
        )
        print(
            f"{len(y) - mask.sum()} amostras de classes únicas serão adicionadas ao treino"
        )

    # Divide em treino e validação (apenas classes com 2+ amostras)
    X_train, X_val, y_train, y_val = train_test_split(
        X[mask], y[mask], test_size=0.2, random_state=42, stratify=y[mask]
    )

    # Adiciona amostras de classes únicas ao treino
    if not mask.all():
        unique_class_mask = ~mask
        X_train = pd.concat([X_train, X[unique_class_mask]], ignore_index=True)
        y_train = pd.concat([y_train, y[unique_class_mask]], ignore_index=True)

    print(f"\nConjunto de treino: {len(X_train)} registros")
    print(f"Conjunto de validação: {len(X_val)} registros")

    # Cria modelos para testar
    models = create_models()

    print(f"\n{'=' * 60}")
    print(f"Testando {len(models)} modelos...")
    print(f"{'=' * 60}")

    if XGBOOST_AVAILABLE:
        print("✓ XGBoost disponível")
    else:
        print("⚠ XGBoost não disponível (instale com: pip install xgboost)")

    if LIGHTGBM_AVAILABLE:
        print("✓ LightGBM disponível")
    else:
        print("⚠ LightGBM não disponível (instale com: pip install lightgbm)")

    # Treina e avalia todos os modelos
    results = []
    for name, model in models.items():
        try:
            result = train_and_evaluate_model(
                model, name, X_train, y_train, X_val, y_val
            )
            results.append(result)
        except Exception as e:
            print(f"\n❌ Erro ao treinar {name}: {str(e)}")
            continue

    # Compara resultados
    print_comparison(results)

    # Salva o melhor modelo
    if results:
        best_result = max(results, key=lambda x: x["accuracy"])
        print(f"\n💾 Salvando melhor modelo: {best_result['name']}...")
        save_model(best_result["model"], best_result["name"], model_dir)

        # Salva também como modelo padrão
        default_path = os.path.join(model_dir, "category_classifier.pkl")
        joblib.dump(best_result["model"], default_path)
        print(f"Modelo padrão salvo em: {default_path}")

        # Relatório detalhado do melhor modelo
        print(f"\n{'=' * 60}")
        print(f"RELATÓRIO DETALHADO - {best_result['name']}")
        print(f"{'=' * 60}")
        print(classification_report(y_val, best_result["y_pred"]))

        # Teste rápido com alguns exemplos usando o melhor modelo
        print("\n\n=== Teste com exemplos do conjunto de validação ===")
        sample_indices = np.random.choice(len(X_val), min(5, len(X_val)), replace=False)
        for idx in sample_indices:
            name = X_val.iloc[idx]["name"]
            amount = X_val.iloc[idx]["amount_cents"]
            true_category = y_val.iloc[idx]
            pred_category = best_result["model"].predict(
                pd.DataFrame([{"name": name, "amount_cents": amount}])
            )[0]
            status = "✓" if true_category == pred_category else "✗"
            print(
                f"{status} Nome: {name}, Valor: {amount}, "
                f"Categoria Real: {true_category}, Predição: {pred_category}"
            )

    print("\n✅ Treinamento concluído!")


if __name__ == "__main__":
    main()
