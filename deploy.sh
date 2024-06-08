#!/bin/bash
set -e

workdir=${workdir:-$PWD}

# be test
# docker build -f ${workdir}/tools/test_be.Dockerfile ${workdir} -t be_test

# docker-compose to serup postgresql, minio, services
