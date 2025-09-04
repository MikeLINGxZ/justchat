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

}

