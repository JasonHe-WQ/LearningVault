# 搭建Ceph集群步骤

## User:
<div>现在我有3台台刚刚装好的ubuntu20.04的电脑，互相之间网络可以联通且在一个局域网之内，并且每一个都有一个未分区的硬盘。我需要在这些电脑上从0到1，搭建一个生产级别的ceph的RBD集群，需要执行哪些步骤？请尽可能地详细说明，不少于3000字</div>

## ChatGPT:
Ceph是一个分布式对象存储系统，能在商业硬件上提供出色的性能、可靠性和可扩展性。下面我将详细介绍如何从0到1在这3台Ubuntu 20.04电脑上搭建一个生产级别的Ceph RBD集群。

首先，我会详细介绍Ceph存储集群的主要组件以及它们的作用，然后我将一步步介绍如何在每个节点上安装和配置Ceph，最后我们将验证我们的集群是否正常工作。
Ceph集群的主要组件
Ceph RBD（RADOS块设备）是Ceph分布式存储系统的一部分，为虚拟机和裸机服务器提供块设备接口。它的主要组件包括：

1. Ceph OSDs（对象存储设备）：这些是存储集群的实际数据存储设备。每个硬盘或SSD都是一个OSD。

2. Ceph Monitor：收集所有OSD和元数据服务器的状态信息，主要负责维护映射信息和集群状态。

3. Ceph Manager：提供额外的监视和接口以外的服务。


下面是如何在Ubuntu 20.04上设置Ceph RBD集群的详细步骤：
1. 配置主机名和主机文件
首先，你需要在每个节点上设置一个唯一的主机名。你可以使用下面的命令设置主机名：

```
sudo hostnamectl set-hostname ceph-node1
```

将“ceph-node1”替换为你希望在每个节点上设置的主机名。在所有节点上都应使用不同的主机名。

然后，你需要在每个节点的`/etc/hosts`文件中添加所有节点的IP地址和主机名。你可以使用以下命令打开hosts文件：

```
sudo nano /etc/hosts
```

然后，添加以下行：

```
192.168.1.1 ceph-node1
192.168.1.2 ceph-node2
192.168.1.3 ceph-node3
```

请将上述IP地址替换为你的实际IP地址。
2. 更新系统
接下来，你需要在每个节点上更新你的系统。你可以使用以下命令更新系统：

```
sudo apt-get update -y
sudo apt-get upgrade -y
```

3. 安装NTP服务
Ceph集群需要时间同步。所以，你需要在每个节点上安装和配置NTP服务。你可以使用以下命令安装NTP服务：

```
sudo apt-get install chrony -y
```

然后，编辑chrony配置文件：

```
sudo nano /etc/chrony/chrony.conf
```

将以下行添加到文件的底部：

```
server ntp.ubuntu.com iburst prefer
```

然后，重新启动Chrony服务：

```
sudo systemctl restart chrony
```

4. 安装和配置Ceph
接下来，你需要在每个节点上安装Ceph。首先，安装Ceph部署工具ceph-deploy。在一个节点（我们称之为管理节点）上执行以下命令：

```
sudo apt-add-repository ppa:ceph/ceph
sudo apt-get update
sudo apt-get install ceph-deploy -y
```

然后，从管理节点创建一个新目录，用于保存配置和关键文件：

```
mkdir my-cluster
cd my-cluster
```

使用ceph-deploy工具创建一个新的Ceph集群：

```
ceph-deploy new ceph-node1 ceph-node2 ceph-node3
```

这将在当前目录下创建一个ceph.conf配置文件。

接着，编辑ceph.conf文件，添加以下内容：

```
sudo nano ceph.conf
```

在文件末尾添加以下行：

```
osd pool default size = 2
osd pool default min size = 1
osd pool default pg num = 128
osd pool default pgp num = 128
osd crush chooseleaf type = 1
radosgw url = "http://ceph-node1:8080"
```

保存并关闭文件。然后，在所有节点上安装Ceph：

```
ceph-deploy install ceph-node1 ceph-node2 ceph-node3
```

