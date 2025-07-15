#!/bin/bash

# proto 文件根目录
PROTO_DIR="lib/rpc"

# 1. 删除 rpc 目录下的所有 .dart 文件
echo "Cleaning existing .dart files under $PROTO_DIR..."
find "$PROTO_DIR" -type f -name "*.dart" | while read file; do
  echo "Removing $file"
  rm -f "$file"
done
echo "Cleanup completed."

# 2. 查找并生成 proto 的 dart 代码
echo "Generating Dart code from .proto files..."
find "$PROTO_DIR" -type f -name "*.proto" | while read file; do
  echo "Processing $file"
  protoc --dart_out=. "$file"
done

echo "All proto files processed."