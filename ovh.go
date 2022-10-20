package main

import (
	"fmt"
	"net/netip"

	"github.com/ovh/go-ovh/ovh"
)

func GetClient() (*ovh.Client, error) {
	return ovh.NewDefaultClient()
}

func FindDynHostRecord(client *ovh.Client, zone string, subDomain string) (*DynHostRecord, error) {
	var ids []int
	if err := client.Get(fmt.Sprintf("/domain/zone/%s/dynHost/record", zone), &ids); err != nil {
		return nil, err
	}

	for _, dynHostId := range ids {
		var dynHostRecord DynHostRecord
		if err := client.Get(fmt.Sprintf("/domain/zone/%s/dynHost/record/%d", zone, dynHostId), &dynHostRecord); err != nil {
			return nil, err
		}
		if subDomain == dynHostRecord.SubDomain {
			return &dynHostRecord, nil
		}
	}

	return nil, fmt.Errorf("No DynHost for zone: %s sub-domain: %s in OVH", zone, subDomain)
}

func UpdateDynHostIP(client *ovh.Client, zone string, dynHostId int, newIp netip.Addr) error {
	params := &DynHostRecordPut{IP: newIp.String()}
	if err := client.Put(fmt.Sprintf("/domain/zone/%s/dynHost/record/%d", zone, dynHostId), params, nil); err != nil {
		return err
	}
	return nil
}
