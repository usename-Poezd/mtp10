import os
from typing import Any

import httpx

GO_SERVICE_URL = os.getenv("GO_SERVICE_URL", "http://localhost:8080")


class GoTodoService:
    def __init__(self, client: httpx.AsyncClient):
        self._client = client
        self._base_url = GO_SERVICE_URL

    async def get_all(self) -> list[dict[str, Any]]:
        response = await self._client.get(f"{self._base_url}/todos")
        response.raise_for_status()
        return response.json()

    async def get_by_id(self, todo_id: str) -> dict[str, Any]:
        response = await self._client.get(f"{self._base_url}/todos/{todo_id}")
        response.raise_for_status()
        return response.json()

    async def create(self, title: str) -> dict[str, Any]:
        response = await self._client.post(
            f"{self._base_url}/todos",
            json={"title": title},
        )
        response.raise_for_status()
        return response.json()

    async def update(self, todo_id: str, title: str, completed: bool) -> dict[str, Any]:
        response = await self._client.put(
            f"{self._base_url}/todos/{todo_id}",
            json={"title": title, "completed": completed},
        )
        response.raise_for_status()
        return response.json()

    async def delete(self, todo_id: str) -> None:
        response = await self._client.delete(f"{self._base_url}/todos/{todo_id}")
        response.raise_for_status()
