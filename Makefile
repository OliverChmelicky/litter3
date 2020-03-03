all:
	echo "Hello, use one of the options dev"

dev:
	docker-compose -f docker-compose-dev.yml up