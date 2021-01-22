-include .env
#####################################################################################
start:
	docker-compose up -d

stop:
	docker-compose stop && docker-compose rm -f

test:
	curl -X POST http://localhost:$(PORT)/convert -H 'content-type: application/json' -d @data.json

