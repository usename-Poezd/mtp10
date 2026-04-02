#!/usr/bin/env python3

import argparse
import asyncio
import importlib
import json
from typing import Any

websockets: Any = importlib.import_module("websockets")


async def run_client(url: str, username: str, message: str) -> None:
    async with websockets.connect(url) as websocket:
        payload = json.dumps({"type": "message", "username": username, "text": message})
        await websocket.send(payload)
        response = await asyncio.wait_for(websocket.recv(), timeout=5.0)
        data = json.loads(response)
        print(f"[{data.get('username')}]: {data.get('text')}")


def main() -> None:
    parser = argparse.ArgumentParser(description="WebSocket chat client")
    parser.add_argument("--url", default="ws://localhost:8081/ws", help="WebSocket URL")
    parser.add_argument("--username", required=True, help="Username")
    parser.add_argument("--message", required=True, help="Message to send")
    args = parser.parse_args()
    asyncio.run(run_client(args.url, args.username, args.message))


if __name__ == "__main__":
    main()
