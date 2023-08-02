from typing import overload


@overload
def get(v: int) -> int:
    ...


@overload
def get(v: str) -> str:
    ...
