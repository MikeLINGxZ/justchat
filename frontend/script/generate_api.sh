#!/bin/bash
npx swagger-typescript-api generate --path ./rpc/common/common.swagger.json --output ./src/api/common/
npx swagger-typescript-api generate --path ./rpc/service/auth.swagger.json --output ./src/api/service/auth/
npx swagger-typescript-api generate --path ./rpc/service/chat.swagger.json --output ./src/api/service/chat/
npx swagger-typescript-api generate --path ./rpc/service/models.swagger.json --output ./src/api/service/models/