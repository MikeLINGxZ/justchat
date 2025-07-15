#!/bin/bash

# 定义仓库URL和分支
REPO_URL="https://gitlab.linhf.cn/project/lemontea/lemon_tea_cli.git"
BRANCH="dev"
TARGET_DIR="rpc"
PROTO_PATTERN="*.proto"

# 创建临时目录
TEMP_DIR=$(mktemp -d)

echo "正在克隆仓库..."
git clone --depth 1 --branch "$BRANCH" "$REPO_URL" "$TEMP_DIR"

# 检查克隆是否成功
if [ $? -ne 0 ]; then
    echo "错误：无法克隆仓库"
    rm -rf "$TEMP_DIR"
    exit 1
fi

# 创建目标目录（如果不存在）
rm -rf "./$TARGET_DIR"
mkdir -p "./$TARGET_DIR"

echo "正在复制.proto文件..."
# 查找并复制所有.proto文件
find "$TEMP_DIR/$TARGET_DIR" -name "$PROTO_PATTERN" -exec cp {} "./$TARGET_DIR/" \;

# 检查复制是否成功
if [ $? -ne 0 ]; then
    echo "警告：未找到任何.proto文件或复制失败"
fi

# 清理临时目录
rm -rf "$TEMP_DIR"

echo "操作完成。.proto文件已保存到 ./$TARGET_DIR/"