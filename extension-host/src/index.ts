import { JsonRpcConnection } from './rpc';
import { ExtensionHost } from './host';

// Ensure stdout is used only for RPC, all logs go to stderr
const rpc = new JsonRpcConnection(process.stdin, process.stdout);
const host = new ExtensionHost(rpc);

host.registerHandlers();

// Global error handlers
process.on('uncaughtException', (err) => {
  process.stderr.write(`[extension-host] Uncaught exception: ${err.stack || err.message}\n`);
  rpc.notify('plugin/error', { error: err.message });
});

process.on('unhandledRejection', (reason) => {
  process.stderr.write(`[extension-host] Unhandled rejection: ${reason}\n`);
  rpc.notify('plugin/error', { error: String(reason) });
});

process.stderr.write('[extension-host] Extension Host started\n');
