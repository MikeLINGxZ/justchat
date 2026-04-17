"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const rpc_1 = require("./rpc");
const host_1 = require("./host");
// Ensure stdout is used only for RPC, all logs go to stderr
const rpc = new rpc_1.JsonRpcConnection(process.stdin, process.stdout);
const host = new host_1.ExtensionHost(rpc);
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
//# sourceMappingURL=index.js.map