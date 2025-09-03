export namespace data_models {
	
	export class Model {
	    id: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	    // Go type: gorm
	    DeletedAt: any;
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
	        this.DeletedAt = this.convertValues(source["DeletedAt"], null);
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

}

