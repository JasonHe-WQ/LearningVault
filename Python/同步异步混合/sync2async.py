import asyncio
import time




async def root():
    loop = asyncio.get_event_loop()

    def do_blocking_work():
        time.sleep(10)
        print("Done blocking work!!")

    await loop.run_in_executor(None, do_blocking_work)
    return {"message": "Hello World"}


async def health():
    return {"status": "ok"}
