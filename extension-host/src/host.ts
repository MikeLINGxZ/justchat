import * as path from 'path';
import { JsonRpcConnection } from './rpc';
import { createPluginAPI } from './api';
import { PluginInstance, PluginModule } from './types';

export class ExtensionHost {
  private plugins = new Map<string, PluginInstance>();
  private rpc: JsonRpcConnection;
  private pluginsDir: string;

  constructor(rpc: JsonRpcConnection) {
    this.rpc = rpc;
    this.pluginsDir = process.env.LEMONTEA_PLUGINS_DIR || '';
  }

  registerHandlers(): void {
    this.rpc.onRequest('plugin/activate', async (params: { pluginId: string; pluginDir: string }) => {
      await this.activatePlugin(params.pluginId, params.pluginDir);
      return { success: true };
    });

    this.rpc.onRequest('plugin/deactivate', async (params: { pluginId: string }) => {
      await this.deactivatePlugin(params.pluginId);
      return { success: true };
    });

    this.rpc.onRequest('tool/execute', async (params: { pluginId: string; toolId: string; input: any }) => {
      return await this.executeTool(params.pluginId, params.toolId, params.input);
    });

    this.rpc.onRequest('hook/onBeforeChat', async (params: { pluginId: string; context: any }) => {
      return await this.runHook(params.pluginId, 'beforeChat', params.context);
    });

    this.rpc.onRequest('hook/onAfterChat', async (params: { pluginId: string; context: any }) => {
      return await this.runHook(params.pluginId, 'afterChat', params.context);
    });

    this.rpc.onRequest('shutdown', async () => {
      await this.shutdownAll();
      process.exit(0);
    });

    // Handle UI messages from frontend → plugin
    this.rpc.onRequest('ui/message', async (params: { pluginId: string; viewId: string; message: any }) => {
      const plugin = this.plugins.get(params.pluginId);
      if (!plugin) return;
      const handlers = plugin.messageHandlers.get(params.viewId);
      if (handlers) {
        for (const handler of handlers) {
          try { handler(params.message); } catch (err) {
            process.stderr.write(`[extension-host] UI message handler error in ${params.pluginId}: ${err}\n`);
          }
        }
      }
    });
  }

  private async activatePlugin(pluginId: string, pluginDir: string): Promise<void> {
    if (this.plugins.has(pluginId)) {
      await this.deactivatePlugin(pluginId);
    }

    const plugin: PluginInstance = {
      id: pluginId,
      dir: pluginDir,
      module: null as any,
      tools: new Map(),
      beforeChatHooks: [],
      afterChatHooks: [],
      messageHandlers: new Map(),
    };

    const api = createPluginAPI(plugin, this.rpc);

    // Load the plugin module
    const pkgJsonPath = path.join(pluginDir, 'package.json');
    const pkgJson = require(pkgJsonPath);
    const mainPath = path.resolve(pluginDir, pkgJson.main || 'index.js');

    // Clear require cache for hot reload
    delete require.cache[require.resolve(mainPath)];
    const mod: PluginModule = require(mainPath);

    plugin.module = mod;
    this.plugins.set(pluginId, plugin);

    if (typeof mod.activate === 'function') {
      await mod.activate(api);
    }

    process.stderr.write(`[extension-host] Plugin activated: ${pluginId}\n`);
  }

  private async deactivatePlugin(pluginId: string): Promise<void> {
    const plugin = this.plugins.get(pluginId);
    if (!plugin) return;

    if (plugin.module && typeof plugin.module.deactivate === 'function') {
      try {
        await plugin.module.deactivate();
      } catch (err) {
        process.stderr.write(`[extension-host] Error deactivating ${pluginId}: ${err}\n`);
      }
    }

    this.plugins.delete(pluginId);
    process.stderr.write(`[extension-host] Plugin deactivated: ${pluginId}\n`);
  }

  private async executeTool(pluginId: string, toolId: string, input: any): Promise<{ content: string }> {
    const plugin = this.plugins.get(pluginId);
    if (!plugin) throw new Error(`Plugin not found: ${pluginId}`);

    const tool = plugin.tools.get(toolId);
    if (!tool) throw new Error(`Tool not found: ${toolId} in plugin ${pluginId}`);

    return await tool.execute(input);
  }

  private async runHook(pluginId: string, hookType: 'beforeChat' | 'afterChat', context: any): Promise<any> {
    const plugin = this.plugins.get(pluginId);
    if (!plugin) return context;

    const hooks = hookType === 'beforeChat' ? plugin.beforeChatHooks : plugin.afterChatHooks;
    let ctx = context;
    for (const hook of hooks) {
      try {
        ctx = await hook(ctx);
      } catch (err) {
        process.stderr.write(`[extension-host] Hook error in ${pluginId}/${hookType}: ${err}\n`);
      }
    }
    return ctx;
  }

  private async shutdownAll(): Promise<void> {
    for (const [pluginId] of this.plugins) {
      await this.deactivatePlugin(pluginId);
    }
  }
}
