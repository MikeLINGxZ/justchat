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
exports.JsonRpcConnection = void 0;
const readline = __importStar(require("readline"));
class JsonRpcConnection {
    constructor(input, output) {
        this.nextId = 1;
        this.pending = new Map();
        this.handlers = new Map();
        this.output = output;
        this.rl = readline.createInterface({ input, crlfDelay: Infinity });
        this.rl.on('line', (line) => this.handleLine(line));
        this.rl.on('close', () => process.exit(0));
    }
    onRequest(method, handler) {
        this.handlers.set(method, handler);
    }
    async call(method, params) {
        const id = this.nextId++;
        const msg = { jsonrpc: '2.0', id, method, params };
        return new Promise((resolve, reject) => {
            this.pending.set(id, { resolve, reject });
            this.send(msg);
            // 30 second timeout
            setTimeout(() => {
                if (this.pending.has(id)) {
                    this.pending.delete(id);
                    reject(new Error(`RPC call ${method} timed out`));
                }
            }, 30000);
        });
    }
    notify(method, params) {
        this.send({ jsonrpc: '2.0', method, params });
    }
    send(msg) {
        this.output.write(JSON.stringify(msg) + '\n');
    }
    async handleLine(line) {
        if (!line.trim())
            return;
        try {
            const msg = JSON.parse(line);
            if (msg.method) {
                // Incoming request or notification
                const handler = this.handlers.get(msg.method);
                if (handler) {
                    try {
                        const result = await handler(msg.params);
                        if (msg.id != null) {
                            this.send({ jsonrpc: '2.0', id: msg.id, result: result ?? null });
                        }
                    }
                    catch (err) {
                        if (msg.id != null) {
                            this.send({ jsonrpc: '2.0', id: msg.id, error: { code: -32000, message: err.message } });
                        }
                    }
                }
                else if (msg.id != null) {
                    this.send({ jsonrpc: '2.0', id: msg.id, error: { code: -32601, message: `Method not found: ${msg.method}` } });
                }
            }
            else if (msg.id != null) {
                // Response to our call
                const pending = this.pending.get(msg.id);
                if (pending) {
                    this.pending.delete(msg.id);
                    if (msg.error) {
                        pending.reject(new Error(msg.error.message));
                    }
                    else {
                        pending.resolve(msg.result);
                    }
                }
            }
        }
        catch (err) {
            process.stderr.write(`[extension-host] JSON parse error: ${err.message}\n`);
        }
    }
}
exports.JsonRpcConnection = JsonRpcConnection;
//# sourceMappingURL=rpc.js.map