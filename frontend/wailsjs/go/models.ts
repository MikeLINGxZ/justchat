export namespace data_models {
	
	export class Model {
	    id: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	    // Go type: gorm
	    deleted_at: any;
	    provider_id: number;
	    model: string;
	    owned_by: string;
	    object: string;
	    enable: boolean;
	    alias?: string;
	
	    static createFrom(source: any = {}) {
	        return new Model(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	        this.deleted_at = this.convertValues(source["deleted_at"], null);
	        this.provider_id = source["provider_id"];
	        this.model = source["model"];
	        this.owned_by = source["owned_by"];
	        this.object = source["object"];
	        this.enable = source["enable"];
	        this.alias = source["alias"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Provider {
	    id: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	    // Go type: gorm
	    deleted_at: any;
	    provider_name: string;
	    base_url: string;
	    api_key: string;
	    enable: boolean;
	    alias?: string;
	
	    static createFrom(source: any = {}) {
	        return new Provider(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	        this.deleted_at = this.convertValues(source["deleted_at"], null);
	        this.provider_name = source["provider_name"];
	        this.base_url = source["base_url"];
	        this.api_key = source["api_key"];
	        this.enable = source["enable"];
	        this.alias = source["alias"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace schema {
	
	export class ChatMessageAudioURL {
	    url?: string;
	    uri?: string;
	    mime_type?: string;
	    extra?: {[key: string]: any};
	
	    static createFrom(source: any = {}) {
	        return new ChatMessageAudioURL(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.uri = source["uri"];
	        this.mime_type = source["mime_type"];
	        this.extra = source["extra"];
	    }
	}
	export class ChatMessageFileURL {
	    url?: string;
	    uri?: string;
	    mime_type?: string;
	    name?: string;
	    extra?: {[key: string]: any};
	
	    static createFrom(source: any = {}) {
	        return new ChatMessageFileURL(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.uri = source["uri"];
	        this.mime_type = source["mime_type"];
	        this.name = source["name"];
	        this.extra = source["extra"];
	    }
	}
	export class ChatMessageImageURL {
	    url?: string;
	    uri?: string;
	    detail?: string;
	    mime_type?: string;
	    extra?: {[key: string]: any};
	
	    static createFrom(source: any = {}) {
	        return new ChatMessageImageURL(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.uri = source["uri"];
	        this.detail = source["detail"];
	        this.mime_type = source["mime_type"];
	        this.extra = source["extra"];
	    }
	}
	export class ChatMessageVideoURL {
	    url?: string;
	    uri?: string;
	    mime_type?: string;
	    extra?: {[key: string]: any};
	
	    static createFrom(source: any = {}) {
	        return new ChatMessageVideoURL(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.uri = source["uri"];
	        this.mime_type = source["mime_type"];
	        this.extra = source["extra"];
	    }
	}
	export class ChatMessagePart {
	    type?: string;
	    text?: string;
	    image_url?: ChatMessageImageURL;
	    audio_url?: ChatMessageAudioURL;
	    video_url?: ChatMessageVideoURL;
	    file_url?: ChatMessageFileURL;
	
	    static createFrom(source: any = {}) {
	        return new ChatMessagePart(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.text = source["text"];
	        this.image_url = this.convertValues(source["image_url"], ChatMessageImageURL);
	        this.audio_url = this.convertValues(source["audio_url"], ChatMessageAudioURL);
	        this.video_url = this.convertValues(source["video_url"], ChatMessageVideoURL);
	        this.file_url = this.convertValues(source["file_url"], ChatMessageFileURL);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class FunctionCall {
	    name?: string;
	    arguments?: string;
	
	    static createFrom(source: any = {}) {
	        return new FunctionCall(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.arguments = source["arguments"];
	    }
	}
	export class TopLogProb {
	    token: string;
	    logprob: number;
	    bytes?: number[];
	
	    static createFrom(source: any = {}) {
	        return new TopLogProb(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.token = source["token"];
	        this.logprob = source["logprob"];
	        this.bytes = source["bytes"];
	    }
	}
	export class LogProb {
	    token: string;
	    logprob: number;
	    bytes?: number[];
	    top_logprobs: TopLogProb[];
	
	    static createFrom(source: any = {}) {
	        return new LogProb(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.token = source["token"];
	        this.logprob = source["logprob"];
	        this.bytes = source["bytes"];
	        this.top_logprobs = this.convertValues(source["top_logprobs"], TopLogProb);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LogProbs {
	    content: LogProb[];
	
	    static createFrom(source: any = {}) {
	        return new LogProbs(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.content = this.convertValues(source["content"], LogProb);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PromptTokenDetails {
	    cached_tokens: number;
	
	    static createFrom(source: any = {}) {
	        return new PromptTokenDetails(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cached_tokens = source["cached_tokens"];
	    }
	}
	export class TokenUsage {
	    prompt_tokens: number;
	    prompt_token_details: PromptTokenDetails;
	    completion_tokens: number;
	    total_tokens: number;
	
	    static createFrom(source: any = {}) {
	        return new TokenUsage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.prompt_tokens = source["prompt_tokens"];
	        this.prompt_token_details = this.convertValues(source["prompt_token_details"], PromptTokenDetails);
	        this.completion_tokens = source["completion_tokens"];
	        this.total_tokens = source["total_tokens"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ResponseMeta {
	    finish_reason?: string;
	    usage?: TokenUsage;
	    logprobs?: LogProbs;
	
	    static createFrom(source: any = {}) {
	        return new ResponseMeta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.finish_reason = source["finish_reason"];
	        this.usage = this.convertValues(source["usage"], TokenUsage);
	        this.logprobs = this.convertValues(source["logprobs"], LogProbs);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ToolCall {
	    index?: number;
	    id: string;
	    type: string;
	    function: FunctionCall;
	    extra?: {[key: string]: any};
	
	    static createFrom(source: any = {}) {
	        return new ToolCall(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.index = source["index"];
	        this.id = source["id"];
	        this.type = source["type"];
	        this.function = this.convertValues(source["function"], FunctionCall);
	        this.extra = source["extra"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Message {
	    role: string;
	    content: string;
	    multi_content?: ChatMessagePart[];
	    name?: string;
	    tool_calls?: ToolCall[];
	    tool_call_id?: string;
	    tool_name?: string;
	    response_meta?: ResponseMeta;
	    reasoning_content?: string;
	    extra?: {[key: string]: any};
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.role = source["role"];
	        this.content = source["content"];
	        this.multi_content = this.convertValues(source["multi_content"], ChatMessagePart);
	        this.name = source["name"];
	        this.tool_calls = this.convertValues(source["tool_calls"], ToolCall);
	        this.tool_call_id = source["tool_call_id"];
	        this.tool_name = source["tool_name"];
	        this.response_meta = this.convertValues(source["response_meta"], ResponseMeta);
	        this.reasoning_content = source["reasoning_content"];
	        this.extra = source["extra"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	

}

export namespace view_models {
	
	export class MatchMessage {
	    role: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new MatchMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.role = source["role"];
	        this.content = source["content"];
	    }
	}
	export class Chat {
	    id: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	    // Go type: gorm
	    deleted_at: any;
	    uuid: string;
	    model_id: number;
	    title: string;
	    prompt: string;
	    content: MatchMessage[];
	    reasoning_content: MatchMessage[];
	
	    static createFrom(source: any = {}) {
	        return new Chat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	        this.deleted_at = this.convertValues(source["deleted_at"], null);
	        this.uuid = source["uuid"];
	        this.model_id = source["model_id"];
	        this.title = source["title"];
	        this.prompt = source["prompt"];
	        this.content = this.convertValues(source["content"], MatchMessage);
	        this.reasoning_content = this.convertValues(source["reasoning_content"], MatchMessage);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ChatList {
	    lists: Chat[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new ChatList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.lists = this.convertValues(source["lists"], Chat);
	        this.total = source["total"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class MessageList {
	    messages: schema.Message[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new MessageList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.messages = this.convertValues(source["messages"], schema.Message);
	        this.total = source["total"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

