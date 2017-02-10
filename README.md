Miniovol
---
**The original Docker volume plugin, for minio, can be found in the [minfs repo](https://github.com/minio/minfs#docker-plugin)  
This repo can serve as an extensive example for how to write docker volume plugins for Docker 1.13.    
If you're looking for running things in production, then always choose the official supported software rather than 3rd party.**  

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
