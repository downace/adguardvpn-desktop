export namespace adguard {
	
	export enum ExclusionMode {
	    GENERAL = "general",
	    SELECTIVE = "selective",
	}
	export interface Subscription {
	    type: string;
	    validUntil: string;
	    maxDevices: number;
	}
	export interface Account {
	    username: string;
	    subscription: Subscription;
	}
	export interface Location {
	    iso: string;
	    country: string;
	    city: string;
	    ping: number;
	}
	export interface Status {
	    connecting: boolean;
	    connected: boolean;
	    location?: Location;
	    mode: string;
	}

}

