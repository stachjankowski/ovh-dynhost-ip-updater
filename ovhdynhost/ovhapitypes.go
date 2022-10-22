package ovhdynhost

// DynHostRecord is an object returned from OVH API endpoint /domain/zone/*/dynHost/record/*
type DynHostRecord struct {
	SubDomain string `json:"subDomain"`
	IP        string `json:"ip"`
	ID        int    `json:"id"`
	Zone      string `json:"zone"`
	TTL       int    `json:"ttl"`
}

// DynHostRecordPut is an object to update in OVH API endpoint /domain/zone/*/dynHost/record/*
type DynHostRecordPut struct {
	IP string `json:"ip"`
}
