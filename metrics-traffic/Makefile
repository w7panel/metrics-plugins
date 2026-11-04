VERSION ?= 0.1.0

.PHONY: test build-ui build-helm build-image package
test:
	mkdir -p ../.go-cache ../.go-tmp ../.go-modcache
	cd .. && GOCACHE="$(abspath ../.go-cache)" GOTMPDIR="$(abspath ../.go-tmp)" GOMODCACHE="$(abspath ../.go-modcache)" go test ./metrics-traffic/backend ./internal/identity

build-ui:
	cd frontend && npm ci --cache ../../.npm-cache && npm run build
	cd frontend/dist && zip -q -FS -r ../dist.zip .

build-helm:
	helm package charts/metrics-traffic --version $(VERSION) --app-version $(VERSION) -d .

build-image:
	docker build -f Dockerfile -t metrics-traffic:$(VERSION) ..

package: test build-ui build-helm build-image
