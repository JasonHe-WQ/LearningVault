import asyncio
import sys

import uvloop

async def main():
    # Main entry-point.
    ...

if sys.version_info >= (3, 11):
    with asyncio.Runner(loop_factory=uvloop.new_event_loop) as runner:
        runner.run(main())
else:
    uvloop.install()
    asyncio.run(main())