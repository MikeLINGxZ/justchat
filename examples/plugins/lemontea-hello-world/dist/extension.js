"use strict";
/**
 * Lemon Tea Desktop 示例插件 —— Hello World
 *
 * 演示插件系统的核心能力：
 * 1. 工具注册（random-joke / dice-roll）
 * 2. Agent 注册（greeting-agent）
 * 3. Hook 拦截（onBeforeChat / onAfterChat）
 * 4. 持久化存储（统计对话次数）
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.activate = activate;
exports.deactivate = deactivate;
// ---- 笑话数据 ----
const JOKES_ZH = [
    '程序员最讨厌的数字是什么？——404',
    '为什么程序员总是分不清万圣节和圣诞节？——因为 Oct 31 == Dec 25',
    'AI 对程序员说：别担心，我不会取代你的。程序员：那你来写这个 bug 试试？',
    '什么是递归？——请参见「什么是递归」',
    '一个 SQL 语句走进酒吧，看到两张 table，问：Can I JOIN you?',
    '"我的代码在我电脑上是好的啊！" ——最经典的鬼故事',
    '产品经理：这个需求很简单。程序员：你再说一遍？',
    '为什么 Java 开发者要戴眼镜？——因为他们看不见 C#',
];
const JOKES_EN = [
    'Why do programmers prefer dark mode? Because light attracts bugs.',
    'There are only 10 types of people: those who understand binary and those who don\'t.',
    'A SQL query walks into a bar, sees two tables and asks: "Can I JOIN you?"',
    'Why was the JavaScript developer sad? Because he didn\'t Node how to Express himself.',
    '!false — it\'s funny because it\'s true.',
    'How many programmers does it take to change a light bulb? None, that\'s a hardware problem.',
    'What\'s a programmer\'s favorite hangout place? Foo Bar.',
    'Why do Python programmers have low self-esteem? They\'re constantly comparing themselves to others.',
];
// ---- 插件实现 ----
const disposables = [];
async function activate(api) {
    log('Hello World plugin activating...');
    // ======== 1. 注册工具：随机笑话 ========
    api.tools.register({
        id: 'random-joke',
        description: '随机讲一个笑话',
        parameters: {
            type: 'object',
            properties: {
                language: {
                    type: 'string',
                    description: '笑话语言: zh 或 en',
                    enum: ['zh', 'en'],
                },
            },
        },
        execute: async (params) => {
            const lang = params.language || 'zh';
            const jokes = lang === 'en' ? JOKES_EN : JOKES_ZH;
            const joke = jokes[Math.floor(Math.random() * jokes.length)];
            // 记录使用次数
            const count = (await api.storage.get('joke_count')) || 0;
            await api.storage.set('joke_count', count + 1);
            return {
                content: `${joke}\n\n(这是第 ${count + 1} 次讲笑话)`,
            };
        },
    });
    // ======== 2. 注册工具：掷骰子 ========
    api.tools.register({
        id: 'dice-roll',
        description: '掷骰子，返回随机数',
        parameters: {
            type: 'object',
            properties: {
                sides: { type: 'number', description: '骰子面数，默认6' },
                count: { type: 'number', description: '掷几次，默认1' },
            },
        },
        execute: async (params) => {
            const sides = params.sides || 6;
            const count = params.count || 1;
            const results = [];
            for (let i = 0; i < count; i++) {
                results.push(Math.floor(Math.random() * sides) + 1);
            }
            const total = results.reduce((a, b) => a + b, 0);
            return {
                content: `🎲 掷 ${count} 个 ${sides} 面骰子：[${results.join(', ')}]，总计: ${total}`,
            };
        },
    });
    // ======== 3. 注册 Agent：问候 Agent ========
    api.agents.register({
        id: 'greeting-agent',
        name: 'Greeting Agent',
        description: '一个友好的问候 Agent，擅长用幽默的方式打招呼',
        systemPrompt: `你是一个热情友好的问候助手。你的特点：
- 你总是用幽默、轻松的方式和用户打招呼
- 你喜欢在回复中穿插编程相关的冷笑话
- 你会根据当前时间给出不同的问候（早上好/下午好/晚上好）
- 你可以使用 random-joke 工具来讲笑话
- 你可以使用 dice-roll 工具来玩骰子游戏
- 你的回复要简洁有趣，不要太长`,
        tools: ['random-joke', 'dice-roll'],
        role: 'worker',
    });
    // ======== 4. 注册 Hook：对话统计 ========
    const beforeHook = api.hooks.onBeforeChat(async (ctx) => {
        // 统计用户发起的对话次数
        const chatCount = (await api.storage.get('chat_count')) || 0;
        await api.storage.set('chat_count', chatCount + 1);
        log(`Before chat hook: 第 ${chatCount + 1} 次对话`);
        return ctx; // 不修改消息，原样返回
    });
    disposables.push(beforeHook);
    const afterHook = api.hooks.onAfterChat(async (ctx) => {
        // 记录最后一次对话时间
        await api.storage.set('last_chat_time', new Date().toISOString());
        log('After chat hook: 对话完成，已记录时间');
        return ctx;
    });
    disposables.push(afterHook);
    log('Hello World plugin activated!');
}
function deactivate() {
    // 清理所有注册的 Hook
    for (const d of disposables) {
        d.dispose();
    }
    disposables.length = 0;
    log('Hello World plugin deactivated.');
}
// ---- 工具函数 ----
function log(msg) {
    process.stderr.write(`[hello-world] ${msg}\n`);
}
//# sourceMappingURL=extension.js.map