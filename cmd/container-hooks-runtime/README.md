## container-hooks-runtime

### 概览

`container-hooks-runtime`是对主机上安装的`runc`的轻量级包装器，通过将指定的`oci hooks`注入容器的运行时规范，然后调用主机本地的`runc`，并传递修改后的带有钩子设置的容器运行时规范。`runc`在启动容器时，会自动运行注入的`oci hooks`。

<img src= "https://cdn.jsdelivr.net/gh/peng-yq/Gallery/202405191729398.png">

### 用法

**`container-hooks-runtime`需要通过`container-hooks-ctk`配置各容器运行时，配置完毕后容器运行时在执行容器命令时会自动调用`container-hooks-runtime`**。

`container-hooks-runtime`的配置文件包含了`container-hooks-runtime`的配置选项，路径为`/etc/container-hooks/config.toml`，支持对其修改从而定义可信容器运行时的日志文件路径、日志级别以及底层运行时：

```toml
[container-hooks-runtime]
debug = "/etc/container-hooks/container-hooks-runtime.log"
log-level = "info"
runtimes = ["runc", "docker-runc"]
```

`container-hooks-runtime`日志记录了容器生命周期的相关记录，默认路径为`/etc/container-hooks/container-hooks-runtime.log`。

