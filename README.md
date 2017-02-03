Miniovol
---
A docker volume plugin for [Minio](https://minio.io/). This plugin provisions
new Minio buckets and mounts them inside Docker volumes.  
Only compatible with [Docker 1.13](https://github.com/docker/docker/releases)
and above.  
See Docker docs about the [managed plugin system](https://docs.docker.com/engine/extend/#/installing-and-using-a-plugin).  




#### Dev stuff
To create a new version of the plugin and register it with docker do:  
```
make clean build rootfs create  
```

clean :
* remove the previously compiled `miniovol` binary and packages from
the `_output` directory.  
* remove the rootfs generated files from plugin.spec.  
* remove the previously docker built image that is used for the rootfs spec.  

build:
* compiles the `miniovol` binary.  

rootfs:
* builds a docker image that we then export to use as an OCI spec image for the
plugin.  

create:
* disables and removes the previous version of the plugin.  
* creates the new plugin based on the rootfs from the plugin.spec folder.  
* enables the newly created plugin.  

devel:
* creates a docker image with everything preinstalled.  
* launches a privileged dev container that you can then attach to to debug the
latest plugin version.  
