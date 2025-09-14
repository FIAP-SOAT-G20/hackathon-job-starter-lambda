.DEFAULT_GOAL := help

export PATH := $(shell go env GOPATH)/bin:$(PATH)


# Variables
LAMBDA_DIR     = .
BINARY_NAME    = bin/bootstrap
ZIP_NAME       = dist/function.zip
BUILD_OS       = linux
BUILD_ARCH     = amd64
MAIN_FILE=main.go

VERSION=$(shell git describe --tags --always --dirty)
TEST_PATH=./internal/...
TEST_COVERAGE_FILE_NAME=coverage.out
MIGRATION_PATH = internal/infrastructure/database/migrations
DB_URL = postgres://postgres:postgres@localhost:5432/fastfood_10soat_g18_tc2?sslmode=disable
LAMBDA_INPUT_FILE=test/data/s3_event_payload.json

# Go commands
AWSLAMBDARPCCMD ?= awslambdarpc
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) $MAIN_FILE
GOTEST=ENVIRONMENT=test $(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
GOTIDY=$(GOCMD) mod tidy
SHCMD=sh

# Looks at comments using ## on targets and uses them to produce a help output.
.PHONY: help
help: ALIGN=22
help: ## 📜 Print this message
	@echo "Usage: make <command>"
	@awk -F '::? .*## ' -- "/^[^':]+::? .*## /"' { printf "  make '$$(tput bold)'%-$(ALIGN)s'$$(tput sgr0)' - %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo

.PHONY: fmt
fmt: ## 🗂️  Format the code
	@echo  "🟢 Formatting the code..."
	$(GOCMD) fmt ./...
	@echo

.PHONY: build
build: fmt ## 🔨 Build the application
	@echo  "🟢 Building the application..."
	#$(GOBUILD) -v -gcflags='all=-N -l' -o bin/$(APP_NAME) $(MAIN_FILE)
	GOOS=$(BUILD_OS) GOARCH=$(BUILD_ARCH) $(GOBUILD) -ldflags="-s -w" -o $(LAMBDA_DIR)/$(BINARY_NAME) $(LAMBDA_DIR)/cmd/job/processor/main.go
	@echo


.PHONY: package
package: build ## 📦 Package the binary into a .zip file for Lambda deployment
	@echo "📦 Packaging Lambda binary into zip..."
	mkdir -p ./dist
	zip -j $(LAMBDA_DIR)/$(ZIP_NAME) $(LAMBDA_DIR)/$(BINARY_NAME)
	@echo

.PHONY: start-lambda
start-lambda:  build  ## ▶  Start the lambda application locally to prepare to receive requests
	@echo "🟢 Starting lambda ..."
	_LAMBDA_SERVER_PORT=3300 AWS_LAMBDA_RUNTIME_API=http://localhost:3300 $(GOCMD) run $(LAMBDA_DIR)/cmd/job/processor/main.go
	@echo

.PHONY: trigger-lambda
trigger-lambda: ## ⚡  Trigger lambda with the input file stored in variable $LAMBDA_INPUT_FILE
	@echo "🟢 Triggering lambda with event: $(LAMBDA_INPUT_FILE)"
	@PATH="$(shell go env GOPATH)/bin:$$PATH" \
		'$(AWSLAMBDARPCCMD)' -a localhost:3300 -e $(LAMBDA_INPUT_FILE)
	@echo


.PHONY: test
test: lint ## 🧪 Run tests
	@echo  "🟢 Running tests..."
	@$(GOFMT) ./...
	@$(GOVET) ./...
	@$(GOTIDY)
	$(GOTEST) $(TEST_PATH) -race -v
	@echo

.PHONY: coverage
coverage: ## 🧪 Run tests with coverage
	@echo  "🟢 Running tests with coverage..."
# remove files that are not meant to be tested
	$(GOTEST) $(TEST_PATH) -coverprofile=$(TEST_COVERAGE_FILE_NAME).tmp
	@cat $(TEST_COVERAGE_FILE_NAME).tmp | grep -v "_mock.go" | grep -v "_request.go" | grep -v "_response.go" \
	| grep -v "_gateway.go" | grep -v "_datasource.go" | grep -v "_presenter.go" | grep -v "middleware" \
	| grep -v "config" | grep -v "route" | grep -v "util" | grep -v "database" \
	| grep -v "server" | grep -v "logger" | grep -v "httpclient" > $(TEST_COVERAGE_FILE_NAME)
	@rm $(TEST_COVERAGE_FILE_NAME).tmp
	$(GOCMD) tool cover -html=$(TEST_COVERAGE_FILE_NAME)
	@echo

.PHONY: clean
clean: ## 🧹 Clean up binaries and coverage files
	@echo "🔴 Cleaning up..."
	$(GOCLEAN)
	rm -f $(APP_NAME)
	rm -f $(TEST_COVERAGE_FILE_NAME)
	rm -f $(LAMBDA_DIR)/$(BINARY_NAME) $(LAMBDA_DIR)/$(ZIP_NAME)
	@echo


.PHONY: lint
lint: ## 🔍 Run linter
	@echo "🟢 Running linter..."
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.7 run --out-format colored-line-number
	@echo

.PHONY: migrate-create
migrate-create: ## 🔄 Create new migration, usage example: make migrate-create name=create_table_products
	@echo "🟢 Creating new migration..."
# if name is not passed, required argument
ifndef name
	$(error name is not set, usage example: make migrate-create name=create_table_products)
endif
	migrate create -ext sql -dir ${MIGRATION_PATH} -seq $(name)
	@echo

.PHONY: migrate-up
migrate-up: ## ⬆️  Run migrations
	@echo "🟢 Running migrations..."
	migrate -path ${MIGRATION_PATH} -database "${DB_URL}" -verbose up
	@echo

.PHONY: migrate-down
migrate-down: ## ⬇️  Roll back migrations
	@echo "🔴 Rolling back migrations..."
	migrate -path ${MIGRATION_PATH} -database "${DB_URL}" -verbose down
	@echo

.PHONY: install
install: ## 📦 Install dependencies
	@echo "🟢 Installing dependencies..."
	go mod download
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.2
	@go install github.com/blmayer/awslambdarpc@latest
	@echo

.PHONY: compose-up
compose-up: ## ▶  Start local database with docker compose
	@echo "🟢 Starting development environment..."
	docker compose pull
	docker-compose up -d --wait --build
	@echo

.PHONY: compose-down
compose-down: ## ■  Stops local database with docker compose
	@echo "🔴 Stopping development environment..."
	docker-compose down
	@echo

.PHONY: compose-clean
compose-clean: ## 🧹 Clean the application with docker compose, removing volumes and images
	@echo "🔴 Cleaning the application..."
	docker compose down --volumes --rmi all
	@echo

.PHONY: scan
scan: ## 🔍 Run security scan
	@echo  "🟠 Running security scan..."
	@go run golang.org/x/vuln/cmd/govulncheck@v1.1.4 -show verbose ./...
	@go run github.com/aquasecurity/trivy/cmd/trivy@latest image --severity HIGH,CRITICAL $(DOCKER_REGISTRY)/$(DOCKER_REGISTRY_APP):latest
	@echo


.PHONY: terraform-init
terraform-init: ## 🔧 Initialize Terraform loading state from AWS S3 bucket
	@echo "🟢 Initializing Terraform..."
	cd terraform && terraform init -force-copy
	@echo

.PHONY: terraform-plan
terraform-plan: ## 💭 Plan Terraform
	@echo "🟢 Planning Terraform..."
	cd terraform && terraform plan -var-file=production.tfvars
	@echo	

.PHONY: terraform-apply
terraform-apply: ## ⚡ Apply Terraform
	@echo "🟢 Applying Terraform..."
	cd terraform && terraform apply -var-file=production.tfvars -auto-approve
	@echo

.PHONY: terraform-destroy
terraform-destroy: ## ⚠️  Destroy Terraform
	@echo "🔴 Destroying Terraform..."
	cd terraform && terraform destroy -var-file=production.tfvars -auto-approve
	cd ..
	@echo

.PHONY: deploy-observability
deploy-observability: ## 🚀 Deploy Elasticsearch stack to observability namespace
	@echo "Deploying Elasticsearch stack..."
	kubectl apply -f k8s/observability/namespace.yaml
	kubectl apply -f k8s/observability/elasticsearch.yaml
	kubectl apply -f k8s/observability/kibana.yaml
	@echo "Waiting for Elasticsearch to be ready..."
	kubectl wait --for=condition=available --timeout=300s deployment/elasticsearch -n observability
	@echo "Elasticsearch stack deployed!"

.PHONY: deploy-fluent-bit
deploy-fluent-bit: ## 📊 Deploy Fluent Bit for log collection
	@echo "Deploying Fluent Bit..."
	kubectl apply -f k8s/observability/fluent-bit-rbac.yaml
	kubectl apply -f k8s/observability/fluent-bit-config.yaml
	kubectl apply -f k8s/observability/fluent-bit-ds.yaml
	@echo "Fluent Bit deployed!"

.PHONY: deploy-logging
deploy-logging: deploy-observability deploy-fluent-bit ## 🎯 Deploy complete logging stack

.PHONY: port-forward-kibana
port-forward-kibana: ## 🌐 Port forward Kibana to localhost
	@echo "Port forwarding Kibana to localhost:5601..."
	kubectl port-forward service/kibana 5601:5601 -n observability

.PHONY: check-elasticsearch
check-elasticsearch: ## 🔍 Check Elasticsearch health
	@echo "Checking Elasticsearch health..."
	kubectl exec -n observability deployment/elasticsearch -- curl -s http://localhost:9200/_cluster/health | jq .

.PHONY: check-logs
check-logs: ## 📋 Check if logs are being indexed
	@echo "Checking for indexed logs..."
	kubectl exec -n observability deployment/elasticsearch -- curl -s "http://localhost:9200/minikube-logs/_search?pretty" | jq .

.PHONY: troubleshoot-elasticsearch
troubleshoot-elasticsearch: ## 🔧 Troubleshoot Elasticsearch issues
	@echo "Troubleshooting Elasticsearch..."
	@echo "Checking pod status..."
	kubectl get pods -n observability
	@echo ""
	@echo "Checking pod logs..."
	kubectl logs -n observability deployment/elasticsearch --tail=50
	@echo ""
	@echo "Checking pod events..."
	kubectl describe pod -n observability -l app=elasticsearch
	@echo ""
	@echo "Checking node resources..."
	kubectl top nodes
	@echo ""
	@echo "Checking if Minikube has enough memory..."
	minikube status

.PHONY: restart-elasticsearch
restart-elasticsearch: ## 🔄 Restart Elasticsearch deployment
	@echo "Restarting Elasticsearch..."
	kubectl rollout restart deployment/elasticsearch -n observability
	kubectl rollout status deployment/elasticsearch -n observability

.PHONY: delete-elasticsearch
delete-elasticsearch: ## 🗑️ Delete Elasticsearch deployment
	@echo "Deleting Elasticsearch deployment..."
	kubectl delete deployment elasticsearch -n observability
	kubectl delete service elasticsearch -n observability
	kubectl delete configmap elasticsearch-config -n observability

.PHONY: deploy-elasticsearch-light
deploy-elasticsearch-light: ## 🚀 Deploy lightweight Elasticsearch
	@echo "Deploying lightweight Elasticsearch..."
	kubectl apply -f k8s/observability/namespace.yaml
	kubectl apply -f k8s/observability/elasticsearch-light.yaml
	@echo "Waiting for Elasticsearch to be ready..."
	kubectl wait --for=condition=available --timeout=300s deployment/elasticsearch -n observability
	@echo "Lightweight Elasticsearch deployed!"

.PHONY: check-minikube-resources
check-minikube-resources: ## 📊 Check Minikube resource allocation
	@echo "Checking Minikube configuration..."
	@echo "Memory allocation:"
	minikube config get memory
	@echo "CPU allocation:"
	minikube config get cpus
	@echo "Disk size:"
	minikube config get disk-size

.PHONY: restart-observability
restart-observability: ## 🔄 Restart the entire observability stack
	@echo "Restarting observability stack..."
	kubectl rollout restart deployment/elasticsearch -n observability
	kubectl rollout restart deployment/kibana -n observability
	@echo "Waiting for deployments to be ready..."
	kubectl wait --for=condition=available --timeout=300s deployment/elasticsearch -n observability
	kubectl wait --for=condition=available --timeout=300s deployment/kibana -n observability
	@echo "Observability stack restarted!"

.PHONY: troubleshoot-kibana
troubleshoot-kibana: ## 🔧 Troubleshoot Kibana access issues
	@echo "Troubleshooting Kibana access..."
	@echo "Checking Kibana pod status..."
	kubectl get pods -n observability -l app=kibana
	@echo ""
	@echo "Checking Kibana logs..."
	kubectl logs -n observability deployment/kibana --tail=20
	@echo ""
	@echo "Checking Kibana service..."
	kubectl get svc kibana -n observability
	@echo ""
	@echo "Checking if Kibana is responding..."
	kubectl exec -n observability deployment/kibana -- curl -s http://localhost:5601/api/status || echo "Kibana not responding"
	@echo ""
	@echo "To access Kibana, try:"
	@echo "  make port-forward-kibana"
	@echo "  Then open http://localhost:5601 in your browser"

.PHONY: get-kibana-url
get-kibana-url: ## 🌐 Get Kibana access URL
	@echo "Getting Kibana access information..."
	@echo "NodePort:"
	kubectl get svc kibana -n observability -o jsonpath='{.spec.ports[0].nodePort}'
	@echo ""
	@echo "To access Kibana:"
	@echo "1. Port forward: make port-forward-kibana"
	@echo "2. Or use NodePort: http://$(minikube ip):$(kubectl get svc kibana -n observability -o jsonpath='{.spec.ports[0].nodePort}')"

.PHONY: check-elasticsearch-kibana-compatibility
check-elasticsearch-kibana-compatibility: ## 🔍 Check version compatibility
	@echo "Checking Elasticsearch and Kibana versions..."
	@echo "Elasticsearch version:"
	kubectl exec -n observability deployment/elasticsearch -- curl -s http://localhost:9200 | jq -r '.version.number' 2>/dev/null || echo "Cannot get Elasticsearch version"
	@echo "Kibana version:"
	kubectl exec -n observability deployment/kibana -- curl -s http://localhost:5601/api/status | jq -r '.version.number' 2>/dev/null || echo "Cannot get Kibana version"