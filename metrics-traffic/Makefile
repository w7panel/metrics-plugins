VERSION ?= 0.1.0
IMAGE_REPO ?= ccr.ccs.tencentyun.com/afan-public/metrics-traffic
CHART_PACKAGE = metrics-traffic-$(VERSION).tgz

.PHONY: test build-ui build-helm build-image docker-build push-zpk package full check-release

check-release:
	@for tool in npm zip helm docker zpk; do command -v "$$tool" >/dev/null || { echo "Missing $$tool" >&2; exit 2; }; done
	@test -n "$${ZPK_IDC_PASSWORD:-$${PASSWORD:-}}" || { echo 'Set ZPK_IDC_PASSWORD or PASSWORD' >&2; exit 2; }

.NOTPARALLEL: full
full: check-release build-ui build-helm build-image push-zpk

test:
	go test ./...

build-ui:
	cd frontend && npm ci && npm run build
	cd frontend/dist && zip -q -FS -r ../dist.zip .

build-helm:
	helm package charts/metrics-traffic --version $(VERSION) --app-version $(VERSION) -d .

docker-build:
	docker build -t $(IMAGE_REPO):$(VERSION) .

build-image: docker-build
	docker push $(IMAGE_REPO):$(VERSION)

push-zpk:
	@test -f "$(CHART_PACKAGE)" && test -f frontend/dist.zip || { echo 'Build Helm and UI first' >&2; exit 1; }
	@test -n "$${ZPK_IDC_PASSWORD:-$${PASSWORD:-}}" || { echo 'Set ZPK_IDC_PASSWORD or PASSWORD' >&2; exit 2; }
	@zpk login --username=w7admin --password="$${ZPK_IDC_PASSWORD:-$${PASSWORD}}" --host=zpk.w7.cc
	zpk use --name=metrics-traffic
	zpk attach add --path="$(CHART_PACKAGE)" --type=helm
	zpk attach add --path=frontend/dist.zip --type=frontend
	zpk push

package: test build-ui build-helm docker-build
