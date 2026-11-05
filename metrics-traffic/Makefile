VERSION ?= 0.1.0

.PHONY: test build-ui build-helm build-image package
test:
	go test ./...

build-ui:
	cd frontend && npm ci && npm run build
	cd frontend/dist && zip -q -FS -r ../dist.zip .

build-helm:
	helm package charts/metrics-traffic --version $(VERSION) --app-version $(VERSION) -d .

build-image:
	docker build -t metrics-traffic:$(VERSION) .

package: test build-ui build-helm build-image
