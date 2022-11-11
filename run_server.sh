#!/bin/bash

go run ./main/tcp_server.go localhost:1234
# go run ./main/tcp_server.go localhost:1234 nodelay # 关闭nagle