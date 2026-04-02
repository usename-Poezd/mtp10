import pytest
from pydantic import ValidationError

from app.models import TodoCreate


def test_create_todo_empty_title_validation():
    with pytest.raises(ValidationError):
        TodoCreate(title="", priority="low")


def test_create_todo_invalid_priority_validation():
    with pytest.raises(ValidationError):
        TodoCreate(title="x", priority="urgent")
