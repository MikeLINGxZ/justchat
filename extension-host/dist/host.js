"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
exports.ExtensionHost = void 0;
const path = __importStar(require("path"));
const api_1 = require("./api");
class ExtensionHost {
    constructor(rpc) {
        this.plugins = new Map();
        this.rpc = rpc;
        this.pluginsDir = process.env.LEMONTEA_PLUGINS_DIR || '';
    }
    registerHandlers() {
        this.rpc.onRequest('plugin/activate', async (params) => {
            await this.activatePlugin(params.pluginId, params.pluginDir);
            return { success: true };
        });
        this.rpc.onRequest('plugin/deactivate', async (params) => {
            await this.deactivatePlugin(params.pluginId);
            return { success: true };
        });
        this.rpc.onRequest('tool/execute', async (params) => {
            return await this.executeTool(params.pluginId, params.toolId, params.input);
        });
        this.rpc.onRequest('hook/onBeforeChat', async (params) => {
            return await this.runHook(params.pluginId, 'beforeChat', params.context);
        });
        this.rpc.onRequest('hook/onAfterChat', async (params) => {
            return await this.runHook(params.pluginId, 'afterChat', params.context);
        });
        this.rpc.onRequest('shutdown', async () => {
            await this.shutdownAll();
            process.exit(0);
        });
        // Handle UI messages from frontend → plugin
        this.rpc.onRequest('ui/message', async (params) => {
            const plugin = this.plugins.get(params.pluginId);
            if (!plugin)
                return;
            const handlers = plugin.messageHandlers.get(params.viewId);
            if (handlers) {
                for (const handler of handlers) {
                    try {
                        handler(params.message);
                    }
                    catch (err) {
                        process.stderr.write(`[extension-host] UI message handler error in ${params.pluginId}: ${err}\n`);
                    }
                }
            }
        });
    }
    async activatePlugin(pluginId, pluginDir) {
        if (this.plugins.has(pluginId)) {
            await this.deactivatePlugin(pluginId);
        }
        const plugin = {
            id: pluginId,
            dir: pluginDir,
            module: null,
            tools: new Map(),
            beforeChatHooks: [],
            afterChatHooks: [],
            messageHandlers: new Map(),
        };
        const api = (0, api_1.createPluginAPI)(plugin, this.rpc);
        // Load the plugin module
        const pkgJsonPath = path.join(pluginDir, 'package.json');
        const pkgJson = require(pkgJsonPath);
        const mainPath = path.resolve(pluginDir, pkgJson.main || 'index.js');
        // Clear require cache for hot reload
        delete require.cache[require.resolve(mainPath)];
        const mod = require(mainPath);
        plugin.module = mod;
        this.plugins.set(pluginId, plugin);
        if (typeof mod.activate === 'function') {
            await mod.activate(api);
        }
        process.stderr.write(`[extension-host] Plugin activated: ${pluginId}\n`);
    }
    async deactivatePlugin(pluginId) {
        const plugin = this.plugins.get(pluginId);
        if (!plugin)
            return;
        if (plugin.module && typeof plugin.module.deactivate === 'function') {
            try {
                await plugin.module.deactivate();
            }
            catch (err) {
                process.stderr.write(`[extension-host] Error deactivating ${pluginId}: ${err}\n`);
            }
        }
        this.plugins.delete(pluginId);
        process.stderr.write(`[extension-host] Plugin deactivated: ${pluginId}\n`);
    }
    async executeTool(pluginId, toolId, input) {
        const plugin = this.plugins.get(pluginId);
        if (!plugin)
            throw new Error(`Plugin not found: ${pluginId}`);
        const tool = plugin.tools.get(toolId);
        if (!tool)
            throw new Error(`Tool not found: ${toolId} in plugin ${pluginId}`);
        return await tool.execute(input);
    }
    async runHook(pluginId, hookType, context) {
        const plugin = this.plugins.get(pluginId);
        if (!plugin)
            return context;
        const hooks = hookType === 'beforeChat' ? plugin.beforeChatHooks : plugin.afterChatHooks;
        let ctx = context;
        for (const hook of hooks) {
            try {
                ctx = await hook(ctx);
            }
            catch (err) {
                process.stderr.write(`[extension-host] Hook error in ${pluginId}/${hookType}: ${err}\n`);
            }
        }
        return ctx;
    }
    async shutdownAll() {
        for (const [pluginId] of this.plugins) {
            await this.deactivatePlugin(pluginId);
        }
    }
}
exports.ExtensionHost = ExtensionHost;
//# sourceMappingURL=host.js.map