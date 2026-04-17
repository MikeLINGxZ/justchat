import { JsonRpcConnection } from './rpc';
import { LemonTeaAPI, PluginInstance } from './types';
export declare function createPluginAPI(plugin: PluginInstance, rpc: JsonRpcConnection): LemonTeaAPI;
