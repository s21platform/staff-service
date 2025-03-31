#!/bin/bash
protoc  --go_out=./  \
        --go-grpc_out=./ \
        api/v0/staff.v0.proto

protoc --doc_out=. --doc_opt=markdown,README.md api/v0/staff.v0.proto