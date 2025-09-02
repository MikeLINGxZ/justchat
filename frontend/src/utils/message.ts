
// 类型转换工具函数
import type {CommonMessage} from "@/api/service/chat/Api.ts";
import type {Message} from "@/types";

export const ConvertApiMessageToMessage = (apiMessage: CommonMessage): Message => {
    return {
        id: `${apiMessage.chatUuid || 'unknown'}-${Date.now()}-${Math.random()}`,
        role: apiMessage.role || 'user',
        content: apiMessage.content || '',
        reasoningContent: apiMessage.reasoningContent, // 添加思考过程字段映射
        timestamp: Date.now(),
    };
};

export const ConvertMessageToApiMessage = (message: Message): CommonMessage => {
    return {
        role: message.role,
        content: message.content,
        chatUuid: message.id,
    };
};