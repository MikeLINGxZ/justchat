import { JsonRpcConnection } from './rpc';
export declare class ExtensionHost {
    private plugins;
    private rpc;
    private pluginsDir;
    constructor(rpc: JsonRpcConnection);
    registerHandlers(): void;
    private activatePlugin;
    private deactivatePlugin;
    private executeTool;
    private runHook;
    private shutdownAll;
}
