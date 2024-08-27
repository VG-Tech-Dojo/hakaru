#!/usr/bin/env bash
cd /root/hakaru || exit 2
make deploy ARTIFACTS_COMMIT=latest
