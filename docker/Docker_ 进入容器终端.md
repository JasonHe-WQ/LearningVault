# Docker: 进入容器终端



## USER:

docker 中如何为每个容器的 cpu/memory 设限，原理是什么

## CHATGPT:

在 Docker 中，可以使用一些参数为容器设置 CPU 和内存限制，从而确保单个容器不会消耗过多的系统资源。以下是如何为 Docker 容器设置 CPU 和内存限制的方法：

- CPU 限制：
- `--cpus`：设置容器可以使用的 CPU 核心数。例如，使用 `--cpus=0.5` 表示容器最多可以使用 50% 的一个 CPU 核心。
- `--cpu-shares`：设置容器的 CPU 份额。默认情况下，所有容器的 CPU 份额相等。通过设置不同的份额，可以调整容器之间的 CPU 资源分配比例。例如，如果容器 A 的份额为 1024，而容器 B 的份额为 512，则容器 A 将获得两倍于容器 B 的 CPU 时间。
- 内存限制：
- `--memory`：设置容器的最大内存使用量。当容器的内存使用量超过这个限制时，Docker 会尝试回收内存或终止容器。例如，使用 `--memory=512m` 表示容器的最大内存使用量为 512 MB。
- `--memory-swap`：设置容器的最大交换空间使用量。交换空间是当物理内存不足时，系统将内存中的一部分数据暂存到磁盘上的空间。例如，使用 `--memory-swap=1g` 表示容器的最大交换空间使用量为 1 GB。
- `--memory-reservation`：设置容器的内存保留量。这是一个软限制，表示容器在内存资源紧张时可以保留的最小内存量。例如，使用 `--memory-reservation=256m` 表示容器的内存保留量为 256 MB。
以下是一个使用这些参数的示例：

```
docker run -it --cpus=0.5 --cpu-shares=1024 --memory=512m --memory-swap=1g --memory-reservation=256m ubuntu:18.04 /bin/bash

```

Docker 通过内核的 cgroups（控制组）功能实现对容器的资源限制。cgroups 可以将一组进程组织到一个层次结构中，并为这些进程设置资源限制，如 CPU、内存、磁盘 I/O 等。Docker 为每个容器创建一个 cgroup，并将这些参数应用到相应的 cgroup 中。当容器运行时，内核会根据 cgroup 中的设置来管理和限制容器的资源使用。

需要注意的是，这些限制只能确保容器在某个水平上共享资源，但不能保证容器在高负载时始终获得足够的资源。为了实现更精确的资源保障，可以考虑使用其他容器编排工具，如 Kubernetes。

## USER:

在 docker 的容器中，如何访问宿主机的 localhost

## CHATGPT:

在 Docker 容器中访问宿主机的 localhost，需要使用宿主机的特殊 IP 地址。根据 Docker 容器所运行的操作系统，可以采用以下方法：

1. 对于 Linux 宿主机：

在容器中，可以使用宿主机的 IP 地址 `172.17.0.1`（默认情况下）访问宿主机的 localhost。这是 Docker 默认网桥的 IP 地址。你可以通过以下命令在宿主机上查看默认网桥的详细信息：

```bash
ip addr show docker0
```
当然，这个 IP 地址可能会发生变化。为了确保正确获取宿主机的 IP 地址，你可以使用以下命令：

```bash
ip -4 route show default | awk '{print $3}'
```

2. 对于 Windows 或 macOS 宿主机：

Docker 在 Windows 和 macOS 上运行在一个轻量级的虚拟机中，而不是直接运行在宿主机上。在这种情况下，可以使用虚拟机的网关 IP 地址访问宿主机的 localhost。通常，这个 IP 地址为 `host.docker.internal`（Docker 18.03 及更高版本）。

例如，如果宿主机上运行了一个 Web 服务，监听在端口 8000 上，那么在容器中可以使用以下命令访问该服务：
```bash
curl http://host.docker.internal:8000
```


需要注意的是，这种方法只适用于访问宿主机上运行的服务。如果需要访问容器之间的服务，可以使用 Docker 网络功能，如创建自定义网络或使用默认的 Docker 网络。


