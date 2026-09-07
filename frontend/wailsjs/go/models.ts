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
	    x_position: string;
	    y_position: string;
	    margin: number;
	
	    static createFrom(source: any = {}) {
	        return new RasterFeatureSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.x_position = source["x_position"];
	        this.y_position = source["y_position"];
	        this.margin = source["margin"];
	    }
	}
	export class RasterKeySettings {
	    type: string;
	    num_char: number;
	
	    static createFrom(source: any = {}) {
	        return new RasterKeySettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.num_char = source["num_char"];
	    }
	}
	export class GeoreferenceSettings {
	    master_map_source: string;
	    master_map: string;
	    attr_key: string;
	    raster_rotation: number;
	    raster_key_settings?: RasterKeySettings;
	    raster_feature_settings?: RasterFeatureSettings;
	
	    static createFrom(source: any = {}) {
	        return new GeoreferenceSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.master_map_source = source["master_map_source"];
	        this.master_map = source["master_map"];
	        this.attr_key = source["attr_key"];
	        this.raster_rotation = source["raster_rotation"];
	        this.raster_key_settings = this.convertValues(source["raster_key_settings"], RasterKeySettings);
	        this.raster_feature_settings = this.convertValues(source["raster_feature_settings"], RasterFeatureSettings);
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

