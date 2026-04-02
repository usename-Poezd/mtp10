from datetime import datetime, timezone

import httpx
from fastapi import FastAPI, HTTPException

from app.models import StatsResponse, TodoCreate, TodoResponse, TodoUpdate
from app.service import GoTodoService

_todo_meta: dict[str, dict] = {}

app = FastAPI(title="FastAPI TODO Service")


@app.get("/todos", response_model=list[TodoResponse])
async def get_todos():
    async with httpx.AsyncClient() as client:
        svc = GoTodoService(client)
        try:
            todos = await svc.get_all()
        except httpx.ConnectError as exc:
            raise HTTPException(status_code=503, detail="Go service unavailable") from exc
    result = []
    for t in todos:
        meta = _todo_meta.get(t["id"], {})
        result.append(
            TodoResponse(
                id=t["id"],
                title=t["title"],
                completed=t["completed"],
                created_at=meta.get("created_at"),
                priority=meta.get("priority"),
            )
        )
    return result


@app.post("/todos", response_model=TodoResponse, status_code=201)
async def create_todo(body: TodoCreate):
    async with httpx.AsyncClient() as client:
        svc = GoTodoService(client)
        try:
            created = await svc.create(body.title)
        except httpx.ConnectError as exc:
            raise HTTPException(status_code=503, detail="Go service unavailable") from exc
    now = datetime.now(timezone.utc).isoformat()
    _todo_meta[created["id"]] = {"created_at": now, "priority": body.priority}
    return TodoResponse(
        id=created["id"],
        title=created["title"],
        completed=created["completed"],
        created_at=now,
        priority=body.priority,
    )


@app.get("/todos/{todo_id}", response_model=TodoResponse)
async def get_todo(todo_id: str):
    async with httpx.AsyncClient() as client:
        svc = GoTodoService(client)
        try:
            todo = await svc.get_by_id(todo_id)
        except httpx.ConnectError as exc:
            raise HTTPException(status_code=503, detail="Go service unavailable") from exc
        except httpx.HTTPStatusError as exc:
            raise HTTPException(status_code=exc.response.status_code, detail=exc.response.text) from exc
    meta = _todo_meta.get(todo_id, {})
    return TodoResponse(
        id=todo["id"],
        title=todo["title"],
        completed=todo["completed"],
        created_at=meta.get("created_at"),
        priority=meta.get("priority"),
    )


@app.put("/todos/{todo_id}", response_model=TodoResponse)
async def update_todo(todo_id: str, body: TodoUpdate):
    async with httpx.AsyncClient() as client:
        svc = GoTodoService(client)
        try:
            updated = await svc.update(todo_id, body.title, body.completed)
        except httpx.ConnectError as exc:
            raise HTTPException(status_code=503, detail="Go service unavailable") from exc
        except httpx.HTTPStatusError as exc:
            raise HTTPException(status_code=exc.response.status_code, detail=exc.response.text) from exc
    meta = _todo_meta.get(todo_id, {})
    return TodoResponse(
        id=updated["id"],
        title=updated["title"],
        completed=updated["completed"],
        created_at=meta.get("created_at"),
        priority=meta.get("priority"),
    )


@app.delete("/todos/{todo_id}", status_code=204)
async def delete_todo(todo_id: str):
    async with httpx.AsyncClient() as client:
        svc = GoTodoService(client)
        try:
            await svc.delete(todo_id)
        except httpx.ConnectError as exc:
            raise HTTPException(status_code=503, detail="Go service unavailable") from exc
        except httpx.HTTPStatusError as exc:
            raise HTTPException(status_code=exc.response.status_code, detail=exc.response.text) from exc
    _todo_meta.pop(todo_id, None)


@app.get("/stats", response_model=StatsResponse)
async def get_stats():
    async with httpx.AsyncClient() as client:
        svc = GoTodoService(client)
        try:
            todos = await svc.get_all()
        except httpx.ConnectError as exc:
            raise HTTPException(status_code=503, detail="Go service unavailable") from exc
    total = len(todos)
    completed_count = sum(1 for t in todos if t["completed"])
    by_priority: dict[str, int] = {"low": 0, "medium": 0, "high": 0}
    for t in todos:
        meta = _todo_meta.get(t["id"], {})
        priority = meta.get("priority", "medium")
        by_priority[priority] = by_priority.get(priority, 0) + 1
    return StatsResponse(
        total=total,
        completed=completed_count,
        pending=total - completed_count,
        by_priority=by_priority,
    )
