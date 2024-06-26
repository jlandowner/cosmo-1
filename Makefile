
# Image URL to use all building/pushing image targets
VERSION ?=
PRERELEASE ?= false
QUICK_BUILD ?= no

MANAGER_VERSION   ?= $(VERSION)
DASHBOARD_VERSION ?= $(VERSION)
COSMOCTL_VERSION  ?= $(VERSION)

CHART_VERSION   ?= $(VERSION)

IMG_MANAGER ?= cosmo-controller-manager:$(MANAGER_VERSION)
IMG_DASHBOARD ?= cosmo-dashboard:$(DASHBOARD_VERSION)
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:generateEmbeddedObjectMeta=true"

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

##---------------------------------------------------------------------
##@ General
##---------------------------------------------------------------------

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: all
all: manager cosmoctl dashboard

##---------------------------------------------------------------------
##@ Development
##---------------------------------------------------------------------

.PHONY: manifests
manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
ifeq ($(QUICK_BUILD),no)
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./api/..." output:crd:artifacts:config=config/crd/bases
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./internal/webhooks" output:crd:artifacts:config=config/crd/bases
endif

.PHONY: generate
generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
ifeq ($(QUICK_BUILD),no)
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./api/..."
endif

.PHONY: proto-generate 
proto-generate:  ## Generate code protocol buffer api.
	make -C proto/ all

.PHONY: fmt
fmt: go ## Run go fmt against code.
ifeq ($(QUICK_BUILD),no)
	$(GO) fmt ./...
endif

.PHONY: vet
vet: go embed-traefik-plugins ## Run go vet against code.
ifeq ($(QUICK_BUILD),no)
	$(GO) vet ./...
endif

##---------------------------------------------------------------------
##@ Test
##---------------------------------------------------------------------
TEST_FILES ?= ./... ./traefik-plugins/src/github.com/cosmo-workspace/cosmoauth/
COVER_PROFILE ?= cover.out
#TEST_OPTS ?= --ginkgo.focus 'Dashboard server \[User\]' -ginkgo.v -ginkgo.progress -test.v > test.out 2>&1

.PHONY: clear-snapshots
clear-snapshots: ## Clear snapshots
	-find . -type f | grep __snapshots__ | grep -v "/web/" | grep -v "/charts/" | xargs rm -f

.PHONY: ingressroute.yaml
ingressroute.yaml: helm config/crd/traefik/traefik.io_ingressroutes.yaml
config/crd/traefik/traefik.io_ingressroutes.yaml:
	mkdir -p config/crd/traefik
	$(HELM) dependency update ./charts/cosmo
	tar -xvf ./charts/cosmo/charts/traefik-*.tgz -O traefik/crds/traefik.io_ingressroutes.yaml > config/crd/traefik/traefik.io_ingressroutes.yaml

.PHONY: go-test.env
go-test.env:
	@echo KUBEBUILDER_ASSETS=$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path) > ./.vscode/go-test.env
	@echo PATH=$(PATH) >> ./.vscode/go-test.env

.PHONY: test
test: manifests generate fmt vet envtest go-test.env go-test ## Run tests.

.PHONY: go-test
go-test: go ingressroute.yaml
ifeq ($(QUICK_BUILD),no)
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" \
	$(GO) test $(TEST_FILES) -coverpkg="./..." -coverprofile $(COVER_PROFILE) $(TEST_OPTS)
endif

.PHONY: test-all-k8s-versions
test-all-k8s-versions: go manifests generate fmt vet envtest ## Run tests on targeting k8s versions.
ifeq ($(QUICK_BUILD),no)
	-KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use 1.30.x -p path)" $(GO) test ./... -coverprofile $(COVER_PROFILE)
	-KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use 1.29.x -p path)" $(GO) test ./... -coverprofile $(COVER_PROFILE)
	-KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use 1.28.x -p path)" $(GO) test ./... -coverprofile $(COVER_PROFILE)
endif

