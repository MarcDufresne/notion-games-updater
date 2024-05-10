FROM python:3.11-slim

RUN apt update \
    && apt install -y build-essential \
    && apt clean

WORKDIR /app
ADD pyproject.toml ./
ADD poetry.lock ./
ADD notion_games/__init__.py ./notion_games/__init__.py

RUN pip install --no-cache-dir -U pip poetry \
    && poetry config virtualenvs.in-project true \
    && poetry install --no-dev

ADD notion_games ./notion_games
ADD run.py ./

ENTRYPOINT ["poetry", "run", "python", "run.py"]
