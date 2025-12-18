export namespace core {
	
	export class LogEntry {
	    id: string;
	    // Go type: time
	    timestamp: any;
	    type: string;
	    domain: string;
	    src_ip: string;
	    dst_ip: string;
	    dst_port: number;
	    protocol: string;
	    process_name: string;
	    process_id: number;
	    bytes_sent: number;
	    bytes_recv: number;
	    status: string;
	    latency: number;
	
	    static createFrom(source: any = {}) {
	        return new LogEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.type = source["type"];
	        this.domain = source["domain"];
	        this.src_ip = source["src_ip"];
	        this.dst_ip = source["dst_ip"];
	        this.dst_port = source["dst_port"];
	        this.protocol = source["protocol"];
	        this.process_name = source["process_name"];
	        this.process_id = source["process_id"];
	        this.bytes_sent = source["bytes_sent"];
	        this.bytes_recv = source["bytes_recv"];
	        this.status = source["status"];
	        this.latency = source["latency"];
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
	export class Rule {
	    id: string;
	    type: string;
	    pattern: string;
	    enabled: boolean;
	    source: string;
	
	    static createFrom(source: any = {}) {
	        return new Rule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.pattern = source["pattern"];
	        this.enabled = source["enabled"];
	        this.source = source["source"];
	    }
	}
	export class PaginatedRulesResponse {
	    rules: Rule[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new PaginatedRulesResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rules = this.convertValues(source["rules"], Rule);
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
	
	export class Stats {
	    total_upload: number;
	    total_download: number;
	    active_connections: number;
	    top_domains: Record<string, number>;
	
	    static createFrom(source: any = {}) {
	        return new Stats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_upload = source["total_upload"];
	        this.total_download = source["total_download"];
	        this.active_connections = source["active_connections"];
	        this.top_domains = source["top_domains"];
	    }
	}
	export class TrafficDataPoint {
	    name: string;
	    upload: number;
	    download: number;
	
	    static createFrom(source: any = {}) {
	        return new TrafficDataPoint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.upload = source["upload"];
	        this.download = source["download"];
	    }
	}

}

export namespace main {
	
	export class AppInfo {
	    name: string;
	    version: string;
	    description: string;
	    author: string;
	    contact?: string;
	
	    static createFrom(source: any = {}) {
	        return new AppInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.version = source["version"];
	        this.description = source["description"];
	        this.author = source["author"];
	        this.contact = source["contact"];
	    }
	}

}

export namespace system {
	
	export class ConnectionInfo {
	    pid: number;
	    process_name: string;
	    local_addr: string;
	    remote_addr: string;
	    status: string;
	    protocol: string;
	
	    static createFrom(source: any = {}) {
	        return new ConnectionInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pid = source["pid"];
	        this.process_name = source["process_name"];
	        this.local_addr = source["local_addr"];
	        this.remote_addr = source["remote_addr"];
	        this.status = source["status"];
	        this.protocol = source["protocol"];
	    }
	}

}

