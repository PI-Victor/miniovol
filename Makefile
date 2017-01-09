PLUGIN_NAME=cloudflavor/miniovol
PLUGIN_TAG=latest

.PHONY: clean
build:
	@echo 'Compiling miniovol'
	@cd cmd/miniovol && go build -o ../../_out/bin/miniovol -v .

clean:
	@echo 'Removing old _out dir'
	@rm -rf _out/
	@echo 'Cleaning out generated rootfs if it exists...'
	@sudo rm -rf plugin.spec/rootfs/*
	@sudo rm -rf plugin.spec/rootfs/.dockerenv

docker:
	@echo 'Building plugin Docker image...'
	@docker build -q -t ${PLUGIN_NAME}:rootfs .

# creates the rootfs needed to distribute the plugin.
rootfs:
	@echo 'Creating new rootfs for plugin...'
	@docker create --name tmp ${PLUGIN_NAME}:rootfs
	@docker export tmp | tar -x -C plugin.spec/rootfs
	@docker rm -vf tmp
	@docker rmi ${PLUGIN_NAME}:rootfs

# creates the plugin based on files in plugin.spec
create:
	@echo 'Removing plugin if it exists...'
	@docker plugin disable ${PLUGIN_NAME}:${PLUGIN_TAG} || true 
	@docker plugin rm -f ${PLUGIN_NAME}:${PLUGIN_TAG} || true
	@echo 'Create new plugin...'
	@sudo docker plugin create ${PLUGIN_NAME}:${PLUGIN_TAG} plugin.spec
	@docker plugin enable ${PLUGIN_NAME}:${PLUGIN_TAG}
