export namespace admin {
	
	export class Ban {
	    netId: string;
	    text: string;
	
	    static createFrom(source: any = {}) {
	        return new Ban(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.netId = source["netId"];
	        this.text = source["text"];
	    }
	}
	export class Player {
	    id: string;
	    name: string;
	    netId: string;
	    ip: string;
	    score: number;
	    bot: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Player(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.netId = source["netId"];
	        this.ip = source["ip"];
	        this.score = source["score"];
	        this.bot = source["bot"];
	    }
	}

}

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
	export class ProfileRef {
	    id: string;
	    name: string;
	    hostname: string;
	    active: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ProfileRef(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.hostname = source["hostname"];
	        this.active = source["active"];
	    }
	}
	export class Settings {
	    language: string;
	    showConsole: boolean;
	    termsAccepted?: string;
	    termsAcceptedAt?: string;
	    logCleanup: boolean;
	    logKeep?: number;
	    startWithWindows: boolean;
	    autoStartServer: boolean;
	
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
	        this.startWithWindows = source["startWithWindows"];
	        this.autoStartServer = source["autoStartServer"];
	    }
	}

}

export namespace desktop {
	
	export class AutomationInfo {
	    // Go type: time
	    nextRestart: any;
	    hasWebhook: boolean;
	    hasSteamCmd: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AutomationInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nextRestart = this.convertValues(source["nextRestart"], null);
	        this.hasWebhook = source["hasWebhook"];
	        this.hasSteamCmd = source["hasSteamCmd"];
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
	export class CrashCheck {
	    mods: number;
	    sizeMismatch: boolean;
	    findings: modcheck.Finding[];
	
