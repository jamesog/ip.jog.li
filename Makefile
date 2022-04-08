IMAGE := ip.jog.li
VERSION := $(shell git describe --tags --always)

all: build

build:
	docker build -t $(IMAGE):$(VERSION) .

run:
	docker run -p 8000:8000 $(IMAGE):$(VERSION)

deploy:
	kubectl apply -f deployment.yaml
