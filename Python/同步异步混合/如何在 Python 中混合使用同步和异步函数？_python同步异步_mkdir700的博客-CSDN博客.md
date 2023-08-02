# 如何在 Python 中混合使用同步和异步函数？_python同步异步_mkdir700的博客

##### 链接

*   [在线程或者进程池中执行代码。](https://docs.python.org/zh-cn/3/library/asyncio-eventloop.html#executing-code-in-thread-or-process-pools)
*   [concurrent.futures.Future](https://docs.python.org/zh-cn/3/library/concurrent.futures.html#concurrent.futures.Future)
*   [add_done_callback](https://docs.python.org/zh-cn/3/library/asyncio-task.html#asyncio.Task.add_done_callback)

##### 标记

前言
--

异步编程可以提高应用程序的性能和吞吐量，因为它可以充分利用 CPU 和 I/O 资源。当某个任务被阻塞时，事件循环可以切换到另一个任务，从而避免浪费 CPU 时间。此外，异步编程还可以简化代码，使其更易于维护和调试。

我们最常用的是同步编程，在同步场景中，某个任务被阻塞时，整个线程都会被挂起，直到该任务完成，所以为了避免整个程序被阻塞的情况，又引入了多线程和锁。同步编程通常需要使用锁和其他同步原语来确保线程安全。

**混合编写的场景**  
在实际开发过程中，通常会遇到同时进行异步和同步操作的场景。例如，在使用异步 Web 框架 FastAPI 时，对于实际的一些业务需求将不得不采用同步编程。例如，可能需要使用同步数据库驱动程序或同步文件系统接口，又或者调用同步接口等等。

在协程函数中调用同步函数
------------

在协程函数中直接调用同步函数会阻塞事件循环，从而影响整个程序的性能。我们先来看一个例子：

以下是使用异步 Web 框架 FastAPI 写的一个例子，FastAPI 是比较快，但不正确的操作将会变得很慢。

```python
import time

from fastapi import FastAPI

app = FastAPI()

@app.get("/")
async def root():
    time.sleep(10)
    return {"message": "Hello World"}

@app.get("/health")
async def health():
    return {"status": "ok"}
```

上面我们写了两个接口，假设 `root` 接口函数耗时 10 秒，在这 10 秒内访问 `health` 接口，想一想会发生什么？

访问 `root` 接口（左），立即访问 `health` 接口（右），`health` 接口被阻塞，直至 `root` 接口返回后，`health` 接口才成功响应。

`time.sleep` 就是一个「同步」函数，它会阻塞整个事件循环。

如何解决呢？想一想以前的处理方法，如果一个函数会阻塞主线程，那么就再开一个线程让这个阻塞函数单独运行。所以，这里也是同理，开一个线程单独去运行那些阻塞式操作，比如读取文件等。

`loop.run_in_executor` 方法将同步函数转换为异步非阻塞方式进行处理。具体来说，`loop.run_in_executor()` 可以将同步函数创建为**一个线程**或**进程**，并在其中执行该函数，从而避免阻塞事件循环。

> 官方例子：[在线程或者进程池中执行代码。](https://docs.python.org/zh-cn/3/library/asyncio-eventloop.html#executing-code-in-thread-or-process-pools)

那么，我们使用 `loop.run_in_executor` 改写上面例子，如下：

```python
import asyncio
import time

from fastapi import FastAPI

app = FastAPI()

@app.get("/")
async def root():
    loop = asyncio.get_event_loop()

    def do_blocking_work():
        time.sleep(10)
        print("Done blocking work!!")

    await loop.run_in_executor(None, do_blocking_work)
    return {"message": "Hello World"}

@app.get("/health")
async def health():
    return {"status": "ok"}
```
效果如下：

`root` 接口被阻塞期间，`health` 依然正常访问互不影响。

> **注意：** 这里都是为了演示，实际在使用 FastAPI 开发时，你可以直接将 `async def root` 更换成 `def root` ，也就是将其换成同步接口函数，FastAPI 内部会自动创建线程处理这个同步接口函数。总的来说，FastAPI 内部也是依靠线程去处理同步函数从而避免阻塞主线程（或主线程中的事件循环）。

在同步函数中调用异步函数
------------

协程只能在「事件循环」内被执行，且同一时刻只能有一个协程被执行。

所以，在同步函数中调用异步函数，其本质就是将协程「扔进」事件循环中，等待该协程执行完获取结果即可。

以下这些函数，都可以实现这个效果：

*   `asyncio.run`
*   `asyncio.run_coroutine_threadsafe`
*   `loop.run_until_complete`
*   `create_task`

接下来，我们将一一讲解这些方法并举例说明。

### asyncio.run

这个方法使用起来最简单，先看下如何使用，然后紧跟着讲一下哪些场景不能直接使用 `asyncio.run`
```python
import asyncio

async def do_work():
    return 1

def main():
    result = asyncio.run(do_work())
    print(result)  # 1

if __name__ == "__main__":
    main()
```
直接 `run` 就完事了，然后接受返回值即可。

但是需要，注意的是 `asyncio.run` 每次调用都会新开一个事件循环，当结束时自动关闭该事件循环。

**一个线程内只存在一个事件循环**，所以如果当前线程已经有存在的事件循环了，就不应该使用 `asyncio.run` 了，否则就会抛出如下异常：

RuntimeError: asyncio.run() cannot be called from a running event loop

因此，`asyncio.run` 用作新开一个事件循环时使用。

### asyncio.run_coroutine_threadsafe

文档： https://docs.python.org/zh-cn/3/library/asyncio-task.html#asyncio.run_coroutine_threadsafe

> 向指定事件循环提交一个协程。（线程安全）  
> 返回一个 [`concurrent.futures.Future`](https://docs.python.org/zh-cn/3/library/concurrent.futures.html#concurrent.futures.Future "concurrent. futures. Future") 以等待来自其他 OS 线程的结果。

换句话说，就是**将协程丢给其他线程中的事件循环去运行**。

值得注意的是这里的「事件循环」应该是其他线程中的事件循环，非当前线程的事件循环。

其返回的结果是一个 future 对象，如果你需要获取协程的执行结果可以使用 `future.result()` 获取，关于 future 对象的更多介绍，见 https://docs.python.org/zh-cn/3/library/concurrent.futures.html#concurrent.futures.Future

下方给了一个例子，一共有两个线程：`thread_with_loop` 和 `another_thread`，分别用于启动事件循环和调用 `run_coroutine_threadsafe`
```python
import asyncio
import threading
import time

loop = None

def get_loop():
    global loop
    if loop is None:
        loop = asyncio.new_event_loop()
    return loop

def another_thread():
    async def coro_func():
        return 1

    loop = get_loop()
    # 将协程提交到另一个线程的事件循环中执行
    future = asyncio.run_coroutine_threadsafe(coro_func(), loop)
    # 等待协程执行结果
    print(future.result())
    # 停止事件循环
    loop.call_soon_threadsafe(loop.stop)

def thread_with_loop():
    loop = get_loop()
    # 启动事件循环，确保事件循环不会退出，直到 loop.stop() 被调用
    loop.run_forever()
    loop.close()

# 启动一个线程，线程内部启动了一个事件循环
threading.Thread(target=thread_with_loop).start()
time.sleep(1)
# 在主线程中启动一个协程, 并将协程提交到另一个线程的事件循环中执行
t = threading.Thread(target=another_thread)
t.start()
t.join()
```
### loop.run_until_complete

文档： https://docs.python.org/zh-cn/3.10/library/asyncio-eventloop.html#asyncio.loop.run_until_complete

> 运行直到 _future_ ( [`Future`](https://docs.python.org/zh-cn/3.10/library/asyncio-future.html#asyncio.Future "asyncio.Future") 的实例 ) 被完成。

> 这个方法和 `asyncio.run` 类似。

具体就是传入一个协程对象或者任务，然后可以直接拿到协程的返回值。

`run_until_complete` 属于 `loop` 对象的方法，所以这个方法的使用前提是有一个事件循环，注意这个事件循环必须是**非运行状态**，如果是运行中就会抛出如下异常：

RuntimeError: This event loop is already running

例子：

loop = asyncio.new_event_loop()
loop.run_until_complete(do_async_work())

### create_task

文档： https://docs.python.org/zh-cn/3/library/asyncio-task.html#creating-tasks

再次准确一点：要**运行一个协程函数的本质是将携带协程函数的任务提交至事件循环中，由事件循环发现、调度并执行。**

其实一共就是满足两个条件：

1.  任务；
2.  事件循环。

我们使用 `async def func` 定义的函数叫做**协程函数**，`func()` 这样调用之后返回的结果是**协程对象**，到这一步协程函数内的代码都没有被执行，直到协程对象被包装成了任务，事件循环才会“正眼看它们”。

所以事件循环调度运行的基本单元就是**任务**，那为什么我们在使用 `async/await` 这些语句时没有涉及到任务这个概念呢？

这是因为 `await` 语法糖在内部将协程对象封装成了任务，再次强调**事件循环只认识任务**。

所以，想要运行一个协程对象，其实就是将协程对象封装成一个任务，至于事件循环是如何发现、调度和执行的，这个我们不用关心。

那将协程封装成的任务的方法有哪些呢？

*   `asyncio.create_task`
*   `asyncio.ensure_future`
*   `loop.create_task`

看着有好几个的，没关系，我们只关心 `loop.create_task`，因为其他方法最终都是调用 `loop.create_task`。

使用起来也是很简单的，将协程对象传入，返回值是一个任务对象。

```python
async def do_work():
    return 222

task = loop.create_task(do_work())
```

`do_work` 会被异步执行，那么 do_work 的结果怎么获取呢，`task.result()` 可以吗？

分情况：

*   如果是在一个协程函数内使用 `await task.result()`，这是可以的；
*   如果是在普通函数内则不行。你不可能立即获得协程函数的返回值，因为协程函数还没有被执行呢。

[asyncio.Task](https://docs.python.org/zh-cn/3/library/asyncio-task.html#asyncio.Task) 运行使用 [add_done_callback](https://docs.python.org/zh-cn/3/library/asyncio-task.html#asyncio.Task.add_done_callback) 添加完成时的回调函数，所以我们可以「曲线救国」，使用回调函数将结果添加到队列、Future 等等。

我这里给个基于 `concurrent.futures.Future` 获取结果的例子，如下：

```python
import asyncio
from asyncio import Task
from concurrent.futures import Future

from fastapi import FastAPI

app = FastAPI()
loop = asyncio.get_event_loop()

async def do_work1():
    return 222

@app.get("/")
def root():
    # 新建一个 future 对象，用于接受结果值
    future = Future()

    # 提交任务至事件循环
    task = loop.create_task(do_work1())

    # 回调函数
    def done_callback(task: Task):
        # 设置结果
        future.set_result(task.result())

    # 为这个任务添加回调函数
    task.add_done_callback(done_callback)

    # future.result 会被阻塞，直到有结果返回为止
    return future.result()  # 222
```

总结
--

*   **协程函数中调用同步函数**  
    核心思想：为了避免同步函数的执行阻塞当前线程（当前线程中的事件循环），应该由其他线程或进程执行后获取结果。
    
*   **同步函数中调用协程函数**  
    核心思想：任务协程对象都需要依赖事件循环。  
    所以要么启动新的事件循环以执行协程，要么依赖已有的事件循环, 至于实现方式就根据实际情况而选择。
    

