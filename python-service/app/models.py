from typing import Literal

from pydantic import BaseModel, field_validator


class TodoCreate(BaseModel):
    title: str
    priority: Literal["low", "medium", "high"] = "medium"

    @field_validator("title")
    @classmethod
    def title_not_empty(cls, v: str) -> str:
        if not v.strip():
            raise ValueError("title must not be empty")
        return v


class TodoUpdate(BaseModel):
    title: str
    completed: bool

    @field_validator("title")
    @classmethod
    def title_not_empty(cls, v: str) -> str:
        if not v.strip():
            raise ValueError("title must not be empty")
        return v


class TodoResponse(BaseModel):
    id: str
    title: str
    completed: bool
    created_at: str | None = None
    priority: str | None = None


class StatsResponse(BaseModel):
    total: int
    completed: int
    pending: int
    by_priority: dict[str, int]
