#!/bin/bash

PROTO_DIR="./proto/pod"  # proto文件目录


# 生成Go代码
protoc --go_out=. --micro_out=. --go_opt=paths=source_relative --micro_opt=paths=source_relative "${PROTO_DIR}"/*.proto
