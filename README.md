[![Go](https://img.shields.io/github/go-mod/go-version/criteo/data-aggregation-api)](https://github.com/criteo/data-aggregation-api)
[![status](https://img.shields.io/badge/status-alpha-orange)](https://github.com/criteo/data-aggregation-api)
[![CI](https://github.com/criteo/data-aggregation-api/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/criteo/data-aggregation-api/actions/workflows/ci.yml)
[![GitHub](https://img.shields.io/github/license/criteo/data-aggregation-api)](https://github.com/criteo/data-aggregation-api/blob/main/LICENSE)

# Criteo Data Aggregation API

This API aggregates data from their sources of truth: the Network CMDB or possibly any other data source you may have.

Then, it computes this data to provide OpenConfig JSON for each device as an output.

[ygot](https://github.com/openconfig/ygot) is used to validate the output against the OpenConfig YANG models.


## HowTo: add new config namespace/option

Note: in future release, we aim to simplify the integration as much as possible.

You need to update the following part in the code:

1. add your ingestor in `ingestors/cmdb/<yournewingestor>.go`:
   - GetBGPGlobal()
   - PrecomputeBGPGlobal()

2. register your ingestor in `internal/ingestors/repository.go`:
   - DataPerDevice struct
   - IngestorRepository struct
   - FetchAll()
   - Precompute()

3. store the preprocessed ingestor data into `internal/convertors/device/device.go`:
   - Device struct
   - NewDevice()

4. add your convertor in `internal/convertors/...`

3. execute your convertor (`internal/convertors/device/device.go`):
   - GenerateOpenconfig()
