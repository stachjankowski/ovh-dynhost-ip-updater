package main

type DynHostRecord struct {
	SubDomain string `json:"subDomain"`
	IP        string `json:"ip"`
	ID        int    `json:"id"`
	Zone      string `json:"zone"`
	TTL       int    `json:"ttl"`
}

type DynHostRecordPut struct {
	IP string `json:"ip"`
}
