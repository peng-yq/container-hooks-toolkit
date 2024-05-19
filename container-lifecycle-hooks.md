## lifecycle

容器中的钩子和容器的生命周期息息相关，钩子能使容器感知其生命周期内的事件，并且当相应的生命周期钩子被调用时运行指定的代码。[oci runtime-spec](https://github.com/opencontainers/runtime-spec/blob/main/runtime.md#lifecycle)对容器生命周期的描述如下（单纯从低级运行时创建容器开始，不包括镜像）。

> The lifecycle describes the timeline of events that happen from when a container is created to when it ceases to exist.
>
> - OCI compliant runtime's create command is invoked with a reference to the location of the bundle and a unique identifier.
> - **The container's runtime environment MUST be created according to the configuration in config.json**. If the runtime is unable to create the environment specified in the config.json, it MUST generate an error. While the resources requested in the config.json MUST be created, the user-specified program (from process) MUST NOT be run at this time. Any updates to config.json after this step MUST NOT affect the container.
> - The **prestart hooks** MUST be invoked by the runtime. If any prestart hook fails, the runtime MUST generate an error, stop the container, and continue the lifecycle at step 12.
> - The **createRuntime hooks** MUST be invoked by the runtime. If any createRuntime hook fails, the runtime MUST generate an error, stop the container, and continue the lifecycle at step 12.
> - The **createContainer hooks** MUST be invoked by the runtime. If any createContainer hook fails, the runtime MUST generate an error, stop the container, and continue the lifecycle at step 12.
> - Runtime's start command is invoked with the unique identifier of the container.
> - The **startContainer hooks** MUST be invoked by the runtime. If any startContainer hook fails, the runtime MUST generate an error, stop the container, and continue the lifecycle at step 12.
> - The runtime MUST run the user-specified program, as specified by process.
> - The **poststart hooks** MUST be invoked by the runtime. If any poststart hook fails, the runtime MUST log a warning, but the remaining hooks and lifecycle continue as if the hook had succeeded.
> - The container process exits. This MAY happen due to erroring out, exiting, crashing or the runtime's kill operation being invoked.
> - Runtime's delete command is invoked with the unique identifier of the container.
> - The container MUST be destroyed by undoing the steps performed during create phase (step 2).
> - The poststop hooks MUST be invoked by the runtime. If any poststop hook fails, the runtime MUST log a warning, but the remaining hooks and lifecycle continue as if the hook had succeeded.

Lifecycle定义了容器从创建到退出之间的时间轴：

1. 容器开始创建：**通常为OCI规范运行时（runc）调用create命令+bundle+container id**
2. 容器运行时环境创建中： 根据容器的config.json中的配置进行创建，此时用户指定程序还未运行，这一步后所有对config.json的更改均不会影响容器
3. prestart hooks
4. createRuntime hooks
5. createContainer hooks
6. 容器启动：**通常为OCI规范运行时（runc）调用start命令+container id**
7. startContainer hooks
8. 容器执行用户指定程序
9. poststart hooks：任何poststart钩子执行失败只会log a warning，不影响其他生命周期（操作继续执行）就好像钩子成功执行一样
10. 容器进程退出：error、正常退出和运行时调用kill命令均会导致
11. 容器删除：**通常为OCI规范运行时（runc）调用delete命令+container id**
12. 容器摧毁：**区别于容器删除，3、4、5、7的钩子执行失败除了生成一个error外，会直接跳到这一步**。撤销第二步创建阶段执行的操作。
13. poststop hooks：执行失败后的操作和poststart一致

可以看到oci定义的容器生命周期中，如果在容器的config.json中定义了钩子，runc必须执行钩子，并且时间节点在前的钩子执行成功后才能执行下一个钩子；若有一个钩子执行失败，则会报错并摧毁容器（在容器创建后执行的钩子失败，并不会删除容器，而是启动失败）。

除了上述oci通过runc创建并启动容器的流程来对容器生命周期的描述外，docker和k8s也有各自对容器生命周期的描述（出于容器的不同状态），均符合oci规范。

Docker

<img src="https://cdn.jsdelivr.net/gh/peng-yq/Gallery/img/202309081015834.png">

[Pod 的生命周期-K8s](https://kubernetes.io/zh-cn/docs/concepts/workloads/pods/pod-lifecycle/#pod-lifetime)

## oci hooks

结合生命周期来看。

hooks (object, OPTIONAL) ：配置与容器生命周期相关的特定操作，按顺序进行调用：

- prestart (array of objects, OPTIONAL, DEPRECATED) :所有类型的钩子均有相同的键
  - path (string, REQUIRED) ：绝对路径
  - args (array of strings, OPTIONAL)
  - env (array of strings, OPTIONAL) 
  - timeout (int, OPTIONAL) ：终止钩子的秒数
- createRuntime (array of objects, OPTIONAL)
- createContainer (array of objects, OPTIONAL)
- startContainer (array of objects, OPTIONAL)
- poststart (array of objects, OPTIONAL)
- poststop (array of objects, OPTIONAL) 

容器的状态必须通过 stdin 传递给钩子，以便它们可以根据容器的当前状态执行相应的工作。

### Prestart

Prestart钩子必须作为创建操作的一部分，在运行时环境创建完成后（根据 config.json 中的配置），但在执行 pivot_root 或任何同等操作之前调用。

**废弃，被后面三个钩子所取代**

### CreateRuntime Hooks

createRuntime钩子必须作为创建操作的一部分，在运行时环境创建完成后（根据 config.json 中的配置），但在执行 pivot_root 或任何同等操作之前调用。

**在容器命名空间被创建后调用**。

> createRuntime 钩子的定义目前未作明确规定，钩子作者只能期望运行时创建挂载命名空间并执行挂载操作。运行时可能尚未执行其他操作，如 cgroups 和 SELinux/AppArmor 标签

### CreateContainer Hooks

createContainer钩子必须作为创建操作的一部分，在运行时环境创建完成后（根据 config.json 中的配置），但在执行 pivot_root 或任何同等操作之前调用。

**在执行 pivot_root 操作之前，但在创建和设置挂载命名空间之后调用**。

### StartContainer Hooks

StartContainer钩子作为启动操作的一部分，**必须在执行用户指定的进程之前调用startContainer挂钩**。此钩子可用于在容器中执行某些操作，例如在容器进程生成之前在linux上运行ldconfig二进制文件。

### Poststart

Poststart钩子**必须在用户指定的进程执行后、启动操作返回前调用**。例如，此钩子可以通知用户容器进程已生成。

### Poststop

Poststart钩子**必须在容器删除后、删除操作返回前调用**。清理或调试函数就是此类钩子的例子。

### summary

**namespace是指path以及钩子必须在指定的namespace中解析或调用**。

| Name                    | Namespace | When                                                         |
| ----------------------- | --------- | ------------------------------------------------------------ |
| `prestart` (Deprecated) | runtime   | After the start operation is called but before the user-specified program command is executed. |
| `createRuntime`         | runtime   | During the create operation, after the runtime environment has been created and before the pivot root or any equivalent operation. |
| `createContainer`       | container | During the create operation, after the runtime environment has been created and before the pivot root or any equivalent operation. |
| `startContainer`        | container | After the start operation is called but before the user-specified program command is executed. |
| `poststart`             | runtime   | After the user-specified process is executed but before the start operation returns. |
| `poststop`              | runtime   | After the container is deleted but before the delete operation returns. |

```json
"hooks": {
    "prestart": [
        {
            "path": "/usr/bin/fix-mounts",
            "args": ["fix-mounts", "arg1", "arg2"],
            "env":  [ "key1=value1"]
        },
        {
            "path": "/usr/bin/setup-network"
        }
    ],
    "createRuntime": [
        {
            "path": "/usr/bin/fix-mounts",
            "args": ["fix-mounts", "arg1", "arg2"],
            "env":  [ "key1=value1"]
        },
        {
            "path": "/usr/bin/setup-network"
        }
    ],
    "createContainer": [
        {
            "path": "/usr/bin/mount-hook",
            "args": ["-mount", "arg1", "arg2"],
            "env":  [ "key1=value1"]
        }
    ],
    "startContainer": [
        {
            "path": "/usr/bin/refresh-ldcache"
        }
    ],
    "poststart": [
        {
            "path": "/usr/bin/notify-start",
            "timeout": 5
        }
    ],
    "poststop": [
        {
            "path": "/usr/sbin/cleanup.sh",
            "args": ["cleanup.sh", "-f"]
        }
    ]
}
```

## hooks supported in docker and k8s

docker并不支持通过参数直接添加使用hooks，或者说现在不支持（调研中发现有一两篇博客写到可以通过docker run --hooks-path添加钩子，但通过`docker --help | grep hook`命令找不到任何和hook相关的参数）。

k8s只支持[postStart和preStop hooks](https://kubernetes.io/zh-cn/docs/tasks/configure-pod-container/attach-handler-lifecycle-event/)，通过在pod的yaml文件中定义即可。但实际上我认为prestart（或者说其细分后的三个hooks）更加实用一些，比如在创建或启动容器前进行一些网络检查和配置。有意思的是在[k8s的issue中有一个关于“PreStart lifecycle hook Required”的讨论](https://github.com/kubernetes/kubernetes/issues/96560)，讨论提到了修改Dockerfile并在entrypoint.sh中执行，以及使用[k8s提供的init容器](https://kubernetes.io/zh-cn/docs/concepts/workloads/pods/init-containers/)对prestart的代替。

补充：[PreStart and PostStop event hooks #140](https://github.com/kubernetes/kubernetes/issues/140)中有老哥回答k8s不会再支持prestart hook了，取而代之的是init容器。