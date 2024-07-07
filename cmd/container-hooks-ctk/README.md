## container-hooks-ctk

### 概览

`container-hooks-ctk` 是一个命令行工具，主要用于配置各容器运行时支持`container hooks runtime`，以向符合`OCI`规范的容器中插入`OCI Hooks`。

### 用法

**`container-hooks-ctk`所有操作均需要`root`权限**。

```shell
container-hooks-ctk [global options] command [command options] [arguments...]
```

### 全局选项

- `--debug, -d`：设置日志输出级别为Debug，默认级别为Info，单次生效
- `--quiet, -q`：设置日志仅输出Error级别的日志，默认级别为Info，单次生效
- `--version, -v`：显示`container-hooks-ctk`版本信息
- `--help, -h`：显示关于如何使用`container-hooks-ctk`的简要帮助信息

### 命令

#### runtime

`runtime`命令用于配置各容器运行时支持/移除`container-hooks-runtime`。

##### 用法

```shell
 container-hooks-ctk runtime command [command options] [arguments...]
```

##### 子命令

`container-hooks-ctk` 的`runtime`命令包括如下子命令：

- `configure`：为指定的容器运行时添加`container-hooks-runtime`支持
- `unconfigure`：为指定的容器运行时移除`container-hooks-runtime`支持

##### configure

**用法**

```shell
container-hooks-ctk runtime configure [command options] [arguments...]
```

**命令选项**

- `--runtime`：需要配置的目标容器运行时，支持`docker`, `containerd`和`crio` (默认为`docker`) 
- `--config-file`：目标容器运行时的配置文件绝对路径（默认为目标容器运行时的默认配置文件路径）
- `--name`：设置`container-hooks-runtime`的别名 （默认为`container-hooks-runtime`）
- `--path`：设置`container-hooks-runtime`执行文件的绝对路径 (默认为 `/usr/bin/container-hooks-runtime`) 
- `--set-as-default, --default`：设置`container-hooks-runtime`为目标容器运行时的默认运行时（默认为`false`）

**用法示例**

```shell
# 配置docker支持container hooks runtime，并设置为默认运行时
sudo container-hooks-ctk runtime configure --runtime=docker --default
# 对容器引擎进行配置后，需要重启容器运行时
# docker
sudo systemctl restart docker
# containerd
sudo systemctl restart containerd
# crio
sudo systemctl restart crio
```

##### unconfigure

**用法**

```shell
container-hooks-ctk runtime unconfigure [command options] [arguments...]
```

**命令选项**

- `--runtime`：需要配置的目标容器运行时，支持`docker`, `containerd`和`crio` (默认为`docker`) 
- `--config-file`：目标容器运行时的配置文件绝对路径（默认为目标容器运行时的默认配置文件路径）
- `--name`：已设置的`container-hooks-runtime`的别名

**用法示例**

```shell
# 配置docker移除container hooks runtime
sudo container-hooks-ctk runtime unconfigure --runtime=docker
# 对容器引擎进行配置后，需要重启容器运行时
# docker
sudo systemctl restart docker
# containerd
sudo systemctl restart containerd
# crio
sudo systemctl restart crio
```

#### config

`config`命令用于生成`container hooks toolkit`的配置文件。

##### 用法

```shell
container-hooks-ctk config [command options] [arguments...]
```

##### 命令选项

- `--set`：使用`key=value`模式设置配置文件中的相关配置，可指定多次

##### 用法示例

```shell
sudo container-hooks-ctk config
```

生成的默认配置文件路径为`/etc/container-hooks/config.toml`，默认内容如下：

```toml
[container-hooks]
path = "/etc/container-hooks/hooks.json"

[container-hooks-ctk]
path = "/usr/bin/container-hooks-ctk"

[container-hooks-runtime]
debug = "/etc/container-hooks/container-hooks-runtime.log"
log-level = "info"
runtimes = ["runc", "docker-runc"]
```

`[container-hooks]`

- `path`：插入的`oci hooks`文件的路径

`[container-hooks-runtime]`

- `debug`：`container hooks runtime`日志文件路径
- `log-level`：日志文件记录级别
- `runtimes`：默认底层容器运行时

`[container-hooks-ctk]`

- `path`：`container-hooks-ctk`工具的路径

**设置配置值**

```shell
# 设置container-hooks-runtime日志模式为debug
sudo container-hooks-ctk config --set container-hooks-runtime.log-level=debug
```

#### install

`install`命令用于将`container hooks toolkit`复制到`/usr/bin`目录。

##### 用法

```shell
container-hooks-ctk install [command options] [arguments...]
```

##### 命令选项

- `--toolkit-root`：`container hooks toolkit`的原始路径（必须指定此参数）

##### 用法示例

```shell
# 从/etc/test目录进行安装
sudo container-hooks-ctk install --toolkit-root=/etc/test
```

