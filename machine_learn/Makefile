install:
	poetry install

train:
	poetry run python train_model.py

predict:
	poetry run python predict.py $(name) $(amount)

dev:
	poetry run fastapi dev src/api/app.py --host 0.0.0.0 --port 8080

run:
	poetry run uvicorn src.api.app:app --reload --host 0.0.0.0 --port 8000

test:
	poetry run pytest -v

