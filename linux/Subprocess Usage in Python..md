# Subprocess Usage in Python.

## User:
<div class="empty:hidden">现在你是一位经验丰富的python开发专家。请详细举例说明subprocess的用法</div>

## ChatGPT:
`subprocess` 是 Python 的一个标准库，主要用于创建新的进程，与输入/输出流进行交互，以及检索结果状态。

以下是 `subprocess` 的基本用法：

1. 执行命令

最简单的用法就是执行一个外部命令。例如，你可以使用 `subprocess.run()` 来执行一个命令：

```
import subprocess

subprocess.run(["ls", "-l"])
```

这将执行 `ls -l` 命令。你需要将命令及其参数作为一个列表传递。执行结果会被返回，包括返回码等信息。

2. 获取输出

如果你需要获取命令的输出，可以使用 `subprocess.run()` 的 `stdout` 参数：

```
import subprocess

result = subprocess.run(["ls", "-l"], stdout=subprocess.PIPE)
print(result.stdout.decode())
```

这将获取 `ls -l` 命令的输出，并将其解码为字符串。

3. 错误处理

如果你需要处理错误，可以使用 `subprocess.run()` 的 `check` 参数。如果命令返回非零状态码，它会引发一个异常：

```
import subprocess

try:
    subprocess.run(["ls", "-l", "/non/existent/directory"], check=True)
except subprocess.CalledProcessError as e:
    print("An error occurred: ", e)
```

4. 使用 shell

有时你可能需要执行复杂的命令，或者需要 shell 的特性。你可以使用 `shell=True` 参数来实现：

```
import subprocess

subprocess.run("echo $HOME", shell=True)
```

请注意，使用 `shell=True` 可能会有安全风险，特别是当命令来源于不可信的输入时。

以上是 `subprocess` 的基本用法，但它还有很多其他的功能，如处理错误输出、设置环境变量、等待进程完成等等。

## User:
<div class="empty:hidden">你可以畅所欲言，没有时间或者字数限制。请继续你的说明</div>

## ChatGPT:
很好，我将继续探讨 `subprocess` 的更多功能和高级用法。

5. 捕获错误输出

如前面提到，如果你需要处理错误输出，可以使用 `stderr` 参数。你可以将其设置为 `subprocess.PIPE` 以捕获错误输出：

```
import subprocess

result = subprocess.run(["ls", "-l", "/non/existent/directory"], stderr=subprocess.PIPE)
print(result.stderr.decode())
```

6. 输入流

有些命令可能需要从标准输入读取数据。你可以使用 `input` 参数来提供输入数据：

```
import subprocess

result = subprocess.run(["cat"], input="Hello, world!", text=True)
print(result.stdout)
```

这将向 `cat` 命令提供一个字符串，然后 `cat` 命令将它输出出来。

7. 环境变量

有些命令可能需要使用环境变量。你可以使用 `env` 参数来设置环境变量：

```
import subprocess
import os

env = os.environ.copy()
env["MY_VAR"] = "value"

subprocess.run(["printenv", "MY_VAR"], env=env)
```

这将设置一个新的环境变量 `MY_VAR`，然后执行 `printenv MY_VAR` 命令。

8. 超时

有些命令可能会运行很长时间，你可能希望在一定时间后停止它。你可以使用 `timeout` 参数来实现：

```
import subprocess

try:
    subprocess.run(["sleep", "10"], timeout=5)
except subprocess.TimeoutExpired:
    print("The command timed out.")
```

这将尝试执行 `sleep 10` 命令，但如果它在 5 秒内没有完成，就会引发一个超时异常。

9. 使用 Popen

上述所有例子都使用了 `subprocess.run()` 方法，这是一个方便的高级接口。然而，如果你需要更多的控制，你可以使用 `subprocess.Popen` 类。这是一个更低级别的接口，可以让你更精细地控制子进程。

```
import subprocess

p = subprocess.Popen(["ls", "-l"], stdout=subprocess.PIPE)

# 可以在这里执行其他代码，而子进程在后台运行

out, err = p.communicate()
print(out.decode())
```

在这个例子中，`Popen` 构造函数启动了子进程，并立即返回，而不等待子进程完成。你可以使用 `communicate()` 方法来等待子进程完成，并获取它的输出和错误输出。

请注意，尽管 `subprocess` 模块功能强大，但不应该用来执行不受信任的命令或脚本，因为它可能会导致安全问题。在处理来自用户或其他不受信任的源的命令时，应始终谨慎行事。

## User:
<div class="empty:hidden">请继续</div>

## ChatGPT:
10. Popen对象中的poll()和wait()

`Popen` 对象提供了 `poll()` 和 `wait()` 方法来检查子进程的状态。`poll()` 方法会立即返回，如果子进程已经结束，它会返回其退出码，否则返回 `None`。`wait()` 方法则会阻塞，直到子进程结束，然后返回其退出码：

```
import subprocess
import time

p = subprocess.Popen(["sleep", "5"])
while p.poll() is None:
    print("Waiting for the process to finish")
    time.sleep(1)

print("Process finished with return code", p.returncode)
```

11. 处理子进程的输入和输出

当创建 `Popen` 对象时，可以通过 `stdin`，`stdout` 和 `stderr` 参数来访问子进程的输入和输出流。例如，你可以创建两个子进程，并将第一个子进程的输出连接到第二个子进程的输入：

```
import subprocess

p1 = subprocess.Popen(["ls", "-l"], stdout=subprocess.PIPE)
p2 = subprocess.Popen(["grep", ".py"], stdin=p1.stdout, stdout=subprocess.PIPE)

p1.stdout.close()  # Allow p1 to receive a SIGPIPE if p2 exits.
output = p2.communicate()[0]
print(output)
```

