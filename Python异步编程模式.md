# Python异步编程模式

## USER:

请详细解释puthon中的async wait模式，请问这是5种io中的哪一种

## CHATGPT:

`asyncio` 是 Python 中一个异步编程框架，它允许你编写并发（concurrent）的代码。在这个框架下，`async` 和 `await` 是两个关键字，它们让你能够编写异步（asynchronous）代码。

`async` 关键字用于定义一个异步函数。异步函数与普通函数不同，它们返回一个协程对象（coroutine object），这个对象可以在后面与 `await` 关键字一起使用。

`await` 关键字用于等待一个异步操作（协程对象）的结果。当一个异步操作被 `await`，程序会在这一点暂停，并等待该操作完成。在操作完成后，程序会回到暂停的地方继续执行。值得注意的是，`await` 只能在 `async` 函数内部使用。

下面是一个简单的示例：

```
import asyncio

async def foo():
    print("Start foo")
    await asyncio.sleep(1)
    print("End foo")

async def main():
    await foo()

asyncio.run(main())

```

在这个例子中，`foo` 函数是一个异步函数，它在开始时打印 "Start foo"，然后等待 1 秒（通过 `await asyncio.sleep(1)`），最后打印 "End foo"。`main` 函数也是一个异步函数，它调用并等待 `foo` 函数完成。最后，我们通过 `asyncio.run(main())` 启动事件循环并运行 `main` 函数。

`async/await` 模式是五种 I/O 模型（阻塞 I/O、非阻塞 I/O、I/O 复用、信号驱动 I/O 和异步 I/O）中的异步 I/O（Asynchronous I/O）模型。在异步 I/O 模型中，当发起 I/O 操作时，进程不会被阻塞，可以继续执行其他任务。当 I/O 操作完成时，操作系统会通知进程，进程可以继续处理 I/O 操作的结果。这种模型允许编写高并发的程序，并有效地利用系统资源。

