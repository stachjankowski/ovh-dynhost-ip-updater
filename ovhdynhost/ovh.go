package ovhdynhost

import (
	"fmt"
	"net/netip"

	"github.com/ovh/go-ovh/ovh"
)

// GetClient returns an ovh.Client
func GetClient() (*ovh.Client, error) {
	return ovh.NewDefaultClient()
}

// FindDynHostRecord searches for a dynhost record in the OVH api
func FindDynHostRecord(client *ovh.Client, zone string, subDomain string) (*DynHostRecord, error) {
	var ids []int
	if err := client.Get(fmt.Sprintf("/domain/zone/%s/dynHost/record", zone), &ids); err != nil {
		return nil, err
	}

	for _, dynHostID := range ids {
		var dynHostRecord DynHostRecord
		if err := client.Get(fmt.Sprintf("/domain/zone/%s/dynHost/record/%d", zone, dynHostID), &dynHostRecord); err != nil {
			return nil, err
		}
		if subDomain == dynHostRecord.SubDomain {
			return &dynHostRecord, nil
		}
	}

	return nil, fmt.Errorf("No DynHost for zone: %s sub-domain: %s in OVH", zone, subDomain)
}

// UpdateDynHostIP updates the IP address in OVH
func UpdateDynHostIP(client *ovh.Client, zone string, dynHostID int, newIP netip.Addr) error {
	params := &DynHostRecordPut{IP: newIP.String()}
	if err := client.Put(fmt.Sprintf("/domain/zone/%s/dynHost/record/%d", zone, dynHostID), params, nil); err != nil {
		return err
	}
	return nil
}
