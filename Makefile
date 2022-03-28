#####################################################################################
## print usage information
help:
	@echo 'Usage:'
	@cat ${MAKEFILE_LIST} | grep -e "^## " -A 1 | grep -v '\-\-' | sed 's/^##//' | cut -f1 -d":" | \
		awk '{info=$$0; getline; print "  " $$0 ": " info;}' | column -t -s ':' | sort 
.PHONY: help
#####################################################################################
## call units tests
test/unit: 
	go test -v -race -count 1 ./...
.PHONY: test/unit
#####################################################################################
## code vet and lint
test/lint: 
	go vet ./...
	go install golang.org/x/lint/golint@latest
	golint -set_exit_status ./...
.PHONY: test/lint
#####################################################################################
run:
	cd cmd/audio-convert/ && go run . -c config.yml	
#####################################################################################
docker/build:
	cd build && $(MAKE) clean dbuild	
#####################################################################################
docker/push:
	cd build && $(MAKE) clean dpush
#####################################################################################
## clean temporary docker build artifacts
clean:
	go mod tidy
	go clean
	rm -f cmd/audio-convert/audio-convert
	cd deploy && $(MAKE) clean