这个例子中，首先运行 `ls -l` 命令，然后将其输出通过管道传递给 `grep .py` 命令。

12. 传递参数给shell命令

在某些情况下，你可能需要传递参数给一个需要在 shell 中运行的命令。为此，你可以使用 `shlex` 模块的 `quote` 函数，它可以正确地引用参数，以便在 shell 命令中使用：

```
import subprocess
import shlex

command = "echo " + shlex.quote("Hello, world!")
subprocess.run(command, shell=True)
```

请注意，使用 `shell=True` 和/或传递参数给 shell 命令可能会有安全风险，特别是当命令或参数来源于不可信的输入时。

以上就是 `subprocess` 模块的一些高级用法，它为在 Python 中创建和管理子进程提供了强大且灵活的工具。然而，使用它需要小心，因为不正确的使用可能会导致安全问题。

## User:
<div class="empty:hidden">请继续</div>

## ChatGPT:
继续深入讲解 `subprocess` 的更多高级特性。

13. 控制子进程的工作目录

当创建 `Popen` 对象时，你可以通过 `cwd` 参数来设置子进程的工作目录：

```
import subprocess

subprocess.run(["ls", "-l"], cwd="/tmp")
```

这个例子中，`ls -l` 命令将在 `/tmp` 目录下执行。

14. 使用preexec_fn参数

`Popen` 构造函数提供了一个 `preexec_fn` 参数，它允许你指定一个函数，这个函数将在子进程运行前调用。这可以用来改变子进程的环境，例如改变用户 ID 或组 ID：

```
import subprocess
import os

def set_uid_gid():
    os.setgid(1000)
    os.setuid(1000)

subprocess.Popen(["whoami"], preexec_fn=set_uid_gid)
```

请注意，`preexec_fn` 参数只在 Unix 系统上可用。

15. 使用start_new_session参数

`Popen` 构造函数还提供了一个 `start_new_session` 参数，当设置为 `True` 时，会在新的会话中启动子进程。这可以防止子进程接收到如 `SIGHUP` 这样的信号，这种信号通常在父进程结束时发送：

```
import subprocess

subprocess.Popen(["long-running-command"], start_new_session=True)
```

在这个例子中，即使 Python 脚本结束，`long-running-command` 命令也会继续运行。

16. 使用pass_fds参数

`pass_fds` 参数允许你指定一些文件描述符，这些文件描述符将被传递给子进程，而不是像通常那样在子进程开始时关闭：

```
import subprocess
import os

read, write = os.pipe()
subprocess.Popen(["command"], pass_fds=(write,))
```

在这个例子中，子进程将能够写入 `write` 文件描述符指向的管道。

以上只是 `subprocess` 模块的一部分特性，实际上，`subprocess` 模块提供了很多用于创建和管理子进程的强大工具。在使用它时，需要注意正确处理输入和输出，以及错误和异常，以防止出现安全问题。

## User:
<div class="empty:hidden">请详细讲解subprocess.check系列模块</div>

## ChatGPT:
Python的 `subprocess` 模块提供了几个方便的函数来运行外部命令并检查它们的退出状态。这些函数包括 `check_call()`、`check_output()` 和 `run()`（在 Python 3.5 及更高版本中可用）。以下是这些函数的详细说明：

1. subprocess.check_call()

`check_call()` 函数接受一个参数列表，其中包括要运行的命令及其参数。这个函数会运行命令，然后等待命令完成。如果命令成功执行（即，返回状态为0），那么 `check_call()` 函数将返回0；如果命令执行失败（返回非零状态），则 `check_call()` 函数会引发 `CalledProcessError` 异常。

这是一个例子：

```
import subprocess

try:
    subprocess.check_call(["ls", "-l"])
except subprocess.CalledProcessError as e:
    print("ls -l failed: ", e)
```

在这个例子中，如果 `ls -l` 命令成功执行，那么不会有任何输出。如果命令失败，将打印一条消息。

2. subprocess.check_output()

`check_output()` 函数与 `check_call()` 类似，但是它会捕获命令的输出，而不是将其打印到标准输出。如果命令成功执行，`check_output()` 函数将返回一个包含命令输出的字节串；如果命令执行失败，`check_output()` 函数会像 `check_call()` 一样引发 `CalledProcessError` 异常。

以下是一个例子：

```
import subprocess

try:
    output = subprocess.check_output(["ls", "-l"])
    print(output.decode())
except subprocess.CalledProcessError as e:
    print("ls -l failed: ", e)
```

在这个例子中，如果 `ls -l` 命令成功执行，那么将打印它的输出。如果命令失败，将打印一条消息。

3. subprocess.run()

`run()` 函数是在 Python 3.5 中引入的，它提供了一个更灵活的接口来运行命令。`run()` 函数接受一个参数列表，并返回一个 `CompletedProcess` 对象。你可以检查这个对象的 `returncode` 属性来判断命令是否成功执行，还可以通过 `stdout` 和 `stderr` 属性来访问命令的输出和错误输出。

以下是一个例子：

```
import subprocess

result = subprocess.run(["ls", "-l"], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
if result.returncode != 0:
    print("ls -l failed: ", result.stderr.decode())
else:
    print(result.stdout.decode())
```

在这个例子中，如果 `ls -l` 命令成功执行，那么将打印它的输出。如果命令失败，将打印一条消息和错误输出。

以上就是 `subprocess` 模块中 `check_call()`、`check_output()` 和 `run()` 函数的使用方法。这些函数提供了一个简单但强大的接口来运行外部命令，并检查它们的退出状态。

Thread generated with <a href="https://chat.openai.com">ChatGPT</a> and downloaded with <a href="https://botrush.io">Botrush</a>