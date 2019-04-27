IMAGE := ip.jog.li
VERSION := 0.3.1

GCR_REGION  := eu
GCR_HOST    := $(GCR_REGION).gcr.io
GCP_PROJECT := ipjogli

all: build

build:
	docker build -t $(IMAGE):$(VERSION) .

run:
	docker run -p 8000:8000 $(IMAGE):$(VERSION)

push:
	docker tag $(IMAGE):$(VERSION) $(GCR_HOST)/$(GCP_PROJECT)/$(IMAGE):$(VERSION)
	docker push $(GCR_HOST)/$(GCP_PROJECT)/$(IMAGE):$(VERSION)
	gsutil acl ch -r -u AllUsers:READ gs://$(GCR_REGION).artifacts.$(GCP_PROJECT).appspot.com/

deploy:
	kubectl apply -f deployment.yaml
