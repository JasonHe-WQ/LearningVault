# 本地构建K8S镜像库

## User:
请问如何构建本地的用于K8S的镜像库？

## ChatGPT:
要在本地构建一个用于 Kubernetes（K8S）的镜像库，你可以使用 Docker Registry，这是一个非常流行的选择。下面是具体步骤：

1. 安装 Docker：你首先需要在机器上安装 Docker。在大多数 Linux 发行版中，你可以使用包管理器来安装 Docker。对于 Ubuntu，你可以使用如下命令：

```
sudo apt update
sudo apt install docker.io
```

1. 运行 Docker Registry：你可以运行 Docker Registry 作为一个容器。下面是如何使用 Docker 命令来运行 Docker Registry 的命令：

```
docker run -d -p 5000:5000 --name=registry --restart=always registry:2
```

这个命令将会下载并运行 Docker Registry，它将监听 5000 端口。

1. 推送镜像到你的本地 Docker Registry：现在你可以将 Docker 镜像推送到你的本地 Docker Registry 了。首先，你需要标记你的 Docker 镜像，然后你可以推送它。下面是如何做到这一点的：

```
docker tag my-image localhost:5000/my-image
docker push localhost:5000/my-image
```

1. 在 Kubernetes 中使用你的 Docker 镜像：现在你可以在 Kubernetes 中使用你的 Docker 镜像了。在你的 Kubernetes 配置文件中，你可以使用如下的镜像路径：

```
image: localhost:5000/my-image
```

注意，由于默认情况下 Docker 不允许非安全的注册表，因此你可能需要配置你的 Docker 守护程序来允许这种连接。你可以在 Docker 守护程序的配置文件（通常在 `/etc/docker/daemon.json`）中添加如下的内容：

```
{
  "insecure-registries" : ["myregistrydomain.com:5000"]
}
```

然后重启 Docker 守护程序：

```
sudo service docker restart
```

以上就是在本地构建用于 Kubernetes 的 Docker 镜像库的步骤。

