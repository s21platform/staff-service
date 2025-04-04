#!/bin/bash
protoc  --go_out=./  \
        --go-grpc_out=./ \
        api/staff.proto

protoc --doc_out=. --doc_opt=markdown,README.md api/staff.proto