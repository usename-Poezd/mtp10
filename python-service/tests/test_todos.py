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


def test_get_todos_empty():
    with respx.mock(base_url=GO_URL) as mock:
        mock.get("/todos").mock(return_value=httpx.Response(200, json=[]))
        response = client.get("/todos")
    assert response.status_code == 200
    assert response.json() == []


def test_create_todo_success():
    go_response = {"id": "abc-123", "title": "Test", "completed": False}
    with respx.mock(base_url=GO_URL) as mock:
        mock.post("/todos").mock(return_value=httpx.Response(201, json=go_response))
        response = client.post("/todos", json={"title": "Test", "priority": "high"})
    assert response.status_code == 201
    data = response.json()
    assert data["id"] == "abc-123"
    assert data["title"] == "Test"
    assert data["priority"] == "high"
    assert data["created_at"] is not None


def test_create_todo_empty_title():
    response = client.post("/todos", json={"title": "", "priority": "low"})
    assert response.status_code == 422


def test_create_todo_invalid_priority():
    response = client.post("/todos", json={"title": "x", "priority": "urgent"})
    assert response.status_code == 422


def test_get_todo_by_id():
    go_response = {"id": "abc-123", "title": "Test", "completed": False}
    with respx.mock(base_url=GO_URL) as mock:
        mock.get("/todos/abc-123").mock(return_value=httpx.Response(200, json=go_response))
        response = client.get("/todos/abc-123")
    assert response.status_code == 200
    assert response.json()["id"] == "abc-123"


def test_get_todo_not_found():
    with respx.mock(base_url=GO_URL) as mock:
        mock.get("/todos/missing").mock(return_value=httpx.Response(404, json={"error": "not found"}))
        response = client.get("/todos/missing")
    assert response.status_code == 404


def test_update_todo():
    go_response = {"id": "abc-123", "title": "Updated", "completed": True}
    with respx.mock(base_url=GO_URL) as mock:
        mock.put("/todos/abc-123").mock(return_value=httpx.Response(200, json=go_response))
        response = client.put("/todos/abc-123", json={"title": "Updated", "completed": True})
    assert response.status_code == 200
    assert response.json()["completed"] is True


def test_delete_todo():
    with respx.mock(base_url=GO_URL) as mock:
        mock.delete("/todos/abc-123").mock(return_value=httpx.Response(204))
        response = client.delete("/todos/abc-123")
    assert response.status_code == 204


def test_go_service_unavailable():
    with respx.mock(base_url=GO_URL) as mock:
        mock.get("/todos").mock(side_effect=httpx.ConnectError("connection refused"))
        response = client.get("/todos")
    assert response.status_code == 503
    assert "unavailable" in response.json()["detail"].lower()
