COMPONENTS = server

IMAGE_BASE=docker.io/mfenwick100/upanddowntheriver/server:master

CURRENT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
OUTDIR=_output

.PHONY: test ${OUTDIR} ${COMPONENTS}

all: compile

compile: ${OUTDIR} ${COMPONENTS}

${COMPONENTS}:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/$@/$@ ./cmd/$@
	docker build -t $(IMAGE_BASE)$@:$(IMAGE_TAG) ./cmd/$@
	mv cmd/$@/$@ $(OUTDIR)
	# gcloud docker -- push $(IMAGE_BASE)$@:$(IMAGE_TAG)

docker-image: $(COMPONENTS)
	$(foreach p,${COMPONENTS},cd ${CURRENT_DIR}/cmd/$p; docker build -t $(IMAGE_BASE)${p}:$(IMAGE_TAG) .;)

test:
	go test ./pkg/...

clean:
	rm -rf ${OUTDIR}
	$(foreach p,${COMPONENTS},rm -f cmd/$p/$p;)

${OUTDIR}:
	mkdir -p ${OUTDIR}

fmt:
	go fmt ./cmd/... ./pkg/...

vet:
	go vet ./cmd/... ./pkg/...
