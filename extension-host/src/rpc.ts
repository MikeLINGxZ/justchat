import { Readable, Writable } from 'stream';
import * as readline from 'readline';

interface JsonRpcMessage {
  jsonrpc: string;
  id?: number;
  method?: string;
  params?: any;
  result?: any;
  error?: { code: number; message: string };
}

type RequestHandler = (params: any) => Promise<any>;

export class JsonRpcConnection {
  private nextId = 1;
  private pending = new Map<number, { resolve: (v: any) => void; reject: (e: Error) => void }>();
  private handlers = new Map<string, RequestHandler>();
  private rl: readline.Interface;
  private output: Writable;

  constructor(input: Readable, output: Writable) {
    this.output = output;
    this.rl = readline.createInterface({ input, crlfDelay: Infinity });
    this.rl.on('line', (line) => this.handleLine(line));
    this.rl.on('close', () => process.exit(0));
  }

  onRequest(method: string, handler: RequestHandler): void {
    this.handlers.set(method, handler);
  }

  async call(method: string, params?: any): Promise<any> {
    const id = this.nextId++;
    const msg: JsonRpcMessage = { jsonrpc: '2.0', id, method, params };
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

  notify(method: string, params?: any): void {
    this.send({ jsonrpc: '2.0', method, params });
  }

  private send(msg: JsonRpcMessage): void {
    this.output.write(JSON.stringify(msg) + '\n');
  }

  private async handleLine(line: string): Promise<void> {
    if (!line.trim()) return;
    try {
      const msg: JsonRpcMessage = JSON.parse(line);

      if (msg.method) {
        // Incoming request or notification
        const handler = this.handlers.get(msg.method);
        if (handler) {
          try {
            const result = await handler(msg.params);
            if (msg.id != null) {
              this.send({ jsonrpc: '2.0', id: msg.id, result: result ?? null });
            }
          } catch (err: any) {
            if (msg.id != null) {
              this.send({ jsonrpc: '2.0', id: msg.id, error: { code: -32000, message: err.message } });
            }
          }
        } else if (msg.id != null) {
          this.send({ jsonrpc: '2.0', id: msg.id, error: { code: -32601, message: `Method not found: ${msg.method}` } });
        }
      } else if (msg.id != null) {
        // Response to our call
        const pending = this.pending.get(msg.id);
        if (pending) {
          this.pending.delete(msg.id);
          if (msg.error) {
            pending.reject(new Error(msg.error.message));
          } else {
            pending.resolve(msg.result);
          }
        }
      }
    } catch (err: any) {
      process.stderr.write(`[extension-host] JSON parse error: ${err.message}\n`);
    }
  }
}
