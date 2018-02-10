IMAGE := ip.jog.li
VERSION := 0.3.0

GCR_HOST := eu.gcr.io
GCP_PROJECT := ipjogli

all: build

build:
	docker build -t $(IMAGE):$(VERSION) .

run:
	docker run -p 8000:8000 $(IMAGE):$(VERSION)

push:
	docker tag $(IMAGE):$(VERSION) $(GCR_HOST)/$(GCP_PROJECT)/$(IMAGE):$(VERSION)
	gcloud docker -- push $(GCR_HOST)/$(GCP_PROJECT)/$(IMAGE):$(VERSION)
