VER=`git tag --points-at HEAD | head -n 1`
USER=`id -u -n`
BUILDDATE=`date`

build:
	go build -ldflags="-X 'github.com/criteo/data-aggregation-api/internal/version.version=$(VER)' \
                   -X 'github.com/criteo/data-aggregation-api/internal/version.buildUser=$(USER)' \
                   -X 'github.com/criteo/data-aggregation-api/internal/version.buildTime=$(BUILDDATE)'" \
                  -o .build/data-aggregation-api ./cmd/data-aggregation-api

run:
	go run -ldflags="-X 'github.com/criteo/data-aggregation-api/internal/version.version=$(VER)' \
                 -X 'github.com/criteo/data-aggregation-api/internal/version.buildUser=$(USER)' \
                 -X 'github.com/criteo/data-aggregation-api/internal/version.buildTime=$(BUILDDATE)'" \
                ./cmd/data-aggregation-api

update_openconfig:
	./update_openconfig.sh