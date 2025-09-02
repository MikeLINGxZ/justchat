# 认证功能使用说明

## 概述

本项目已更新为使用生成的RPC接口进行认证功能。主要变更包括：

1. 使用生成的TypeScript类型定义
2. 使用生成的RPC客户端进行API调用
3. 密码使用MD5加密
4. 支持邮箱验证码注册

## 主要文件

### 1. 生成的类型定义 (`src/generated/types.ts`)
- `LoginRequest`: 登录请求类型
- `LoginResponse`: 登录响应类型
- `RegisterRequest`: 注册请求类型
- `SendEmailVerificationCodeRequest`: 发送验证码请求类型

### 2. 生成的RPC客户端 (`src/generated/client.ts`)
- `getAuthClient()`: 获取认证服务客户端
- 提供登录、注册、发送验证码、登出等方法

### 3. 认证状态管理 (`src/stores/authStore.ts`)
- 使用Zustand管理认证状态
- 集成生成的RPC客户端
- 自动处理token存储和清除

### 4. 工具函数 (`src/utils/crypto.ts`)
- `hashPassword()`: 密码MD5加密
- `validatePassword()`: 密码强度验证

## 使用方法

### 登录功能

```typescript
import { useAuthStore } from '@/stores/authStore';
import { hashPassword } from '@/utils/crypto';

const { login } = useAuthStore();

// 登录
const handleLogin = async (username: string, password: string) => {
  try {
    await login({
      username,
      password_md5: hashPassword(password)
    });
    // 登录成功，自动跳转
  } catch (error) {
    // 处理错误
  }
};
```

### 注册功能

```typescript
import { useAuthStore } from '@/stores/authStore';
import { getAuthClient } from '@/generated/client';
import { hashPassword } from '@/utils/crypto';

const { register } = useAuthStore();
const authClient = getAuthClient();

// 发送验证码
const handleSendCode = async (email: string, username: string) => {
  await authClient.sendEmailVerificationCode({
    email,
    username,
    code_type: 1 // 注册验证码
  });
};

// 注册
const handleRegister = async (userData: {
  username: string;
  email: string;
  password: string;
  emailVerificationCode: number;
}) => {
  try {
    await register({
      username: userData.username,
      email: userData.email,
      password_md5: hashPassword(userData.password),
      email_verification_code: userData.emailVerificationCode
    });
    // 注册成功，自动跳转
  } catch (error) {
    // 处理错误
  }
};
```

### 登出功能

```typescript
import { useAuthStore } from '@/stores/authStore';

const { logout } = useAuthStore();

const handleLogout = async () => {
  await logout();
  // 自动清除token和用户信息
};
```

## 页面组件

### 登录页面 (`src/pages/auth/Login.tsx`)
- 用户名/密码登录
- 自动MD5加密
- 错误处理和重定向

### 注册页面 (`src/pages/auth/Register.tsx`)
- 用户名/邮箱/密码注册
- 邮箱验证码功能
- 密码强度验证

### 测试页面 (`src/pages/TestAuth.tsx`)
- API功能测试
- 状态显示
- 开发调试用

## 配置说明

### RPC客户端配置
## 环境变量配置

项目支持通过环境变量配置后端API地址和相关参数：

### 配置文件

1. 复制 `.env.example` 文件为 `.env`：
```bash
cp .env.example .env
```

2. 修改 `.env` 文件中的配置：
```bash
# 后端API服务器地址
VITE_API_BASE_URL=http://localhost:8080

# API请求超时时间（毫秒）
VITE_API_TIMEOUT=30000

# 是否启用凭证传递
VITE_API_WITH_CREDENTIALS=true
```

### 环境变量说明

- `VITE_API_BASE_URL`: 后端API服务器地址
  - 开发环境: `http://localhost:8080`
  - 生产环境: `https://your-api-server.com`
- `VITE_API_TIMEOUT`: API请求超时时间（毫秒），默认30000
- `VITE_API_WITH_CREDENTIALS`: 是否启用凭证传递，默认true

### 测试配置

访问 `/env-test` 页面可以查看当前的环境变量配置状态。

### 代码中使用

```typescript
import { getApiConfig, env } from '@/config/env';

// 获取API配置
const apiConfig = getApiConfig();
console.log(apiConfig.baseUrl); // 输出配置的API地址

// 获取环境信息
console.log(env.isDev); // 是否为开发模式
```

## API客户端配置

API客户端会自动读取环境变量配置：

### 认证拦截器
RPC客户端已配置：
- 请求拦截器：自动添加Authorization头
- 响应拦截器：401错误自动跳转登录页

## 测试

1. 启动开发服务器：`npm run dev`
2. 访问测试页面：`http://localhost:5173/test-auth`
3. 测试各项认证功能

## 注意事项

1. 确保后端服务正在运行
2. 密码使用MD5加密，后端需要相应处理
3. 邮箱验证码功能需要后端支持
4. Token自动存储在localStorage中
5. 401错误会自动清除token并跳转登录页

## 错误处理

所有API调用都有统一的错误处理：
- 网络错误
- 服务器错误
- 业务逻辑错误
- 自动显示错误信息给用户
