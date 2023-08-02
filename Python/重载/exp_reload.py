def get(v):
    if isinstance(v, int):
        return v + 1
    elif isinstance(v, str):
        return v + "1"
    else:
        raise TypeError("Invalid type")


b = get(1)
a = get("1")