5. 配置Monitor
接下来，你需要在第一个节点上配置Ceph monitor。你可以使用以下命令在ceph-node1上创建一个monitor：

```
ceph-deploy mon create-initial
```

然后，使用以下命令将admin密钥复制到每个节点：

```
ceph-deploy admin ceph-node1 ceph-node2 ceph-node3
```

然后更改ceph.client.admin.keyring文件的权限：

```
sudo chmod 644 /etc/ceph/ceph.client.admin.keyring
```

6. 创建OSD
你可以使用以下命令在每个节点上为每个硬盘创建一个OSD：

```
ceph-deploy osd create --data /dev/sdb ceph-node1
ceph-deploy osd create --data /dev/sdb ceph-node2
ceph-deploy osd create --data /dev/sdb ceph-node3
```

将"/dev/sdb"替换为你的硬盘名称。
7. 验证集群
你可以使用以下命令检查集群的健康状况：

```
ceph health
```

如果一切顺利，你将看到HEALTH_OK的输出。

至此，一个生产级别的Ceph RBD集群就已经在你的Ubuntu 20.04电脑上搭建完成了。

## User:
<div>你在ceph.conf中修改的内容是什么意思？请详细解释。只在一个管理节点上修改后，是否会在其他节点上生效？还有其他什么字段在生产中常用吗？请详细解释。</div>

## ChatGPT:
在`ceph.conf`中，我们设定了一些参数来调整Ceph集群的行为。以下是对这些参数的详细解释：

1. osd pool default size： 这个参数设置数据的副本数量，即Ceph将在集群中的多少个OSD上存储对象的副本。在我们的例子中，我们将其设置为2，这意味着每个数据对象都将有两个副本。

2. osd pool default min size： 这个参数定义了在标记池为"active"之前需要的最小副本数。如果副本数低于这个值，Ceph将停止写入数据直到有足够的副本。

3. osd pool default pg num： 这个参数设置了Placement Group的数量，PG是Ceph将数据划分的逻辑单元。增加PG的数量可以改善数据的分布，但是如果PG的数量过多，也会增加集群的开销。

4. osd pool default pgp num： 这个参数设置了用于数据放置的Placement Group的数量。通常，这个值应该和pg num的值一样。

5. osd crush chooseleaf type： 这个参数设置了CRUSH算法的类型。CRUSH是Ceph用来决定数据应该存储在哪里的算法。在我们的例子中，我们将其设置为1，这意味着CRUSH将在主机级别选择存储位置。

6. radosgw url： 这个参数设置了RADOS Gateway的URL，RADOS Gateway是Ceph提供的对象存储接口。


这些修改只在我们执行`ceph-deploy new`命令的管理节点上生效。但是通过`ceph-deploy`工具，这些修改会被推送到整个集群，所以这些修改会在整个集群中生效。

在生产环境中，我们可能还会使用其他一些参数来优化Ceph的性能和行为，例如：

1. public network： 定义Ceph集群用于公共通信的网络。

2. cluster network： 定义Ceph集群用于集群内部通信的网络。

3. osd journal size： 设置OSD日志的大小。增加日志的大小可以提高写性能，但是也会消耗更多的磁盘空间。

4. debug ms： 设置消息模块的调试级别。在生产环境中，我们通常会关闭调试以提高性能。

5. osd max backfills： 控制一个OSD可以同时进行的最大backfill操作数量。Backfill是Ceph用来在添加或替换OSD后，将数据复制回新的OSD的过程。


以上就是`ceph.conf`中一些重要的参数，以及它们的含义和如何使用。这个文件提供了许多选项来调整Ceph集群的行为，具体需要设置哪些参数取决于你的具体需求和硬件配置。

## User:
<div>现在，我有一个已经搭建起来包含两个健康的node的K8S集群，请从0到1详细教我在K8S中配置并使用按照上述步骤搭建的ceph。请尽可能地详细描述每一个步骤，不少于3000字</div>

