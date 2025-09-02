#!/bin/bash

# 配置部分
REPO_URL="https://gitlab.linhf.cn/project/lemontea/lemon_tea_server.git"
BRANCH="dev"
RPC_DIR_IN_REPO="rpc"
LOCAL_RPC_DIR="rpc"

# 创建本地 rpc 目录（如果不存在）
mkdir -p "$LOCAL_RPC_DIR"

# 创建临时目录
TMP_DIR=$(mktemp -d)

echo "克隆仓库 $REPO_URL 的 $BRANCH 分支到临时目录..."
git clone --depth 1 -b "$BRANCH" "$REPO_URL" "$TMP_DIR"

if [ $? -ne 0 ]; then
  echo "❌ 克隆仓库失败，请检查网络或仓库地址。"
  exit 1
fi

# 查找并复制 .proto 文件，保留目录结构
PROTO_FILES=$(find "$TMP_DIR/$RPC_DIR_IN_REPO" -type f -name "*.proto")

if [ -z "$PROTO_FILES" ]; then
  echo "⚠️ 在 $RPC_DIR_IN_REPO 中未找到 .proto 文件。"
else
  echo "✅ 找到以下 .proto 文件："
  echo "$PROTO_FILES"

  # 使用 rsync 保留目录结构复制所有 .proto 文件
  rsync -av --include="*/" --include="*.proto" --exclude="*" "$TMP_DIR/$RPC_DIR_IN_REPO/" "$LOCAL_RPC_DIR/"
  echo "✅ .proto 文件已复制到 $LOCAL_RPC_DIR/，保留了原目录结构"
  
  # 替换 proto 文件中的字符串
  echo "🔄 正在替换 proto 文件中的字符串..."
  find "$LOCAL_RPC_DIR" -type f -name "*.proto" -exec sed -i '' 's/gitlab\.linhf\.cn\/project\/lemontea\/lemon_tea_server/gitlab.linhf.cn\/project\/lemontea\/lemon_tea_desktop/g' {} \;
  echo "✅ 字符串替换完成"
fi

# 清理临时目录
rm -rf "$TMP_DIR"

echo "✅ 完成！"