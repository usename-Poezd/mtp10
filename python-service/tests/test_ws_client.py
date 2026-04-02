import asyncio
import importlib
import json
import threading
from typing import Any

import pytest

websockets: Any = importlib.import_module("websockets")


async def echo_server(websocket):
    async for message in websocket:
        await websocket.send(message)


@pytest.fixture
def ws_server():
    loop = asyncio.new_event_loop()
    server = None
    host = "127.0.0.1"
    ready = threading.Event()
    ws_url = None

    async def start():
        nonlocal server, ws_url
        server = await websockets.serve(echo_server, host, 0)
        sock = server.sockets[0]
        ws_url = f"ws://{sock.getsockname()[0]}:{sock.getsockname()[1]}"
        ready.set()

    def run():
        loop.run_until_complete(start())
        loop.run_forever()

    thread = threading.Thread(target=run, daemon=True)
    thread.start()
    if not ready.wait(timeout=3):
        raise RuntimeError("websocket test server did not start")

    yield ws_url

    if server is not None:
        close_future = asyncio.run_coroutine_threadsafe(server.wait_closed(), loop)
        loop.call_soon_threadsafe(server.close)
        close_future.result(timeout=1)
    loop.call_soon_threadsafe(loop.stop)
    thread.join(timeout=1)


def test_ws_client_sends_and_receives(ws_server):
    received = []

    async def test():
        async with websockets.connect(ws_server) as ws:
            payload = json.dumps(
                {"type": "message", "username": "TestUser", "text": "Hello test"}
            )
            await ws.send(payload)
            response = await asyncio.wait_for(ws.recv(), timeout=3.0)
            data = json.loads(response)
            received.append(data)

    asyncio.run(test())
    assert len(received) == 1
    assert received[0]["username"] == "TestUser"
    assert received[0]["text"] == "Hello test"


def test_ws_client_json_format(ws_server):
    async def test():
        async with websockets.connect(ws_server) as ws:
            payload = json.dumps({"type": "message", "username": "Alice", "text": "Hi"})
            await ws.send(payload)
            raw = await asyncio.wait_for(ws.recv(), timeout=3.0)
            data = json.loads(raw)
            assert data.get("type") == "message"
            assert data.get("username") == "Alice"
            assert data.get("text") == "Hi"

    asyncio.run(test())
