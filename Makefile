# build app
.PHONY: build
build:
	@go build -o ./app ./

# build app alpine
.PHONY: build-alpine
build-alpine:
	@go mod tidy && \
	GOOS=linux GOARCH=amd64 go build -o ./app ./

.PHONY: mock-install
mock-install:
	@go install github.com/golang/mock/mockgen@v1.6.0

# make startProm
.PHONY: start-prom
start-prom:
	docker run -d \
	--rm \
	--network="host" \
	-p 9090:9090 \
	--name=prometheus \
	-v $(shell pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
	prom/prometheus

# make startGrafana
# for first timers, the username & password is both `admin`
.PHONY: start-grafana
start-grafana:
	docker volume create grafana-storage
	docker volume inspect grafana-storage
	docker run -p 3000:3000 --name=grafana grafana/grafana-oss || docker start grafana