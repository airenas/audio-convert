test: 
	go test ./...

build:
	cd cmd/audio-convert/ && go build .

run:
	cd cmd/audio-convert/ && go run . -c config.yml	

docker-build:
	cd deploy && $(MAKE) clean dbuild	

docker-push:
	cd deploy && $(MAKE) clean dpush

clean:
	rm -f cmd/audio-convert/audio-convert
	cd deploy && $(MAKE) clean

