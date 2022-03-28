-include ../version 
#####################################################################################
service=airenas/audio-convert-service
########### DOCKER ##################################################################
tag=$(service):$(version)
dbuild:
	cd ../ && docker build -t $(tag) --build-arg BUILD_VERSION=$(version) -f build/Dockerfile .

dpush: dbuild
	docker push $(tag)

dscan: dbuild
	docker scan --accept-license $(tag)	
#####################################################################################
.PHONY: dbuild dpush
