PLUGIN_NAME=cloudflavor/miniovol
PLUGIN_TAG=latest

.PHONY: clean
build:
	@echo 'Compiling miniovol...'
	@cd cmd/miniovol && go build -o ../../_out/bin/miniovol -v .

# remove miniovol bin, rootfs artifacts and base docker rootfs image.
clean:
	@echo 'Removing old _out dir...'
	@rm -rf _out/
	@echo 'Cleaning out generated rootfs artifacts if they exist...'
	@sudo rm -rf plugin.spec/rootfs/*
	@sudo rm -rf plugin.spec/rootfs/.dockerenv
	@echo 'Removing ${PLUGIN_NAME}:rootfs image if it exists'
	@docker rmi -f ${PLUGIN_NAME}:rootfs || true
	@docker rmi -f ${PLUGIN_NAME}:dev || true
	@docker rm -vf tmp || true

# creates the rootfs needed to distribute the plugin.
rootfs:
	@echo 'Building plugin Docker image...'
	@docker build -t ${PLUGIN_NAME}:rootfs .
	@echo 'Creating new rootfs for plugin...'
	@docker create --name tmp ${PLUGIN_NAME}:rootfs .
	@docker export tmp | tar -x -C plugin.spec/rootfs
	@docker rm -vf tmp
	@docker rmi ${PLUGIN_NAME}:rootfs

# creates the plugin based on files in plugin.spec.
create:
	@echo 'Removing any previous versions of the plugin...'
	@docker plugin disable ${PLUGIN_NAME}:${PLUGIN_TAG} || true
	@docker plugin rm -f ${PLUGIN_NAME}:${PLUGIN_TAG} || true
	@echo 'Creating and enabling new docker plugin...'
	@sudo docker plugin create ${PLUGIN_NAME}:${PLUGIN_TAG} plugin.spec
	@docker plugin enable ${PLUGIN_NAME}:${PLUGIN_TAG}

# create a dev ready docker image to test out manually plugin functionality.
devel: build
	@echo 'Removing previously built dev docker image...'
	@docker rmi -f ${PLUGIN_NAME}:dev || true
	@echo 'Building dev docker image...'
	@docker build -t ${PLUGIN_NAME}:dev -f dev/Dockerfile.dev .
