[\[K8S\] client-go 中 informer 的一些技巧](https://juejin.cn/post/7218120593623236667)
================================================================================

ShadowYD

前言
--

本次介绍一下 **informer** 实际使用中遇到的问题与一些技巧进行补充, 如需了解 **Informer** 比较详细的使用方法请查阅 [\[K8S\] client-go 的正确打开方式](https://juejin.cn/post/7203690731276517432 "https://juejin.cn/post/7203690731276517432").

资源对象会丢失的数据 Kind 和 Version
-------------------------

*   如果你在生产环境中经常使用 **Informer**，那么你可能已经遇到了从 Informer 中获取的资源对象不包含`Kind` 和`ApiVersion`这两个字段的情况。此时，你可能会怀疑是自己的用法出了问题导致出现这个问题。实际上，并不是你的用法有误，而是 Informer 的一种兼容性设计问题。
    
*   为了恢复这两个字段，我们可以使用客户端中的注册表(**Scheme**)。注册表兼容了当前client-go版本之前的所有API版本，因此可以用于获取Kind和ApiVersion信息。
    
*   示例代码:
    
    ```
    import (
    "k8s.io/client-go/kubernetes/scheme"
    )
    
    ...
    // GVKForObject 传入k8s资源对象, 并且传入一个k8s的注册表
    // 如果是自定义资源则需要用该资源的注册表, 一般都是由k8s codegen 自动生成;
    func GVKForObject(obj runtime.Object, scheme \*runtime.Scheme) (schema.GroupVersionKind, error) {
            \_, isPartial := obj.(\*metav1.PartialObjectMetadata) //nolint:ifshort
            \_, isPartialList := obj.(\*metav1.PartialObjectMetadataList)
            // 如果存在部分metadata 信息会尝试通过, runtime.Object 自带的接口去获取GVK
            // 我们都需要恢复大概率就肯定是因为我们没有 😄
            if isPartial || isPartialList {
    
                    gvk := obj.GetObjectKind().GroupVersionKind()
                    if len(gvk.Kind) == 0 {
                            return schema.GroupVersionKind{}, runtime.NewMissingKindErr("unstructured object has no kind")
                    }
                    if len(gvk.Version) == 0 {
                            return schema.GroupVersionKind{}, runtime.NewMissingVersionErr("unstructured object has no version")
                    }
                    return gvk, nil
            }
            // 我们会尝试从注册表中去解释 GVK 列表
            gvks, isUnversioned, err := scheme.ObjectKinds(obj)
            if err != nil {
                    return schema.GroupVersionKind{}, err
            }
            if isUnversioned {
                    return schema.GroupVersionKind{}, fmt.Errorf("cannot create group-version-kind for unversioned type %T", obj)
            }
    
            if len(gvks) < 1 {
                    return schema.GroupVersionKind{}, fmt.Errorf("no group-version-kinds associated with type %T", obj)
            }
            if len(gvks) > 1 {
                    return schema.GroupVersionKind{}, fmt.Errorf(
                            "multiple group-version-kinds associated with type %T, refusing to guess at one", obj)
            }
            return gvks\[0\], nil
    }
    
    //RestoreGVKForList 恢复资源对象列表, 
    func RestoreGVKForList\[T ObjectType\](objects \[\]T) error {
        if len(objects) != 0 {
                gvk, err := GVKForObject(objects\[0\], scheme.Scheme)
                if err != nil {
                        return err
                }
                for \_, obj := range objects {
                        // 把获取到的 GVK 设置回 k8s 资源
                        obj.SetGroupVersionKind(gvk)
                }
        }
        return nil
    }
    
    
    
    ...
    ```
    
*   如何进行测试, 只需要构造一个k8s资源对象, 或者通过 Informer 获取一个任意资源也是可以的;
    
    ```
    func Test\_RestoreGVKForList(t \*test.T) {
         podObject := corev1.Pod{}
         if err := RestoreGVKForList\[corev1.Pod\](\[\]corev1.Pod{podObject}); err != nil {
             panic(err)
         }
         fmt.Println(podObject)
    
    }
    ```
    

如何重启 Informer 时清理缓存数据
---------------------

*   实际上，**Informer** 就是在做本地缓存。随着你监听的资源越来越多，它所占用的内存也会随之增加，这是无法避免的。但是，如果需要添加新的资源监听或使用新的 KubeConfig 时，我们需要重启 Informer。Informer 并没有提供显示清除缓存的动作，而缓存数据主要来自于 `Store` 对象中的 `map` 数据结构。既然我们已经知道数据存在于 `Store` 中，那么可以通过显示调用 `Store`对象来清理缓存数据。
    
*   以下通过一段代码片段, 提供一个删除的思路, 每个人的写法还是有点不一样, 如果内存占用过多清除时间太久则可以通过goroutine进行异步删除;
    
    ```
    func StopAllInformer () {
        ...
        // 清除缓存
        tmpClient := \[\]\*client.InformerClient{}
        // 这是所有集群的 informer, 遍历后获取informer 的 client, 该client 经过封装(不是重点)
        for \_, informerCli := range K8sInformerMap {
                tmpClient = append(tmpClient, informerCli)
        }
    
        for \_, cli := range tmpClient {
                // 通过 NewPreferredGVR 获取到我现在监听的资源版本
                gvrList, err := NewPreferredGVR(cli.CliSet())
                if err != nil {
                        utility.Logger.Warningf(context.TODO(), "k8s.Stop get gvr list  error. err %s", err.Error())
                        continue
                }
                // 遍历 gvrList
                for \_, gvr := range gvrList {
                        // 从informer中获取 Store 对象 (重点)
                        store := cli.Informer(gvr).GetStore()
                        // 遍历对象进行删除
                        for \_, obj := range store.List() {
                                // delete 对象
                                if err := store.Delete(obj); err != nil {
                                        utility.Logger.Warningf(context.TODO(), "k8s.Stop delete object error. err %s", err.Error())
                                }
                        }
                }
        }
    
    }
    ```
    

动态监听任意资源的事件
-----------

*   在进行集群管理时，我们可能需要对集群中某些资源的变更进行感知。这种情况下，我们可以通过监听资源变更事件来实现。众所周知，Kubernetes 的每次资源变更都会创建一个新的版本号(ResourceVersion)，随着时间的推移，**ResourceVersion** 就会过期。而我们使用的 `watch` 函数实际上就是在监听资源的版本，因此在监听过程中，我们需要记录资源的最新版本，当旧版本过期时重新进行监听。
    
*   示例代码, 展示如何监听任意资源, 及如何记录新版本重新监听, 留意代码注释.
    
    ```
    package main
    
    import (
            "context"
            "encoding/json"
            corev1 "k8s.io/api/core/v1"
            "k8s.io/apimachinery/pkg/api/meta"
            metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
            "k8s.io/apimachinery/pkg/runtime/schema"
            "k8s.io/apimachinery/pkg/watch"
            "k8s.io/client-go/dynamic"
            "log"
    )
    
    type EventObject struct {
            Metadata   metav1.ObjectMeta \`json:"metadata"\`
            ApiVersion string            \`json:"apiVersion"\`
            Kind       string            \`json:"kind"\`
    }
    
    func main() {
    
            restConfig, err := config.NewRestConfig()
            if err != nil {
                panic(err)
            }
    
            // 获取动态客户端， 如果不清楚如何获取动态客户端，可查阅我上一篇文章
            dynamicClient, err := client.NewK8sClient(restConfig).DynamicClient()
            if err != nil {
                panic(err)
            }
    
            // 因为动态客户端是可以组件任意资，所以你可以随时根据集群的资源列表增加监听资源
            gvr := schema.GroupVersionResource{
                Group:    "",
                Version:  "v1",
                Resource: "namespaces",
            }
    
            // 启动前先获取一次监听资源的版本号；
            nsList, err := dynamicClient.Resource(gvr).List(context.TODO(), metav1.ListOptions{Limit: 500})
            if err != nil {
                panic(err)
            }
            nsObjList := &corev1.NamespaceList{}
            b, \_ := nsList.MarshalJSON()
            if err := json.Unmarshal(b, nsObjList); err != nil {
                panic(err)
            }
            // 获取 resource version
            rv, err := meta.NewAccessor().ResourceVersion(nsObjList)
            if err != nil {
                panic(err)
            }
    
            // 启动监听直到天荒地老
            for {
                rv = WatchEvent(dynamicClient, gvr, rv)
            }
    
    }
    
    // WatchEvent 监听函数
    func WatchEvent(dynamicClient dynamic.Interface, gvr schema.GroupVersionResource, resourceVersion string) (rv string) {
            watcher, err := dynamicClient.Resource(gvr).Watch(context.TODO(), metav1.ListOptions{
                    TimeoutSeconds: nil,
                    // 监听的资源版本
                    ResourceVersion: resourceVersion,
                    // true 使用http 的 trunk 协议
                    Watch: true,
            })
            if err != nil {
                    panic(err)
            }
            defer watcher.Stop()
            // 最近一个版本
            latestResourceVersion := resourceVersion
            ch := watcher.ResultChan()
            for {
                event, ok := <-ch
                if ok {
                        // 当遇到错误时返回之前记录的最新版本号，这样就可以避免漏掉变更事件。
                        if event.Type == watch.Error {
                                log.Printf("watch error: %+v \\n", event.Object)
                                return latestResourceVersion
                        }
                        // 记录新的版本号
                        latestResourceVersion, err = meta.NewAccessor().ResourceVersion(event.Object)
                        if err != nil {
                                panic(err)
                        }
    
                        marshal, \_ := json.Marshal(event.Object)
                        obj := &EventObject{}
                        if err := json.Unmarshal(marshal, obj); err != nil {
                                panic(err)
                        }
                        log.Println(event.Type)
                        log.Printf("kind: %+v\\n", obj.Kind)
                        log.Printf("name: %+v\\n", obj.Metadata.Name)
                        log.Printf("created: %+v\\n", obj.Metadata.CreationTimestamp.Unix())
                        log.Printf("resourceVersion: %+v\\n", obj.Metadata.ResourceVersion)
                }
            }
    }
    ```
    
*   如何测试 ? 只需要对命名空间进行增删即可.
    
    ```
    \# kubectl create ns aaa
    \# kubectl delete ns aaa
    ```
    
*   事件输出
    
    ```
    2023/04/04 18:41:05 ADDED
    2023/04/04 18:41:05 kind: Namespace
    2023/04/04 18:41:05 name: aaa
    2023/04/04 18:41:05 created: 1680604865
    2023/04/04 18:41:05 resourceVersion: 16117101
    2023/04/04 18:41:14 MODIFIED
    2023/04/04 18:41:14 kind: Namespace
    2023/04/04 18:41:14 name: aaa
    2023/04/04 18:41:14 created: 1680604865
    2023/04/04 18:41:14 resourceVersion: 16117124
    2023/04/04 18:41:19 MODIFIED
    2023/04/04 18:41:19 kind: Namespace
    2023/04/04 18:41:19 name: aaa
    2023/04/04 18:41:19 created: 1680604865
    2023/04/04 18:41:19 resourceVersion: 16117133
    2023/04/04 18:41:19 DELETED
    2023/04/04 18:41:19 kind: Namespace
    2023/04/04 18:41:19 name: aaa
    2023/04/04 18:41:19 created: 1680604865
    2023/04/04 18:41:19 resourceVersion: 16117134
    ```
    

写在最后
----

*   3月份懒了没有保持每月一发的节奏, 4月得发两篇补一补;
*   下期分享下关于 **Envoy** 的XDS;
*   **点赞👍、关注➕、留言💬. 下篇文章见**

_**本文正在参加[「金石计划」](https://juejin.cn/post/7207698564641996856/ "https://juejin.cn/post/7207698564641996856/")**_

[Kubernetes](javascript:void(0))[Go](javascript:void(0))[掘金·金石计划](javascript:void(0))