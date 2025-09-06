package example
package example

import (
	"context"
	"log"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/global"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/rpc/service"
)

// ExampleUsage 展示如何使用新的 gRPC 客户端
func ExampleUsage() {
	ctx := context.Background()

	// 1. 使用默认配置
	log.Println("使用默认配置的 gRPC 客户端")
	
	// 获取认证服务客户端（自动连接）
	authClient, err := global.GRPC.Auth(ctx)
	if err != nil {
		log.Printf("获取认证客户端失败: %v", err)
		return
	}

	// 调用登录接口
	loginReq := &service.LoginRequest{
		LoginField:  "test@example.com",
		PasswordMd5: "test_password_hash",
	}
	
	loginResp, err := authClient.Login(ctx, loginReq)
	if err != nil {
		log.Printf("登录失败: %v", err)
	} else {
		log.Printf("登录成功: %+v", loginResp)
	}

	// 2. 使用自定义配置
	log.Println("使用自定义配置")
	
	customConfig := &global.GrpcClientConfig{
		Host:           "192.168.1.100",
		Port:           8080,
		ConnectTimeout: 15 * time.Second,
		KeepAlive:      60 * time.Second,
		MaxRetries:     5,
		Insecure:       true,
	}
	
	// 应用新配置（会自动重新连接）
	global.GRPC.SetConfig(customConfig)
	
	// 获取聊天服务客户端
	chatClient, err := global.GRPC.Chat(ctx)
	if err != nil {
		log.Printf("获取聊天客户端失败: %v", err)
		return
	}

	// 调用聊天接口
	chatReq := &service.ListChatsRequest{
		// 填入相应的请求参数
	}
	
	chatResp, err := chatClient.ListChats(ctx, chatReq)
	if err != nil {
		log.Printf("获取聊天列表失败: %v", err)
	} else {
		log.Printf("获取聊天列表成功: %+v", chatResp)
	}

	// 3. 检查连接状态
	log.Printf("连接状态: %s", global.GRPC.GetConnectionState())
	log.Printf("是否已连接: %v", global.GRPC.IsConnected())

	// 4. 获取模型服务客户端
	modelsClient, err := global.GRPC.Models(ctx)
	if err != nil {
		log.Printf("获取模型客户端失败: %v", err)
		return
	}

	// 调用模型接口
	modelsReq := &service.ModelsRequest{}
	modelsResp, err := modelsClient.Models(ctx, modelsReq)
	if err != nil {
		log.Printf("获取模型列表失败: %v", err)
	} else {
		log.Printf("获取模型列表成功: %+v", modelsResp)
	}
}

// ExampleStreamingChat 展示如何使用流式聊天功能
func ExampleStreamingChat() {
	ctx := context.Background()

	// 获取聊天客户端
	chatClient, err := global.GRPC.Chat(ctx)
	if err != nil {
		log.Printf("获取聊天客户端失败: %v", err)
		return
	}

	// 创建流式聊天请求
	completionsReq := &service.CompletionsRequest{
		// 填入相应的请求参数
		// Model: "gpt-3.5-turbo",
		// Messages: []*service.Message{...},
	}

	// 调用流式聊天接口
	stream, err := chatClient.Completions(ctx, completionsReq)
	if err != nil {
		log.Printf("创建聊天流失败: %v", err)
		return
	}

	// 接收流式响应
	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Printf("接收流式响应结束或出错: %v", err)
			break
		}
		
		log.Printf("收到流式响应: %+v", resp)
	}
}

// ExampleGracefulShutdown 展示如何优雅关闭 gRPC 客户端
func ExampleGracefulShutdown() {
	// 在应用关闭时调用
	if err := global.GRPC.Close(); err != nil {
		log.Printf("关闭 gRPC 客户端失败: %v", err)
	} else {
		log.Println("gRPC 客户端已优雅关闭")
	}
}