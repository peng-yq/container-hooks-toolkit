# container-hooks-toolkit

[English](./README-EN.md) | 简体中文

<img src= "https://cdn.jsdelivr.net/gh/peng-yq/Gallery/202405191729398.png">

`container-hooks-toolkit`用于向容器的配置文件(`config.json`)中插入自定义的`oci hooks`，组件包括：

- `container-hooks-runtime`
- `container-hooks-ctk`
- `container-hooks`

## container lifecycle and oci hooks

[container lifecycle and oci hooks](./container-lifecycle-hooks.md)

## container-hooks-runtime

`container-hooks-runtime`是对主机上安装的`runc`的轻量级包装器，通过将指定的`oci hooks`注入容器的运行时规范，然后调用主机本地的`runc`，并传递修改后的带有钩子设置的容器运行时规范。`runc`在启动容器时，会自动运行注入的`oci hooks`。

`container-hooks-runtime`的详细介绍和用法见[README of container hooks runtime](./cmd/container-hooks-runtime/README.md)。

## container-hooks-ctk

`container-hooks-ctk` 是一个命令行工具，主要用于配置各容器运行时支持`container hooks runtime`，以向符合`OCI`规范的容器中插入`OCI Hooks`。

`container-hooks-ctk`的详细介绍和用法见[README of container-hooks-ctk](./cmd/container-hooks-ctk/README.md)。

## container-hooks 

`container-hooks`是一个空程序，仅用于判断当前容器是否已经添加自定义`hooks`，避免重复添加。

## how to use

`container-toolkit`支持`docker`，`containerd`和`cri-o`，简单用例介绍如下。

1. 下载

```shell
git clone https://github.com/peng-yq/container-hooks-toolkit.git
```

2. 编译

```shell
make all
```

> 下面的所有操作均需要`root`权限

1. 复制到`/usr/bin`

```shell
cd bin
container-hooks-ctk install --toolkit-root=$(pwd)
```

4. 生成配置文件

```shell
container-hooks-ctk config
```

5. 以`docker`为例，进行配置

```shell
container-hooks-ctk runtime configure --runtime=docker --default
systemctl restart docker
```

6. 编写自定义`oci hooks`，格式如下，必须添加第一个`prestart hook`中的`container-hooks`用于避免重复添加定义`hooks`，需要写入至`/etc/container-hooks/hooks.json`文件中（此路径可在配置文件中修改）

```json
{
  "hooks": {
    "prestart": [
        {
            "path": "/usr/bin/container-hooks",
        }
    ],
    "createRuntime": [
        {
            "path": "/usr/bin/fix-mounts",
            "args": ["fix-mounts", "arg1", "arg2"],
            "env":  [ "key1=value1"]
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
}
```

7. 运行容器，执行自定义`hook`

```shell
docker run image:tag
```

## customized usage

提供一些更加定制化的思路：

案例1：在容器启动前自动对容器进行签名验证和完整性校验

此时直接编写`hooks`至`/etc/container-hooks/hooks.json`就行不通了，因为我们无法提前预知每个容器的启动镜像信息。可以对项目进行二次开发，不采用读取文件中的钩子的形式，而是直接在代码中进行插入并根据容器的配置进行调整参数。

需要修改的代码部分：

1. `/internel/runtime`
2. `/internel/modifier`

[可参考的资料](https://peng-yq.github.io/2023/09/07/runc/)

案例2：输出每个容器的`bundle`路径

...

## 参考/致谢

1. [nvidia-container-toolkit](https://github.com/NVIDIA/nvidia-container-toolkit)
2. [oci-add-hooks](https://github.com/awslabs/oci-add-hooks)