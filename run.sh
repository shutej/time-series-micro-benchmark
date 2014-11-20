#!/bin/sh
echo time-series-micro-benchmark-mongodb
MONGODB_URI=mongodb://localhost time-series-micro-benchmark-mongodb
echo time-series-micro-benchmark-postgresql
POSTGRESQL_URI='dbname=tsmb sslmode=disable' time-series-micro-benchmark-postgresql