## ChatGPT:
要在Kubernetes集群中使用Ceph，你需要在Kubernetes集群中配置Ceph RBD作为Persistent Volume。以下是配置步骤：
1. 在Ceph节点上创建Pool和用户
首先，我们需要在Ceph集群上创建一个pool。Rados块设备（RBD）存储在RADOS池中，这些池是将物理磁盘空间组合在一起的逻辑分组。在Ceph节点上执行以下命令来创建一个名为`kubernetes`的pool：

```
ceph osd pool create kubernetes 128
```

然后，创建一个新的Ceph用户以供Kubernetes使用，并给予其相应的权限。在这个例子中，我们将创建一个名为`kubernetes`的用户：

```
ceph auth get-or-create client.kubernetes mon 'allow r' osd 'allow rwx pool=kubernetes'
```

然后，获取这个新用户的密钥：

```
ceph auth get-key client.kubernetes
```

记下这个密钥，因为稍后我们将在Kubernetes集群中使用它。
2. 在Kubernetes集群中创建Secret
Kubernetes需要Ceph用户的密钥才能访问Ceph集群。在Kubernetes集群中，我们可以使用Secret来存储这个密钥。

首先，登录到Kubernetes集群中的一个节点，然后创建一个名为`ceph-secret`的Secret：

```
kubectl create secret generic ceph-secret --type="kubernetes.io/rbd" --from-literal=key=AQD9o0Fd6hQRChAAsbYBe2gSGB3sNwvlhV+Z2g== --namespace=kube-system
```

将`AQD9o0Fd6hQRChAAsbYBe2gSGB3sNwvlhV+Z2g==`替换为你实际的Ceph密钥。
3. 创建PersistentVolume和PersistentVolumeClaim
在Kubernetes中，PersistentVolume（PV）是集群中的一块存储，它已被管理员预先配置好。PersistentVolumeClaim（PVC）是用户在PV中申请的存储空间。我们可以创建一个PV来表示我们的Ceph RBD，然后创建一个PVC来使用这个PV。

首先，创建一个名为`ceph-pv.yaml`的文件，然后添加以下内容：

```
apiVersion: v1
kind: PersistentVolume
metadata:
  name: ceph-pv
spec:
  capacity:
    storage: 1Gi
  volumeMode: Filesystem
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  storageClassName: slow
  rbd:
    monitors:
      - '192.168.1.1:6789'
      - '192.168.1.2:6789'
      - '192.168.1.3:6789'
    pool: kubernetes
    image: image-name
    user: kubernetes
    secretRef:
      name: ceph-secret
    fsType: ext4
    readOnly: false
```

将`192.168.1.1:6789`、`192.168.1.2:6789`和`192.168.1.3:6789`替换为你的Ceph monitor的IP和端口。将`image-name`替换为你的RBD image名称。

然后，创建这个PV：

```
kubectl apply -f ceph-pv.yaml
```

接下来，创建一个名为`ceph-pvc.yaml`的文件，然后添加以下内容：

```
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: ceph-pvc
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: slow
  resources:
    requests:
      storage: 1Gi
```

然后，创建这个PVC：

```
kubectl apply -f ceph-pvc.yaml
```

4. 使用Ceph RBD
现在，你可以在Pod中使用这个PVC了。创建一个名为`ceph-pod.yaml`的文件，然后添加以下内容：

```
apiVersion: v1
kind: Pod
metadata:
  name: ceph-pod
spec:
  containers:
  - name: ceph-pod
    image: nginx
    volumeMounts:
    - mountPath: "/usr/share/nginx/html"
      name: mypd
  volumes:
  - name: mypd
    persistentVolumeClaim:
      claimName: ceph-pvc
```

然后，创建这个Pod：

```
kubectl apply -f ceph-pod.yaml
```

这样，你就在Kubernetes集群中使用Ceph RBD作为存储了。你可以进入Pod检查挂载的卷，或者尝试写入一些数据，然后删除Pod，重新创建一个新的Pod并挂载同一个卷，你会发现之前写入的数据依然存在，这说明Ceph RBD已经成功作为Kubernetes的持久存储使用了。

