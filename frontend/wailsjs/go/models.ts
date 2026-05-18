export namespace backend {
	
	export class AppConfig {
	    apiKey: string;
	    models: string[];
	    systemPrompt: string;
	    outputDir: string;
	
	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.apiKey = source["apiKey"];
	        this.models = source["models"];
	        this.systemPrompt = source["systemPrompt"];
	        this.outputDir = source["outputDir"];
	    }
	}

}

