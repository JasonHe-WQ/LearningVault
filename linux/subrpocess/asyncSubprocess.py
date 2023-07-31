import asyncio


# This is the function that uses asyncio.create_subprocess_shell
async def run_cmd(command):
    process = await asyncio.create_subprocess_shell(
        command,
        stdout=asyncio.subprocess.PIPE,
        stderr=asyncio.subprocess.PIPE
    )

    stdout, stderr = await process.communicate()

    if process.returncode != 0:
        print(f"Command failed with exit code {process.returncode}")
        print(f"stderr: {stderr.decode()}")
        return None

    return stdout.decode()


# This is another function that calls the first function
async def execute_commands(commands):
    tasks = [run_cmd(cmd) for cmd in commands]
    results = await asyncio.gather(*tasks)

    for i, result in enumerate(results):
        print(f"Result of command {i}: {result}")


# List of commands to run
commands = [
    'ls',
    'pwd',
    # Add more commands as needed
]

# Run the execute_commands function
asyncio.run(execute_commands(commands))
