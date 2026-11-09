export namespace backupsvc {
	
	export class Summary {
	    // Go type: time
	    createdAt: any;
	    appVersion: string;
	    machine: string;
	    hostname: string;
	    oldDir: string;
	    files: string[];
	    secrets: string[];
	    hasModio: boolean;
	    admins: number;
	    mods: number;
	
	    static createFrom(source: any = {}) {
	        return new Summary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.appVersion = source["appVersion"];
	        this.machine = source["machine"];
	        this.hostname = source["hostname"];
	        this.oldDir = source["oldDir"];
	        this.files = source["files"];
	        this.secrets = source["secrets"];
	        this.hasModio = source["hasModio"];
	        this.admins = source["admins"];
	        this.mods = source["mods"];
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

export namespace config {
	
	export class ApplyResult {
	    profile: server.Profile;
	    applied: string[];
	    kept: string[];
	    skipped: string[];
	    rconInIni: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ApplyResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.profile = this.convertValues(source["profile"], server.Profile);
	        this.applied = source["applied"];
	        this.kept = source["kept"];
	        this.skipped = source["skipped"];
	        this.rconInIni = source["rconInIni"];
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
	export class Settings {
	    language: string;
	    showConsole: boolean;
	    termsAccepted?: string;
	    termsAcceptedAt?: string;
	    logCleanup: boolean;
	    logKeep?: number;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.language = source["language"];
	        this.showConsole = source["showConsole"];
	        this.termsAccepted = source["termsAccepted"];
	        this.termsAcceptedAt = source["termsAcceptedAt"];
	        this.logCleanup = source["logCleanup"];
	        this.logKeep = source["logKeep"];
	    }
	}

}

export namespace desktop {
	
	export class Connect {
	    ip: string;
	    port: number;
	    command: string;
	
	    static createFrom(source: any = {}) {
	        return new Connect(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ip = source["ip"];
	        this.port = source["port"];
	        this.command = source["command"];
	    }
	}
	export class CreateResult {
	    saved: boolean;
	    file: string;
	    summary: backupsvc.Summary;
	
	    static createFrom(source: any = {}) {
	        return new CreateResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.saved = source["saved"];
	        this.file = source["file"];
	        this.summary = this.convertValues(source["summary"], backupsvc.Summary);
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
	export class ExistingScan {
	    found: boolean;
	    logFile: string;
	    diffs: server.Diff[];
	    rconInIni: boolean;
	    running: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ExistingScan(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.found = source["found"];
	        this.logFile = source["logFile"];
	        this.diffs = this.convertValues(source["diffs"], server.Diff);
	        this.rconInIni = source["rconInIni"];
	        this.running = source["running"];
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
	export class Info {
	    appVersion: string;
	    goVersion: string;
	    languages: string[];
	    problems: string[];
	
	    static createFrom(source: any = {}) {
	        return new Info(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.appVersion = source["appVersion"];
	        this.goVersion = source["goVersion"];
	        this.languages = source["languages"];
	        this.problems = source["problems"];
	    }
	}
	export class LogCleanup {
	    enabled: boolean;
	    keep: number;
	    min: number;
	    max: number;
	    runs: number;
	    old: number;
	    oldSize: number;
	    deleted: number;
	    freed: number;
	
	    static createFrom(source: any = {}) {
	        return new LogCleanup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.keep = source["keep"];
	        this.min = source["min"];
	        this.max = source["max"];
	        this.runs = source["runs"];
	        this.old = source["old"];
	        this.oldSize = source["oldSize"];
	        this.deleted = source["deleted"];
	        this.freed = source["freed"];
	    }
	}
	export class LogRun {
	    id: string;
	    // Go type: time
	    started: any;
	    // Go type: time
	    ended: any;
	    size: number;
	    current: boolean;
	
	    static createFrom(source: any = {}) {
	        return new LogRun(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.started = this.convertValues(source["started"], null);
	        this.ended = this.convertValues(source["ended"], null);
	        this.size = source["size"];
	        this.current = source["current"];
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
	export class MapsInfo {
	    maps: server.Map[];
	    modes: server.Mode[];
	
	    static createFrom(source: any = {}) {
	        return new MapsInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.maps = this.convertValues(source["maps"], server.Map);
	        this.modes = this.convertValues(source["modes"], server.Mode);
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
	export class ModRow {
	    id: number;
	    name: string;
	    size: number;
	    onDisk: number;
	    installed: boolean;
	    subscribed: boolean;
	    live: string;
	    done: number;
	
	    static createFrom(source: any = {}) {
	        return new ModRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.size = source["size"];
	        this.onDisk = source["onDisk"];
	        this.installed = source["installed"];
	        this.subscribed = source["subscribed"];
	        this.live = source["live"];
	        this.done = source["done"];
	    }
	}
	export class ModsView {
	    session: modio.Session;
	    hasApiKey: boolean;
	    hasEmail: boolean;
	    root: string;
	    mods: ModRow[];
	    // Go type: time
	    subscribedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new ModsView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.session = this.convertValues(source["session"], modio.Session);
	        this.hasApiKey = source["hasApiKey"];
	        this.hasEmail = source["hasEmail"];
	        this.root = source["root"];
	        this.mods = this.convertValues(source["mods"], ModRow);
	        this.subscribedAt = this.convertValues(source["subscribedAt"], null);
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
	export class Step {
	    id: string;
	    done: boolean;
	    optional: boolean;
	    page: string;
	
	    static createFrom(source: any = {}) {
	        return new Step(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.done = source["done"];
	        this.optional = source["optional"];
	        this.page = source["page"];
	    }
	}
	export class Terms {
	    version: string;
	    accepted: boolean;
	    acceptedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Terms(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.accepted = source["accepted"];
	        this.acceptedAt = source["acceptedAt"];
	    }
	}

}

export namespace installer {
	
	export class Info {
	    steamcmdFound: boolean;
	    serverInstalled: boolean;
	    buildId: string;
	    sizeOnDisk: number;
	    lastUpdated: number;
	
	    static createFrom(source: any = {}) {
	        return new Info(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.steamcmdFound = source["steamcmdFound"];
	        this.serverInstalled = source["serverInstalled"];
	        this.buildId = source["buildId"];
	        this.sizeOnDisk = source["sizeOnDisk"];
	        this.lastUpdated = source["lastUpdated"];
	    }
	}
	export class State {
	    task: string;
	    phase: string;
	    percent: number;
	    done: number;
	    total: number;
	    speed: number;
	    eta: number;
	    // Go type: time
	    started: any;
	    error?: problems.Error;
	    finished: string;
	
	    static createFrom(source: any = {}) {
	        return new State(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.task = source["task"];
	        this.phase = source["phase"];
	        this.percent = source["percent"];
	        this.done = source["done"];
	        this.total = source["total"];
	        this.speed = source["speed"];
	        this.eta = source["eta"];
	        this.started = this.convertValues(source["started"], null);
	        this.error = this.convertValues(source["error"], problems.Error);
	        this.finished = source["finished"];
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

export namespace logbuf {
	
	export class Line {
	    seq: number;
	    // Go type: time
	    time: any;
	    source: string;
	    level: string;
	    noise: boolean;
	    event: boolean;
	    text: string;
	    run: number;
	    kind?: string;
	
	    static createFrom(source: any = {}) {
	        return new Line(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.seq = source["seq"];
	        this.time = this.convertValues(source["time"], null);
	        this.source = source["source"];
	        this.level = source["level"];
	        this.noise = source["noise"];
	        this.event = source["event"];
	        this.text = source["text"];
	        this.run = source["run"];
	        this.kind = source["kind"];
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

export namespace modio {
	
	export class Session {
	    found: boolean;
	    username: string;
	    // Go type: time
	    expires: any;
	    valid: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Session(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.found = source["found"];
	        this.username = source["username"];
	        this.expires = this.convertValues(source["expires"], null);
	        this.valid = source["valid"];
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

export namespace problems {
	
	export class Error {
	    code: string;
	    field?: string;
	    msg?: string;
	    params?: Record<string, string>;
	    detail?: string;
	    links?: string[];
	
	    static createFrom(source: any = {}) {
	        return new Error(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.field = source["field"];
	        this.msg = source["msg"];
	        this.params = source["params"];
	        this.detail = source["detail"];
	        this.links = source["links"];
	    }
	}

}

export namespace runner {
	
	export class ModState {
	    id: number;
	    name: string;
	    state: string;
	    done: number;
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new ModState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.state = source["state"];
	        this.done = source["done"];
	        this.total = source["total"];
	    }
	}
	export class Status {
	    phase: string;
	    pids: number[];
	    // Go type: time
	    since: any;
	    map: string;
	    scenario: string;
	    matchState: string;
	    port: number;
	    mods: ModState[];
	    problem?: problems.Error;
	    cmdLine: string;
	    external: boolean;
	    authorizing: boolean;
	    gameStatsMissing: boolean;
	    remount: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Status(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.phase = source["phase"];
	        this.pids = source["pids"];
	        this.since = this.convertValues(source["since"], null);
	        this.map = source["map"];
	        this.scenario = source["scenario"];
	        this.matchState = source["matchState"];
	        this.port = source["port"];
	        this.mods = this.convertValues(source["mods"], ModState);
	        this.problem = this.convertValues(source["problem"], problems.Error);
	        this.cmdLine = source["cmdLine"];
	        this.external = source["external"];
	        this.authorizing = source["authorizing"];
	        this.gameStatsMissing = source["gameStatsMissing"];
	        this.remount = source["remount"];
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

export namespace server {
	
	export class Admin {
	    steamId: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new Admin(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.steamId = source["steamId"];
	        this.name = source["name"];
	    }
	}
	export class Diff {
	    key: string;
	    kind: string;
	    mine: string;
	    server: string;
	    state: string;
	
	    static createFrom(source: any = {}) {
	        return new Diff(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.kind = source["kind"];
	        this.mine = source["mine"];
	        this.server = source["server"];
	        this.state = source["state"];
	    }
	}
	export class Map {
	    label: string;
	    map: string;
	    scenario: string;
	
	    static createFrom(source: any = {}) {
	        return new Map(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.label = source["label"];
	        this.map = source["map"];
	        this.scenario = source["scenario"];
	    }
	}
	export class Mode {
	    id: string;
	    label: string;
	
	    static createFrom(source: any = {}) {
	        return new Mode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.label = source["label"];
	    }
	}
	export class Profile {
	    serverDir: string;
	    steamcmdDir: string;
	    hostname: string;
	    port: number;
	    queryPort: number;
	    maxPlayers: number;
	    startMap: string;
	    startScenario: string;
	    travelMap: string;
	    travelScenario: string;
	    useMods: boolean;
	    mutators: string[];
	    rcon: boolean;
	    rconPort: number;
	    connectIp: string;
	    admins: Admin[];
	    mapCycle: string[];
	    modIds: number[];
	
	    static createFrom(source: any = {}) {
	        return new Profile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.serverDir = source["serverDir"];
	        this.steamcmdDir = source["steamcmdDir"];
	        this.hostname = source["hostname"];
	        this.port = source["port"];
	        this.queryPort = source["queryPort"];
	        this.maxPlayers = source["maxPlayers"];
	        this.startMap = source["startMap"];
	        this.startScenario = source["startScenario"];
	        this.travelMap = source["travelMap"];
	        this.travelScenario = source["travelScenario"];
	        this.useMods = source["useMods"];
	        this.mutators = source["mutators"];
	        this.rcon = source["rcon"];
	        this.rconPort = source["rconPort"];
	        this.connectIp = source["connectIp"];
	        this.admins = this.convertValues(source["admins"], Admin);
	        this.mapCycle = source["mapCycle"];
	        this.modIds = source["modIds"];
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

export namespace updater {
	
	export class State {
	    status: string;
	    current: string;
	    version: string;
	    progress: number;
	    stage: string;
	
	    static createFrom(source: any = {}) {
	        return new State(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.current = source["current"];
	        this.version = source["version"];
	        this.progress = source["progress"];
	        this.stage = source["stage"];
	    }
	}

}

