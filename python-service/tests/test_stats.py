import httpx
import pytest
import respx
from fastapi.testclient import TestClient

from app.main import _todo_meta, app

client = TestClient(app)

GO_URL = "http://localhost:8080"


@pytest.fixture(autouse=True)
def clear_meta():
    _todo_meta.clear()
    yield
    _todo_meta.clear()


def test_stats_empty():
    with respx.mock(base_url=GO_URL) as mock:
        mock.get("/todos").mock(return_value=httpx.Response(200, json=[]))
        response = client.get("/stats")
    assert response.status_code == 200
    data = response.json()
    assert data["total"] == 0
    assert data["completed"] == 0
    assert data["pending"] == 0


def test_stats_with_todos():
    todos = [
        {"id": "1", "title": "A", "completed": True},
        {"id": "2", "title": "B", "completed": False},
        {"id": "3", "title": "C", "completed": False},
    ]
    _todo_meta["1"] = {"priority": "high", "created_at": "2026-01-01T00:00:00+00:00"}
    _todo_meta["2"] = {"priority": "low", "created_at": "2026-01-02T00:00:00+00:00"}
    with respx.mock(base_url=GO_URL) as mock:
        mock.get("/todos").mock(return_value=httpx.Response(200, json=todos))
        response = client.get("/stats")
    assert response.status_code == 200
    data = response.json()
    assert data["total"] == 3
    assert data["completed"] == 1
    assert data["pending"] == 2
    assert data["by_priority"]["high"] == 1
    assert data["by_priority"]["low"] == 1
    assert data["by_priority"]["medium"] == 1
