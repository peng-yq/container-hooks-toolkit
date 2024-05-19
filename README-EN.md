# container-hooks-toolkit

<img src= "https://cdn.jsdelivr.net/gh/peng-yq/Gallery/202405191729398.png">

`container-hooks-toolkit` is used to insert custom `oci hooks` into the container configuration file (`config.json`). Components include:

- `container-hooks-runtime`
- `container-hooks-ctk`
- `container-hooks`

## container lifecycle and oci hooks

[container lifecycle and oci hooks](./container-lifecycle-hooks.md)

## container-hooks-runtime

`container-hooks-runtime` is a lightweight wrapper for the `runc` installed on the host. It injects specified `oci hooks` into the container runtime specification, then calls the local `runc` on the host, passing the modified container runtime specification with hooks settings. `runc` automatically executes the injected `oci hooks` when starting the container.

Detailed introduction and usage of `container-hooks-runtime` can be found in [README of container hooks runtime](./cmd/container-hooks-runtime/README.md).

## container-hooks-ctk

`container-hooks-ctk` is a command-line tool primarily used for configuring various container runtimes to support `container hooks runtime`, in order to insert `OCI Hooks` into containers that comply with the `OCI` specification.

Detailed introduction and usage of `container-hooks-ctk` can be found in [README of container-hooks-ctk](./cmd/container-hooks-ctk/README.md).

## container-hooks 

`container-hooks` is a dummy program, used only to determine if custom `hooks` have already been added to the current container to avoid duplication.

## how to use

`container-toolkit` supports `docker`, `containerd`, and `cri-o`. A simple use case is described below.

1. Download

```shell
git clone https://github.com/peng-yq/container-hooks-toolkit.git
```

2. Compile

```shell
make all
```

> All the following operations require `root` permissions

3. Copy to `/usr/bin`

```shell
cd bin
container-hooks-ctk install --toolkit-root=$(pwd)
```

4. Generate configuration file

```shell
container-hooks-ctk config
```

5. Configure using `docker` as an example

```shell
container-hooks-ctk runtime configure --runtime=docker --default
systemctl restart docker
```

6. Write custom `oci hooks`, format as follows. The first `prestart hook` must include `container-hooks` to prevent re-adding defined `hooks`, to be written into the `/etc/container-hooks/hooks.json` file (this path can be modified in the configuration file):

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

7. Run the container, execute the custom `hook`

```shell
docker run image:tag
```

## customized usage

Some more customized ideas:

Case 1: Automatically perform signature verification and integrity check on the container before it starts

In this case, simply writing `hooks` to `/etc/container-hooks/hooks.json` will not work, as we cannot know the startup image information of each container in advance. You can modify the project for secondary development, not using the method of reading hooks from files, but directly inserting them in the code and adjusting parameters according to the container's configuration.

Code parts that need modification:

1. `/internal/runtime`
2. `/internal/modifier`

[Reference material](https://peng-yq.github.io/2023/09/07/runc/)

Case 2: Output the `bundle` path of each container

...

## References/Acknowledgements

1. [nvidia-container-toolkit](https://github.com/NVIDIA/nvidia-container-toolkit)
2. [oci-add-hooks](https://github.com/awslabs/oci-add-hooks)