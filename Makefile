.PHONY: clean
build:
	@echo 'Compiling miniovol'
	@cd cmd/miniovol && go build -o ../../_out/bin/miniovol -v .

clean:
	@echo 'Removing old _out dir'
	@rm -rf _out/
