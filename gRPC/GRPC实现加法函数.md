# GRPC实现加法函数

## User:
请写一个python程序，要求使用gRPC，实现add函数。这个函数能够将两个整数相加并返回他们的和

## ChatGPT:
为了使用gRPC在Python中实现add函数，你需要进行几个步骤。

步骤一，先定义gRPC服务。gRPC使用Protocol Buffers作为其接口定义语言，因此，你需要创建一个`.proto`文件来定义你的服务。

```
syntax = "proto3";

package calculator;

// The calculator service definition.
service Calculator {
  // Sends two numbers and returns their sum.
  rpc Add (AddRequest) returns (AddReply) {}
}

// The request message containing two numbers.
message AddRequest {
  int32 num1 = 1;
  int32 num2 = 2;
}

// The response message containing their sum.
message AddReply {
  int32 sum = 1;
}
```

保存这个文件为`calculator.proto`。

步骤二，使用gRPC Python工具生成客户端和服务器端代码。执行以下命令：

```
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. calculator.proto
```

这将生成`calculator_pb2.py`和`calculator_pb2_grpc.py`两个文件。

步骤三，实现服务器端代码：

```
from concurrent import futures
import grpc
import calculator_pb2
import calculator_pb2_grpc

class Calculator(calculator_pb2_grpc.CalculatorServicer):

    def Add(self, request, context):
        return calculator_pb2.AddReply(sum=request.num1 + request.num2)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    calculator_pb2_grpc.add_CalculatorServicer_to_server(Calculator(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()
```

将以上代码保存为`calculator_server.py`。

步骤四，实现客户端代码：

```
import grpc
import calculator_pb2
import calculator_pb2_grpc

def run():
    channel = grpc.insecure_channel('localhost:50051')
    stub = calculator_pb2_grpc.CalculatorStub(channel)
    response = stub.Add(calculator_pb2.AddRequest(num1=3, num2=5))
    print("3 + 5 = " + str(response.sum))

if __name__ == '__main__':
    run()
```

将以上代码保存为`calculator_client.py`。

最后，先运行服务器端代码（`python calculator_server.py`），然后运行客户端代码（`python calculator_client.py`）。这样你就能看到客户端输出两个数的和了。

## User:
那么，请详细解释python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. calculator.proto

## ChatGPT:
这个命令是用来生成Python版本的gRPC代码的。

命令中的各个部分解释如下：

1. `python -m grpc_tools.protoc`：这个命令使用Python运行gRPC的Protocol Buffer编译器。这是因为gRPC使用Protocol Buffers作为其接口定义语言。

2. `-I.`：这个参数是用来指定查找import文件的目录，`.`表示当前目录。

3. `--python_out=.`：这个参数用来指定生成的基础Python代码（定义了所有消息类）的输出目录，`.`表示当前目录。生成的文件名会是`<你的.proto文件名>_pb2.py`，在这个例子中就是`calculator_pb2.py`。

4. `--grpc_python_out=.`：这个参数用来指定生成的gRPC特定Python代码（定义了服务器和客户端类）的输出目录，`.`表示当前目录。生成的文件名会是`<你的.proto文件名>_pb2_grpc.py`，在这个例子中就是`calculator_pb2_grpc.py`。

5. `calculator.proto`：这是你要编译的Protocol Buffer定义文件。


这个命令将会读取`calculator.proto`文件，然后生成两个Python文件：`calculator_pb2.py`和`calculator_pb2_grpc.py`，这两个文件中包含了你在`.proto`文件中定义的服务和消息的Python版本。

## User:
请问能不能将上述的代码更改为async异步调用以提高性能？

## ChatGPT:
当然可以。在Python的gRPC实现中，你可以使用异步IO (AIO) 模块来实现异步调用。以下是相应的服务器端和客户端代码。

首先是服务器端代码：

```
import asyncio
import grpc
import calculator_pb2
import calculator_pb2_grpc

class Calculator(calculator_pb2_grpc.CalculatorServicer):

    async def Add(self, request, context):
        return calculator_pb2.AddReply(sum=request.num1 + request.num2)

async def serve():
    server = grpc.aio.server()
    calculator_pb2_grpc.add_CalculatorServicer_to_server(Calculator(), server)
    server.add_insecure_port('[::]:50051')
    await server.start()
    await server.wait_for_termination()

if __name__ == '__main__':
    asyncio.run(serve())
```

将以上代码保存为`calculator_server.py`。

然后是客户端代码：

```
import asyncio
import grpc
import calculator_pb2
import calculator_pb2_grpc

async def run():
    async with grpc.aio.insecure_channel('localhost:50051') as channel:
        stub = calculator_pb2_grpc.CalculatorStub(channel)
        response = await stub.Add(calculator_pb2.AddRequest(num1=3, num2=5))
        print("3 + 5 = " + str(response.sum))

if __name__ == '__main__':
    asyncio.run(run())
```

将以上代码保存为`calculator_client.py`。

然后同样地，先运行服务器端代码（`python calculator_server.py`），然后运行客户端代码（`python calculator_client.py`）。这样你就能看到客户端输出两个数的和了。

请注意，你需要Python 3.7或更高版本才能运行这些代码，因为`asyncio.run()`函数是在Python 3.7中添加的。同样的，gRPC Python AIO API也需要Python 3.6或更高版本。

## User:
那么，如何关闭这个RPC调用呢？如何设置超时机制？

## ChatGPT:
在gRPC中，你可以使用超时参数或者取消操作来结束一个RPC调用。

1. 设置超时：你可以在调用RPC方法时传入一个超时参数（单位为秒）。如果在超时时间内服务还未响应，gRPC将会中断该请求。下面是一个客户端代码示例，展示如何为请求设置5秒的超时时间：

```
import asyncio
import grpc
import calculator_pb2
import calculator_pb2_grpc

async def run():
    async with grpc.aio.insecure_channel('localhost:50051') as channel:
        stub = calculator_pb2_grpc.CalculatorStub(channel)
        try:
            response = await stub.Add(calculator_pb2.AddRequest(num1=3, num2=5), timeout=5.0)
            print("3 + 5 = " + str(response.sum))
        except grpc.RpcError as e:
            print(f"RPC failed: {e}")

if __name__ == '__main__':
    asyncio.run(run())
```

2. 取消操作：你也可以主动取消一个操作。在异步gRPC中，你可以使用`asyncio.CancelledError`异常来捕获取消的操作。以下是客户端代码的一个示例，展示了如何在RPC调用开始后1秒取消操作：

```
import asyncio
import grpc
import calculator_pb2
import calculator_pb2_grpc

async def run():
    async with grpc.aio.insecure_channel('localhost:50051') as channel:
        stub = calculator_pb2_grpc.CalculatorStub(channel)
        call = stub.Add(calculator_pb2.AddRequest(num1=3, num2=5))

        await asyncio.sleep(1.0)  # 等待一秒
        call.cancel()  # 取消RPC调用

        try:
            response = await call
            print("3 + 5 = " + str(response.sum))
        except asyncio.CancelledError:
            print("The call was cancelled.")

if __name__ == '__main__':
    asyncio.run(run())
```


在这个例子中，客户端将会在调用开始后的一秒内取消操作。注意，如果服务在这一秒内已经响应了，取消操作将不会生效。

Thread generated with <a href="https://chat.openai.com">ChatGPT</a> and downloaded with <a href="https://botrush.io">Botrush</a>