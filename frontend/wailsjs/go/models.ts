export namespace adguard {
	
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
	    connected: boolean;
	    location?: Location;
	    mode: string;
	}

}

