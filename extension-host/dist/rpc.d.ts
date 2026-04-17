import { Readable, Writable } from 'stream';
type RequestHandler = (params: any) => Promise<any>;
export declare class JsonRpcConnection {
    private nextId;
    private pending;
    private handlers;
    private rl;
    private output;
    constructor(input: Readable, output: Writable);
    onRequest(method: string, handler: RequestHandler): void;
    call(method: string, params?: any): Promise<any>;
    notify(method: string, params?: any): void;
    private send;
    private handleLine;
}
export {};
