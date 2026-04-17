"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.createPluginAPI = createPluginAPI;
function createPluginAPI(plugin, rpc) {
    return {
        tools: {
            register(tool) {
                plugin.tools.set(tool.id, tool);
                // Notify Go side about the tool registration (schema info)
                rpc.notify('app/registerTool', {
                    pluginId: plugin.id,
                    toolId: tool.id,
                    description: tool.description,
                    parameters: tool.parameters,
                });
            },
        },
        agents: {
            register(agent) {
                rpc.notify('app/registerAgent', {
                    pluginId: plugin.id,
                    agentId: agent.id,
                    name: agent.name,
                    description: agent.description,
                    systemPrompt: agent.systemPrompt,
                    tools: agent.tools,
                    role: agent.role || 'worker',
                });
            },
        },
        hooks: {
            onBeforeChat(handler) {
                plugin.beforeChatHooks.push(handler);
                return {
                    dispose() {
                        const idx = plugin.beforeChatHooks.indexOf(handler);
                        if (idx >= 0)
                            plugin.beforeChatHooks.splice(idx, 1);
                    },
                };
            },
            onAfterChat(handler) {
                plugin.afterChatHooks.push(handler);
                return {
                    dispose() {
                        const idx = plugin.afterChatHooks.indexOf(handler);
                        if (idx >= 0)
                            plugin.afterChatHooks.splice(idx, 1);
                    },
                };
            },
        },
        ui: {
            postMessage(viewId, data) {
                rpc.notify('app/emitEvent', {
                    event: `ui:message:${plugin.id}:${viewId}`,
                    data,
                });
            },
            onMessage(viewId, handler) {
                if (!plugin.messageHandlers.has(viewId)) {
                    plugin.messageHandlers.set(viewId, []);
                }
                plugin.messageHandlers.get(viewId).push(handler);
                return {
                    dispose() {
                        const handlers = plugin.messageHandlers.get(viewId);
                        if (handlers) {
                            const idx = handlers.indexOf(handler);
                            if (idx >= 0)
                                handlers.splice(idx, 1);
                        }
                    },
                };
            },
            renderChatCard(cardId, data) {
                rpc.notify('app/emitEvent', {
                    event: 'chatCard',
                    data: { pluginId: plugin.id, cardId, data },
                });
            },
        },
        storage: {
            async get(key) {
                const result = await rpc.call('app/storage/get', { pluginId: plugin.id, key });
                if (result == null)
                    return null;
                try {
                    return typeof result === 'string' ? JSON.parse(result) : result;
                }
                catch {
                    return result;
                }
            },
            async set(key, value) {
                const encoded = JSON.stringify(value);
                await rpc.call('app/storage/set', { pluginId: plugin.id, key, value: encoded });
            },
            async delete(key) {
                await rpc.call('app/storage/delete', { pluginId: plugin.id, key });
            },
        },
    };
}
//# sourceMappingURL=api.js.map