在生产环境中，你可能还需要配置一些其他的参数，例如调整存储的大小、调整读写权限等。更多的配置选项可以在Kubernetes的官方文档中找到。

## User:
<div>请详细解释ceph-pv.yaml的参数。另外，请问monitors的端口默认情况下是6789号吗？如何修改</div>

## ChatGPT:
在`ceph-pv.yaml`中，我们定义了一个PersistentVolume（PV）对象，该对象描述了Kubernetes集群中的一个存储卷，及其相关的配置信息。以下是参数的详细解释：

- `capacity`： 定义了PV的存储容量，例如`1Gi`表示1GB。

- `volumeMode`： 定义了PV的卷模式，可以是`Filesystem`（用于存储文件的文件系统）或`Block`（裸设备，没有文件系统）。

- `accessModes`： 定义了PV的访问模式，例如`ReadWriteOnce`（单个节点读写）、`ReadOnlyMany`（多节点只读）和`ReadWriteMany`（多节点读写）。

- `persistentVolumeReclaimPolicy`： 定义了当PVC释放PV后，PV的回收策略。`Retain`表示保留卷及其数据，`Delete`表示删除卷及其数据，`Recycle`表示清空卷的数据。

- `storageClassName`： 提供了PV的类名。StorageClass是用于描述PV的一种方式，可以用于动态创建PV。

- `rbd`： 定义了RBD的相关配置。RBD是Ceph的一种块设备，可以被挂载到Kubernetes的Pod上。

`monitors`： Ceph集群的monitor节点的地址列表，可以包含IP地址和端口。

`pool`： 在Ceph集群中使用的pool的名称。

`image`： RBD的image名称。

`user`： 用于访问Ceph集群的用户名称。

`secretRef`： 包含Ceph用户密钥的Kubernetes Secret的名称。

`fsType`： 在卷上使用的文件系统类型。

`readOnly`： 是否以只读模式挂载卷。

	- `monitors`： Ceph集群的monitor节点的地址列表，可以包含IP地址和端口。

	- `pool`： 在Ceph集群中使用的pool的名称。

	- `image`： RBD的image名称。

	- `user`： 用于访问Ceph集群的用户名称。

	- `secretRef`： 包含Ceph用户密钥的Kubernetes Secret的名称。

	- `fsType`： 在卷上使用的文件系统类型。

	- `readOnly`： 是否以只读模式挂载卷。


在默认情况下，Ceph monitor的端口确实是6789。但如果你想修改这个端口，你需要在Ceph的配置文件`ceph.conf`中修改`mon_port`参数，然后重启Ceph monitor。

例如，如果你想将端口修改为1234，你可以在`ceph.conf`文件中添加或修改以下行：

```
[mon]
mon_port = 1234
```

然后，在每个monitor节点上重启Ceph monitor：

```
sudo systemctl restart ceph-mon.target
```

修改后，你也需要在Kubernetes的PV配置中更新monitor的端口。

## User:
<div>请问RBD的image名称是什么意思？请详细解释并举例分析</div>

## ChatGPT:
Ceph的RADOS块设备（RBD）是一个分布式块设备，它将数据存储在Ceph的对象存储池中。每个RBD都对应一个image（镜像）。这个image并不是Docker那样的容器镜像，而是类似于一个硬盘的镜像。你可以将它理解为一块虚拟的硬盘，Kubernetes或其他系统可以将其挂载，并像使用物理硬盘一样使用它。

每个RBD image都有一个唯一的名称，我们在创建RBD时指定这个名称。例如，我们可以使用以下命令在名为`kubernetes`的pool中创建一个名为`my-image`的RBD image：

```
rbd create kubernetes/my-image --size 1024
```

这个命令会在`kubernetes`这个pool中创建一个名为`my-image`，大小为1GB的RBD image。然后，我们可以在Kubernetes的PV配置中指定`image: my-image`，来让PV使用这个RBD image。

使用RBD的好处是，我们可以在多个节点上同时挂载同一个RBD，实现数据的共享。同时，因为RBD是分布式的，所以它的容量可以非常大，也能提供很高的I/O性能。

