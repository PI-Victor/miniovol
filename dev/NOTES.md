{"version":"1","accessKey":"OKWUYY6OHNKMJB67WPU0","secretKey":"MWShkNCAlhmeDXdL6M3qb8N8lhAjM2ih83+pgjzD"}
mount -t minfs http://192.168.0.115:9000/docker-volumes /mnt


see how to mount /dev/fuse in a privileged container.
https://github.com/docker/docker/issues/9448
```
docker run -ti --cap-add SYS_ADMIN --device /dev/fuse
```

how to create a new volume with this driver.
docker volume create -d cloudflavor/miniovol:latest -o server="192.168.0.115:9000" -o accessKey=OKWUYY6OHNKMJB67WPU0 -o secretKey=MWShkNCAlhmeDXdL6M3qb8N8lhAjM2ih83+pgjzD miniovol

{"version":"1","accessKey":"26W01YWCHGG3U5XLL0YB","secretKey":"DmvQ8AvidnrkkJFRUVjEXKALlkx7G12XO6sdCrJe"}
