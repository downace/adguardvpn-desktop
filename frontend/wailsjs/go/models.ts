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
	export interface Status {
	    connected: boolean;
	}

}