.PHONY: clear-snapshots-ui
clear-snapshots-ui: ## Clear snapshots ui
	-find ./web -type f | grep __snapshots__ | xargs rm -f

.PHONY: ui-test
ui-test: ## Run UI tests.
	cd web/dashboard-ui && yarn install && yarn test  --coverage --run

##---------------------------------------------------------------------
##@ Build
##---------------------------------------------------------------------
.PHONY: manager
manager: go generate fmt vet ## Build manager binary.
	CGO_ENABLED=0 $(GO) build -o bin/manager ./cmd/controller-manager/main.go

.PHONY: cosmoctl
cosmoctl: go generate fmt vet ## Build cosmoctl binary.
	CGO_ENABLED=0 $(GO) build -o bin/cosmoctl ./cmd/cosmoctl/main.go

.PHONY: dashboard
dashboard: go generate fmt vet ## Build dashboard binary.
	CGO_ENABLED=0 $(GO) build -o bin/dashboard ./cmd/dashboard/main.go

.PHONY: update-version
update-version: kustomize ## Update version in version.go.
ifndef VERSION
	@echo "Usage: make update-version VERSION=v9.9.9"
	@exit 9
else
ifeq ($(shell expr $(VERSION) : '^v[0-9]\+\.[0-9]\+\.[0-9]\+'),0)
	@echo "Usage: make update-version VERSION=v9.9.9"
	@exit 9
endif
endif
	sed -i.bk -e "s/v[0-9]\+.[0-9]\+.[0-9]\+.* cosmo-workspace/${MANAGER_VERSION} cosmo-workspace/" ./cmd/controller-manager/main.go
	sed -i.bk -e "s/v[0-9]\+.[0-9]\+.[0-9]\+.* cosmo-workspace/${DASHBOARD_VERSION} cosmo-workspace/" ./internal/dashboard/root.go
	sed -i.bk -e "s/v[0-9]\+.[0-9]\+.[0-9]\+.* cosmo-workspace/${COSMOCTL_VERSION} cosmo-workspace/" ./internal/cmd/version/version.go
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG_MANAGER}
	cd config/dashboard && $(KUSTOMIZE) edit set image dashboard=${IMG_DASHBOARD}
	sed -i.bk \
		-e "s/version: [0-9]\+.[0-9]\+.[0-9]\+.*/version: ${CHART_VERSION:v%=%}/" \
		-e "s/appVersion: v[0-9]\+.[0-9]\+.[0-9]\+.*/appVersion: ${VERSION}/" \
		-e 's;artifacthub.io/prerelease: "\(true\|false\)";artifacthub.io/prerelease: "$(PRERELEASE)";' \
		charts/cosmo/Chart.yaml

##---------------------------------------------------------------------
##@ Run
##---------------------------------------------------------------------

LOG_LEVEL ?= 3

COOKIE_DOMAIN ?= 
COOKIE_SESSION_NAME ?= test-cosmo-auth
COOKIE_HASHkEY  ?= 12345678901234567890123456789012
COOKIE_BLOCKKEY ?= ----+----1----+----2----+----3--


# Run against the configured Kubernetes cluster in ~/.kube/config
.PHONY: run-dashboard
run-dashboard: go generate fmt vet manifests ## Run dashboard against the configured Kubernetes cluster in ~/.kube/config.
ifndef COOKIE_DOMAIN
	@echo "Usage: make run-dashboard COOKIE_DOMAIN=xxxx.xxx"
	@exit 9
endif
	$(GO) run ./cmd/dashboard/main.go \
		--zap-log-level $(LOG_LEVEL) \
		--zap-time-encoding=iso8601 \
		--cookie-session-name=$(COOKIE_SESSION_NAME) \
		--cookie-domain=$(COOKIE_DOMAIN) \
		--cookie-hashkey=$(COOKIE_HASHkEY) \
		--cookie-blockkey=$(COOKIE_BLOCKKEY) \
		--insecure

