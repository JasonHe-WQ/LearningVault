[生产环境Go程序内存泄露，用pprof如何快速定位](https://mp.weixin.qq.com/s?__biz=MzkyNzI1NzM5NQ==&mid=2247486203&idx=1&sn=b75d45f23bc63b514a3ebf0788b3a4c1&sessionid=1699288514#rd)
===============================================================================================================================================================

Golang梦工厂

内存泄漏可以在整个系统中以多种形式出现，除了在写代码上的疏忽，忘了关闭该关闭的资源外，更多的时候导致系统发生内存泄露原因可能是设计上决策不对、或者业务逻辑上的疏忽没有考虑到一些边界条件。

比如查数据库时，有个查询条件在一定情况下应用不到，导致程序被迫持有一个超大的结果集，这样持续一段时间，执行相同任务的线程一多，就会造成内存泄露。

Golang 为我们提供了 pprof 工具。掌握之后，可以帮助排查程序的内存泄露问题，当然除了排查内存，它也能排查 CPU 占用过高，线程死锁的这些问题，不过这篇文章我们会聚焦在怎么用 pprof 排查程序的内存泄露问题。

Go 开发的系统中，怎么 添加 pprof 进行采样的步骤，在这里我就不再细说了，因为我之前的文章，对 pprof 的安装和使用做了详细的说明，文章链接我放在这里：

*   [Golang 程序性能分析（一）pprof 和 火焰图分析](https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247486618&idx=1&sn=bb5e76e011ba99ebc2ffb8f9d3c00b89&scene=21#wechat_redirect)
    
*   [Golang程序性能分析（二）在Echo和Gin框架中使用pprof](https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247486654&idx=1&sn=ea7171f58254dfecebc61cfbec7b64e5&scene=21#wechat_redirect)
    
*   [Golang程序性能分析（三）用pprof分析gRPC服务的性能](https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247486669&idx=1&sn=2be1d2cb1b85c8d0f59a314cafb4637f&scene=21#wechat_redirect)
    

当然如果你想尝试点更智能的，让程序能自己监控自己，并在出现抖动的时候自己采样，Dump 出导致内存、CPU的问题调用栈信息，可以看一下下面两篇文章里我介绍的方法和实用的类库。

*   [学会这几招让 Go 程序自己监控自己](https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247490745&idx=1&sn=6a04327f98a734fd50e509362fc04d48&scene=21#wechat_redirect)
    
*   [Go 服务进行自动采样性能分析的方案设计与实现](https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247491713&idx=1&sn=3735d8f028823eaca6f42b28a9e7d817&scene=21#wechat_redirect)
    

内存泄露该看哪个指标
----------

`pprof`工具集，提供了`Go`程序内部多种性能指标的采样能力，我们常会用到的性能采样指标有这些：

*   profile：CPU采样
    
*   heap：堆中活跃对象的内存分配情况的采样
    
*   goroutine：当前所有goroutine的堆栈信息
    
*   allocs: 会采样自程序启动所有对象的内存分配信息（包括已经被GC回收的内存）
    
*   threadcreate：采样导致创建新系统线程的堆栈信息
    

上面 heap 和 allocs 是两个与内存相关的指标， allocs 指标会采样自程序启动所有对象的内存分配信息。一般是在想要分析哪些代码能优化提高效率时，查看的指标。针对查看内存泄露问题的分析，使用则的是 heap 指标里的采样信息。

Heap
----

pprof 的 heap 信息，是对堆中活跃对象的内存分配情况的采样。Go 里边哪些对象会被分配到堆上？一般概况就是，被多个函数引用的对象、全局变量、超过一定体积（32KB）的对象都会被分配到堆上，当然对于 Go 来说还会有其他的一些情况会让对象逃逸到堆上。

具体哪些变量会被分配到堆上、以及内存逃逸的事儿，就不多说了，想看详细情况的，看下面这两篇文章。

*   [图解Go内存管理器的内存分配策略](https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247485736&idx=1&sn=921a9dfe3d638074b68a4fd072ea3cb9&scene=21#wechat_redirect)
    
*   [Go内存管理之代码的逃逸分析](https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247485579&idx=1&sn=f481cff4ffccacc186a020e45e884924&scene=21#wechat_redirect)
    

Heap 采样
-------

要使用 pprof 获取 heap 指标的采样信息，一种情况是使用 "net/http/pprof" 包

```go
import (
 "net/http/pprof"
)

func main() {
 http.HandleFunc("/debug/pprof/heap", pprof.Index)
 ......
 http.ListenAndServe(":80", nil)
}
```

然后通过 HTTP 请求的方式获得

```go
curl -sK -v https://example.com/debug/pprof/profile > heap.out
```

还有一种主要的方法是使用`runtime.pprof`  提供的方法，把采样信息保存到文件。

```go
pprof.Lookup("heap").WriteTo(profile_file, 0)
```

关于这两个包的使用方式，以及怎么把信息采样到文件，上面介绍自动采样的文章里有详细的介绍，这里就不再花过多篇幅了。

下面进入文章的正题， 拿到采样文件后，怎么用 pprof 排查出代码哪里导致了内存泄露。

用 pprof 找出内存泄露的地方
-----------------

pprof 在采样 heap 指标的信息时，使用的是 runtime.MemProfile 函数，该函数默认收集每个 512KB 已分配字节的分配信息。我们可以设置让 runtime.MemProfile 收集所有对象的信息，不过这会对程序的性能造成影响。

当我们拿到采样文件后，就可以通过 go tool pprof 将信息加载到一个交互模式的控制台中。

```go
> go tool pprof heap.out
```

进入，交互式控制台后，一般会有如下的提示：

```go
File: heap.out
Type: inuse_space
Time: Feb 1, 2022 at 10:11am (CST)
Entering interactive mode (type "help" for commands, "o" for options)
```

这里的 `Type: inuse_space`  指明了文件内采样信息的类型， Type 可能的值有：

*   inuse_space — 已分配但尚未释放的内存空间
    
*   inuse_objects——已分配但尚未释放的对象数量
    
*   alloc_space — 分配的内存总量（已释放的也会统计）
    
*   alloc_objects — 分配的对象总数（无论是否释放）
    

接下来，介绍一个 pprof 交互式模式下的命令 **top**，也可以是 **topN**，比如 top10。这个跟Linux 系统的 top 命令类似，输出 Top N 个最占用内存的函数。

```go
(pprof) top10
Showing nodes accounting for 134.55MB, 92.16% of 145.99MB total
Dropped 60 nodes (cum <= 0.73MB)
Showing top 10 nodes out of 117
      flat  flat%   sum%        cum   cum%
   60.53MB 41.46% 41.46%    85.68MB 58.69%  github.com/jinzhu/gorm.glob..func2
   18.65MB 12.77% 54.24%    18.65MB 12.77%  regexp.(*Regexp).Split
   16.95MB 11.61% 65.84%    16.95MB 11.61%  github.com/jinzhu/gorm.(*Scope).AddToVars
    8.67MB  5.94% 71.78%   129.05MB 88.39%  example.com/xxservice/dummy.GetLargeData
    7.50MB  5.14% 82.63%     7.50MB  5.14%  reflect.packEface
    6.50MB  4.45% 87.08%     6.50MB  4.45%  fmt.Sprintf
       4MB  2.74% 89.82%        4MB  2.74%  runtime.malg
    1.91MB  1.31% 91.13%     1.91MB  1.31%  strings.Replace
    1.51MB  1.03% 92.16%     1.51MB  1.03%  bytes.makeSlice
```

在这两个里边，最占用内存的前三是 gorm 库的一个方法，gorm 是个 ORM 库，但是导致它内存泄露的原因应该是后面一个有业务逻辑的代码，dummy.GetLargeData 方法。

在 top 指令输出的列表中，我们可以看到两个值，flat 和 cum。

*   flat：表示此函数分配、并由该函数持有的内存空间。
    
*   cum：表示由这个函数或它调用堆栈下面的函数分配的内存总量。
    

此外 sum % 表示前面几行输出的 flat百分比之和， 比如上面第四行 sum% 列的值是， 71.78%  实际上就是它以及它上面三行输出的 flat% 的总和。

定位到导致内存泄露的函数后，后面要做的优化问题就是，深入函数内部，看哪里使用不当或者有逻辑上的疏忽，比如我开头举得那个查询条件在有些情况下应用不上的例子。

当然如果你想在函数内部再精确的定位到底是哪段代码导致的内存溢出，也是有办法的，这时候就需要用到 list 指令了。

**list** 指令可以列出函数内部，每一行代码运行时分配的内存（如果分析CPU的采样文件，则会显示CPU使用时间）

```go
(pprof) list dummy.GetLargeData
Total: 814.62MB
ROUTINE ======================== dummy.GetLargeData in /home/xxx/xxx/xxx.go
  814.62MB   814.62MB (flat, cum)   100% of Total
         .          .     20:    }()
         .          .     21:
         .          .     22:    tick := time.Tick(time.Second / 100)
         .          .     23:    var buf []byte
         .          .     24:    for range tick {
  814.62MB   814.62MB     25:        buf = append(buf, make([]byte, 1024*1024)...)
         .          .     26:    }
         .          .     27:}
         .          .     28:
```

总结
--

这里把用 pprof 怎么排查程序的内存泄露做了个简单的总结，当然如果你们公司有条件上持续采样，或者我之前文章说的自动采样方案的话，最好还是用上，让机器帮我们做这些事情。

不过不管是用什么办法，最终只能是帮我们定位出来哪里造成了内存泄露，至于要怎么优化解决这个问题，还得具体情况具体分析，如果是一些业务逻辑实现上的问题，那就得跟团队商量一下实现方式，可能还会涉及到产品上的一些改动。

- END -

扫码关注公众号「网管叨bi叨」

![图片](https://mmbiz.qpic.cn/mmbiz_png/z4pQ0O5h0f4icJbGAQ8RjXUUVdUZsGADuMBVWePgn7tfrWjjHyc6b8kXTQ7Sdkp0QQFFK4mel5tniczqooMna1CA/640?wx_fmt=png)

给网管个星标，第一时间吸我的知识 👆

网管为大家整理了一本超实用的《Go 开发参考书》收集了70多条开发实践。去公众号回复【gocookbook】即刻领取！

觉得有用就点个在看👇👇👇
