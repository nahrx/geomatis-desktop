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