	    static createFrom(source: any = {}) {
	        return new CrashCheck(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mods = source["mods"];
	        this.sizeMismatch = source["sizeMismatch"];
	        this.findings = this.convertValues(source["findings"], modcheck.Finding);
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
	export class CycleDefault {
	    lines: string[];
	    serverHasFile: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CycleDefault(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.lines = source["lines"];
	        this.serverHasFile = source["serverHasFile"];
	    }
	}
	export class CycleLine {
	    scenario: string;
	    lighting: string;
	    missing: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CycleLine(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.scenario = source["scenario"];
	        this.lighting = source["lighting"];
	        this.missing = source["missing"];
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
	export class FirewallInfo {
	    ready: boolean;
	    rules: boolean;
	
	    static createFrom(source: any = {}) {
	        return new FirewallInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ready = source["ready"];
	        this.rules = source["rules"];
	    }
	}
	export class ImageInfo {
	    name: string;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new ImageInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.size = source["size"];
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
	export class ModCheck {
	    mods: number;
	    conflicts: modcheck.Conflict[];
	    unreadable: modcheck.ModRef[];
	
	    static createFrom(source: any = {}) {
	        return new ModCheck(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mods = source["mods"];
	        this.conflicts = this.convertValues(source["conflicts"], modcheck.Conflict);
	        this.unreadable = this.convertValues(source["unreadable"], modcheck.ModRef);
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
	export class ReportRequest {
	    kind: string;
	    title: string;
	    description: string;
	    includeLog: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ReportRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.kind = source["kind"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.includeLog = source["includeLog"];
	    }
	}
	export class ReportResult {
	    file: string;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new ReportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.file = source["file"];
	        this.size = source["size"];
	    }
	}
	export class RulesInfo {
	    groups: string[];
	    rules: server.RuleDef[];
	
	    static createFrom(source: any = {}) {
	        return new RulesInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.groups = source["groups"];
	        this.rules = this.convertValues(source["rules"], server.RuleDef);
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
	    foldersNested: boolean;
	
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
	        this.foldersNested = source["foldersNested"];
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
	    last: string;
	    validate: boolean;
	    upToDate: boolean;
	    present: boolean;
	
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
	        this.last = source["last"];
	        this.validate = source["validate"];
	        this.upToDate = source["upToDate"];
	        this.present = source["present"];
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

export namespace mapchange {
	
	export class Info {
	    rcon: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Info(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rcon = source["rcon"];
	    }
	}

}

export namespace modcheck {
	
	export class ModRef {
	    id: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new ModRef(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}
	export class Conflict {
	    a: ModRef;
	    b: ModRef;
	    count: number;
	    example: string;
	
	    static createFrom(source: any = {}) {
	        return new Conflict(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.a = this.convertValues(source["a"], ModRef);
	        this.b = this.convertValues(source["b"], ModRef);
	        this.count = source["count"];
	        this.example = source["example"];
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
	export class Finding {
	    object: string;
	    mods: ModRef[];
	
	    static createFrom(source: any = {}) {
	        return new Finding(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.object = source["object"];
	        this.mods = this.convertValues(source["mods"], ModRef);
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
	
	export class ModMutator {
	    id: string;
	    modId: number;
	    modName: string;
	
	    static createFrom(source: any = {}) {
	        return new ModMutator(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.modId = source["modId"];
	        this.modName = source["modName"];
	    }
	}
	export class ModScenario {
	    map: string;
	    scenario: string;
	    modId: number;
	    modName: string;
	
	    static createFrom(source: any = {}) {
	        return new ModScenario(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.map = source["map"];
	        this.scenario = source["scenario"];
	        this.modId = source["modId"];
	        this.modName = source["modName"];
	    }
	}
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

export namespace netcheck {
	
	export class Item {
	    key: string;
	    status: string;
	    params?: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new Item(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.status = source["status"];
	        this.params = source["params"];
	    }
	}
	export class Result {
	    publicIp: string;
	    gamePort: number;
	    queryPort: number;
	    verdict: string;
	    items: Item[];
	
	    static createFrom(source: any = {}) {
	        return new Result(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.publicIp = source["publicIp"];
	        this.gamePort = source["gamePort"];
	        this.queryPort = source["queryPort"];
	        this.verdict = source["verdict"];
	        this.items = this.convertValues(source["items"], Item);
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

export namespace players {
	
	export class View {
	    players: admin.Player[];
	    bans: admin.Ban[];
	
	    static createFrom(source: any = {}) {
	        return new View(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.players = this.convertValues(source["players"], admin.Player);
	        this.bans = this.convertValues(source["bans"], admin.Ban);
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
	    refusedMutators: string[];
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
	        this.refusedMutators = source["refusedMutators"];
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
	export class AutoRestart {
	    daily: string;
	    warnMessage: string;
	    update: boolean;
	    onCrash: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AutoRestart(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.daily = source["daily"];
	        this.warnMessage = source["warnMessage"];
	        this.update = source["update"];
	        this.onCrash = source["onCrash"];
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
	export class Discord {
	    ready: string;
	    stopped: string;
	    crashed: string;
	    join: string;
	    leave: string;
	
	    static createFrom(source: any = {}) {
	        return new Discord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ready = source["ready"];
	        this.stopped = source["stopped"];
	        this.crashed = source["crashed"];
	        this.join = source["join"];
	        this.leave = source["leave"];
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
	export class Mutator {
	    id: string;
	    label: string;
	
	    static createFrom(source: any = {}) {
	        return new Mutator(source);
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
	    startLighting: string;
	    travelMap: string;
	    travelScenario: string;
	    useMods: boolean;
	    mutators: string[];
	    rcon: boolean;
	    rconPort: number;
	    bindIp: string;
	    connectIp: string;
	    admins: Admin[];
	    mapCycle: string[];
	    modIds: number[];
	    cmdModList: boolean;
	    motd: string;
	    joinMessage: string;
	    leaveMessage: string;
	    announcements: string[];
	    announceEvery: number;
	    rules: Record<string, string>;
	    autoRestart: AutoRestart;
	    discord: Discord;
	
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
	        this.startLighting = source["startLighting"];
	        this.travelMap = source["travelMap"];
	        this.travelScenario = source["travelScenario"];
	        this.useMods = source["useMods"];
	        this.mutators = source["mutators"];
	        this.rcon = source["rcon"];
	        this.rconPort = source["rconPort"];
	        this.bindIp = source["bindIp"];
	        this.connectIp = source["connectIp"];
	        this.admins = this.convertValues(source["admins"], Admin);
	        this.mapCycle = source["mapCycle"];
	        this.modIds = source["modIds"];
	        this.cmdModList = source["cmdModList"];
	        this.motd = source["motd"];
	        this.joinMessage = source["joinMessage"];
	        this.leaveMessage = source["leaveMessage"];
	        this.announcements = source["announcements"];
	        this.announceEvery = source["announceEvery"];
	        this.rules = source["rules"];
	        this.autoRestart = this.convertValues(source["autoRestart"], AutoRestart);
	        this.discord = this.convertValues(source["discord"], Discord);
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
	export class RuleDef {
	    id: string;
	    group: string;
	    file: string;
	    section: string;
	    key: string;
	    kind: string;
	    min: number;
	    max: number;
	    default: string;
	
	    static createFrom(source: any = {}) {
	        return new RuleDef(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.group = source["group"];
	        this.file = source["file"];
	        this.section = source["section"];
	        this.key = source["key"];
	        this.kind = source["kind"];
	        this.min = source["min"];
	        this.max = source["max"];
	        this.default = source["default"];
	    }
	}
	export class Scenario {
	    id: string;
	    map: string;
	    mode: string;
	    source: string;
	
	    static createFrom(source: any = {}) {
	        return new Scenario(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.map = source["map"];
	        this.mode = source["mode"];
	        this.source = source["source"];
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

export namespace winsys {
	
	export class Address {
	    ip: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new Address(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ip = source["ip"];
	        this.name = source["name"];
	    }
	}

}