## User:
<div>那么你在上面的教程中没有涉及到这一点。请问是默认有什么配置还是你犯了错误？</div>

## ChatGPT:
对不起，这确实是我的疏忽。在创建RBD对应的Kubernetes PV之前，我们应该在Ceph集群中创建对应的RBD image。

请按照以下步骤创建RBD image：

在Ceph集群中的一个节点上执行以下命令：

```
rbd create kubernetes/my-image --size 1024 --pool=kubernetes --id kubernetes --keyring=/etc/ceph/keyring
```

在这条命令中：

- `kubernetes/my-image`是新创建的RBD image的名称，`kubernetes`是pool的名称，`my-image`是image的名称。

- `--size 1024`指定了新创建的RBD image的大小为1024MB，也就是1GB。你可以根据需要调整这个大小。

- `--pool=kubernetes`指定了存储RBD image的pool的名称。

- `--id kubernetes`指定了访问Ceph集群的用户的ID。

- `--keyring=/etc/ceph/keyring`指定了包含用户密钥的文件的路径。


然后，你就可以在Kubernetes的PV配置中使用这个新创建的RBD image了。在PV配置的`rbd`字段中，将`image`的值改为你新创建的image的名称，例如`my-image`。

再次对之前的疏忽表示歉意，希望这个解答能帮到你。

## User:
<div>请再次检查上述教程中有没有错误或遗漏。如果有，请详细说明</div>

## ChatGPT:
我再次检查了上述的教程，我发现了一处需要改正的地方，那就是在创建PV的时候使用了已被弃用的RBD的volume插件，应该使用Ceph CSI插件。对此我表示歉意。

Ceph CSI插件为Kubernetes提供了对Ceph的支持，它实现了Container Storage Interface（CSI），使得Kubernetes可以使用Ceph作为存储。Ceph CSI插件支持RBD和CephFS。

以下是安装和配置Ceph CSI插件的步骤：

1. 安装Ceph CSI插件：

首先，你需要在Kubernetes集群中部署Ceph CSI插件。你可以在Ceph CSI的GitHub页面找到部署的YAML文件。根据你的Kubernetes版本和Ceph版本，选择对应的YAML文件，然后使用kubectl apply命令进行部署。

2. 创建CephCSIStorageClass：

创建一个名为`csi-rbd-sc.yaml`的文件，并添加以下内容：

```
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: csi-rbd-sc
provisioner: rbd.csi.ceph.com
parameters:
  clusterID: your-cluster-id
  pool: kubernetes
  csi.storage.k8s.io/provisioner-secret-name: ceph-secret
  csi.storage.k8s.io/provisioner-secret-namespace: default
  csi.storage.k8s.io/node-stage-secret-name: ceph-secret
  csi.storage.k8s.io/node-stage-secret-namespace: default
  imageFeatures: layering
  csi.storage.k8s.io/fstype: ext4
```

在这个文件中，将`your-cluster-id`替换为你的Ceph集群的ID，你可以在Ceph的配置文件或者通过`ceph fsid`命令获取到。然后，使用以下命令创建StorageClass：

```
kubectl apply -f csi-rbd-sc.yaml
```

3. 创建PVC：

创建一个名为`ceph-pvc.yaml`的文件，然后添加以下内容：

```
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: ceph-pvc
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
  storageClassName: csi-rbd-sc
```

然后，使用以下命令创建PVC：

```
kubectl apply -f ceph-pvc.yaml
```

4. 在Pod中使用PVC：

创建一个名为`ceph-pod.yaml`的文件，然后添加以下内容：

```
apiVersion: v1
kind: Pod
metadata:
  name: ceph-pod
spec:
  containers:
  - name: ceph-pod
    image: nginx
    volumeMounts:
    - mountPath: "/usr/share/nginx/html"
      name: mypd
  volumes:
  - name: mypd
    persistentVolumeClaim:
      claimName: ceph-pvc
```

然后，使用以下命令创建Pod：

```
kubectl apply -f ceph-pod.yaml
```