.PHONY: run-dashboard-ui
run-dashboard-ui: ## Run dashboard-ui.
	cd web/dashboard-ui && yarn install && yarn start

.PHONY: run
run: go generate fmt vet manifests ## Run controller-manager against the configured Kubernetes cluster in ~/.kube/config.
	$(GO) run ./cmd/controller-manager/main.go \
		--zap-log-level $(LOG_LEVEL) \
		--zap-time-encoding=iso8601 \
		--metrics-bind-address :8085 \
		--cert-dir .

embed-traefik-plugins: cmd/controller-manager/traefik-plugins.tar.gz
cmd/controller-manager/traefik-plugins.tar.gz:
	make -C traefik-plugins/src/github.com/cosmo-workspace/cosmoauth vendor
	cd traefik-plugins && tar zcvf traefik-plugins.tar.gz src/
	mv traefik-plugins/traefik-plugins.tar.gz cmd/controller-manager/

##---------------------------------------------------------------------
##@ Docker build
##---------------------------------------------------------------------
.PHONY: docker-build
docker-build: docker-build-manager docker-build-dashboard ## Build the docker image.

.PHONY: docker-build-manager
docker-build-manager: test ## Build the docker image for controller-manager.
	DOCKER_BUILDKIT=1 docker build . -t ${IMG_MANAGER} -f dockerfile/controller-manager.Dockerfile

.PHONY: docker-build-dashboard
docker-build-dashboard: test ## Build the docker image for dashboard.
	DOCKER_BUILDKIT=1 docker build . -t ${IMG_DASHBOARD} -f dockerfile/dashboard.Dockerfile

.PHONY: docker-push docker-push-manager docker-push-dashboard
docker-push: docker-push-manager docker-push-dashboard ## Build the docker image.

REGISTORY ?= ghcr.io/cosmo-workspace

docker-push-manager: docker-build-manager ## push cosmo contoller-manager image.
	docker tag ${IMG_MANAGER} ${REGISTORY}/${IMG_MANAGER}
	docker push ${REGISTORY}/${IMG_MANAGER}

docker-push-dashboard: docker-build-dashboard ## push cosmo dashboard image.
	docker tag ${IMG_DASHBOARD} ${REGISTORY}/${IMG_DASHBOARD}
	docker push ${REGISTORY}/${IMG_DASHBOARD}

##---------------------------------------------------------------------
##@ Deployment
##---------------------------------------------------------------------
ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: install
install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

.PHONY: uninstall
uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/crd | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

.PHONY: deploy
deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/default | kubectl apply -f -

.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/default | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

##---------------------------------------------------------------------
##@ Build Dependencies
##---------------------------------------------------------------------

## Tool Versions
GO_VERSION ?= 1.22.3
KUSTOMIZE_VERSION ?= v5.4.1
CONTROLLER_TOOLS_VERSION ?= v0.15.0

# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION ?= 1.30.x

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

export PATH := $(LOCALBIN):$(PATH)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif
GO ?= $(GOBIN)/go$(GO_VERSION)

## Tool Binaries
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest
HELM ?= $(LOCALBIN)/helm

KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
.PHONY: kustomize
kustomize: $(LOCALBIN) $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE):
	curl -s $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN)

.PHONY: controller-gen
controller-gen: go $(LOCALBIN) $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN):
	GOBIN=$(LOCALBIN) $(GO) install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: envtest
envtest: go $(LOCALBIN) $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST):
	GOBIN=$(LOCALBIN) $(GO) install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

.PHONY: helm
helm: $(LOCALBIN) $(HELM) ## Download helm locally if necessary.
$(HELM):
	export HELM_INSTALL_DIR=$(LOCALBIN) && \
	curl -s curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

.PHONY: go
go: $(GO)
$(GO): 
	go install golang.org/dl/go$(GO_VERSION)@latest
	$(GO) download

.PHONY: configure
configure: kustomize controller-gen envtest
