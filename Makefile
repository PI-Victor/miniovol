PLUGIN_NAME=cloudflavor/minivol
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
# this will check if we already have exported a rootfs and removes it if so.
	@sudo rm -rf plugin.spec/rootfs/*
	@sudo rm -rf plugin.spec/rootfs/.dockerenv
	@docker create --name tmp ${PLUGIN_NAME}:rootfs
	@docker export tmp | tar -x -C plugin.spec/rootfs
	@docker rm -vf tmp

# creates the plugin based on files in plugin.spec
create:
	@docker plugin rm ${PLUGIN_NAME}:${PLUGIN_TAG} || true