这样，你就可以在Kubernetes中使用Ceph RBD作为存储了。如果你需要更详细的步骤或者遇到了问题，你可以查看Ceph CSI插件的文档。

## User:
<div>csi.storage.k8s.io/provisioner-secret-name是什么？csi.storage.k8s.io/node-stage-secret-name是什么？他们有区别吗？还是都对应ceph的密码？请详细解答。另外，请问我是否还需要安装provisioner还是Ceph CSI插件已经包括了这个？</div>

## ChatGPT:
`csi.storage.k8s.io/provisioner-secret-name`和`csi.storage.k8s.io/node-stage-secret-name`是两个Ceph CSI插件用于获取Ceph集群密钥的参数。他们的值都应该是一个Kubernetes Secret的名称，这个Secret中应该包含了Ceph集群的访问密钥。

- `csi.storage.k8s.io/provisioner-secret-name`： 这个参数指定了在创建或删除存储卷（也就是执行provision操作）时，Ceph CSI插件如何获取Ceph集群的访问密钥。Provisioner运行在Kubernetes的master节点上，用于响应用户创建或删除PVC的请求。

- `csi.storage.k8s.io/node-stage-secret-name`： 这个参数指定了在将存储卷挂载到Pod（也就是执行node stage操作）时，Ceph CSI插件如何获取Ceph集群的访问密钥。Node Plugin运行在Kubernetes的工作节点上，用于将PVC挂载到Pod。


通常情况下，我们会创建一个包含了Ceph集群访问密钥的Kubernetes Secret，并将这个Secret的名称设置为`csi.storage.k8s.io/provisioner-secret-name`和`csi.storage.k8s.io/node-stage-secret-name`的值。这样，Ceph CSI插件在执行provision操作和node stage操作时，都可以使用这个Secret获取到Ceph集群的访问密钥。

至于你的第二个问题，Ceph CSI插件已经包含了Provisioner，所以你不需要单独安装Provisioner。Ceph CSI插件会根据StorageClass的配置自动创建和删除存储卷，你只需要在PVC中指定使用哪个StorageClass即可。

## User:
<div>请重新给出创建一个包含了Ceph集群访问密钥的Kubernetes Secret的示例，并检查上述教程是否还有其他错误或遗漏</div>

## ChatGPT:
对不起之前的解答中确实漏掉了如何创建Ceph访问密钥的Secret。以下是创建该Secret的步骤。

首先，你需要获取Ceph的访问密钥。你可以在Ceph集群的管理节点上，通过以下命令获取：

```
ceph auth get-key client.admin
```

这个命令会输出一个长字符串，这就是Ceph的访问密钥。

然后，创建一个名为`ceph-secret.yaml`的文件，并添加以下内容：

```
apiVersion: v1
kind: Secret
metadata:
  name: ceph-secret
data:
  userID: YWRtaW4= # "admin"的Base64编码
  userKey: YOUR_CEPH_ACCESS_KEY # Ceph访问密钥的Base64编码
```

在这个文件中，将`YOUR_CEPH_ACCESS_KEY`替换为你的Ceph访问密钥的Base64编码。你可以通过`echo -n YOUR_CEPH_ACCESS_KEY | base64`命令获取这个编码。

然后，使用以下命令创建Secret：

```
kubectl apply -f ceph-secret.yaml
```

这样，你就创建了一个包含了Ceph访问密钥的Secret。你可以在StorageClass的配置中，将`csi.storage.k8s.io/provisioner-secret-name`和`csi.storage.k8s.io/node-stage-secret-name`的值都设置为`ceph-secret`，这样Ceph CSI插件就可以使用这个Secret来获取Ceph的访问密钥了。

在进行一次详细的回顾和检查后，我认为上述教程现在应该是正确的，已经包括了所有需要的步骤。不过，由于具体环境可能存在差异，实际操作中可能还需要根据具体情况进行一些调整。再次对之前的疏忽表示歉意，并感谢你的耐心和理解。

Thread generated with <a href="https://chat.openai.com">ChatGPT</a> and downloaded with <a href="https://botrush.io">Botrush</a>