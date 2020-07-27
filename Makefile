PLUGIN_ID=srcf.redact
PLUGIN_VERSION=1.1.1

CWD := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

export GO111MODULE=on

.PHONY: build
build:
	rm -rf dist/

	for os in linux darwin windows; do \
		export TARGET=$(CWD)/dist/$$os/$(PLUGIN_ID); \
		mkdir -p $$TARGET; \
		cd $(CWD)/src && env GOOS=$$os GOARCH=amd64 go build -o $$TARGET/plugin-$$os-amd64; cd ..; \
		cp plugin.json $$TARGET; \
		tar -C $$TARGET/.. -cvzf $(CWD)/dist/$(PLUGIN_ID)-$(PLUGIN_VERSION)-$$os-amd64.tar.gz $(PLUGIN_ID); \
	done

.PHONY: clean
clean:
	rm -rf dist/
