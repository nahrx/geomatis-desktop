export namespace storage {
	
	export class Config {
	    DB_HOST: string;
	    DB_PORT: number;
	    DB_DATABASE: string;
	    DB_USERNAME: string;
	    DB_PASSWORD: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.DB_HOST = source["DB_HOST"];
	        this.DB_PORT = source["DB_PORT"];
	        this.DB_DATABASE = source["DB_DATABASE"];
	        this.DB_USERNAME = source["DB_USERNAME"];
	        this.DB_PASSWORD = source["DB_PASSWORD"];
	    }
	}

}

export namespace types {
	
	export class Extent {
	    minX: number;
	    minY: number;
	    maxX: number;
	    maxY: number;
	
	    static createFrom(source: any = {}) {
	        return new Extent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.minX = source["minX"];
	        this.minY = source["minY"];
	        this.maxX = source["maxX"];
	        this.maxY = source["maxY"];
	    }
	}
	export class RasterFeatureSettings {
	    XPosition: string;
	    YPosition: string;
	    Margin: number;
	
	    static createFrom(source: any = {}) {
	        return new RasterFeatureSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.XPosition = source["XPosition"];
	        this.YPosition = source["YPosition"];
	        this.Margin = source["Margin"];
	    }
	}
	export class RasterKeySettings {
	    Type: string;
	    NumChar: number;
	    // Go type: regexp
	    Regex?: any;
	
	    static createFrom(source: any = {}) {
	        return new RasterKeySettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Type = source["Type"];
	        this.NumChar = source["NumChar"];
	        this.Regex = this.convertValues(source["Regex"], null);
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
	export class GeoreferenceSettings {
	    MasterMapSource: string;
	    MasterMap: string;
	    AttrKey: string;
	    RasterKeySettings?: RasterKeySettings;
	    RasterFeatureSettings?: RasterFeatureSettings;
	
	    static createFrom(source: any = {}) {
	        return new GeoreferenceSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.MasterMapSource = source["MasterMapSource"];
	        this.MasterMap = source["MasterMap"];
	        this.AttrKey = source["AttrKey"];
	        this.RasterKeySettings = this.convertValues(source["RasterKeySettings"], RasterKeySettings);
	        this.RasterFeatureSettings = this.convertValues(source["RasterFeatureSettings"], RasterFeatureSettings);
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
	export class MasterMap {
	    name: string;
	    dimension: number;
	    srid: number;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new MasterMap(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.dimension = source["dimension"];
	        this.srid = source["srid"];
	        this.type = source["type"];
	    }
	}
	

}

