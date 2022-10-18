package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/netip"
	"os"
	"time"

	arg "github.com/alexflint/go-arg"
	"github.com/ovh/go-ovh/ovh"
	"github.com/sirupsen/logrus"
	"github.com/spyzhov/ajson"
)

var log = logrus.New()

func findDynHostRecord(client *ovh.Client, zone string, subDomain string) (*DynHostRecord, error) {
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

func updateDynHostIP(client *ovh.Client, zone string, dynHostId int, newIp netip.Addr) error {
	params := &DynHostRecordPut{IP: newIp.String()}
	if err := client.Put(fmt.Sprintf("/domain/zone/%s/dynHost/record/%d", zone, dynHostId), params, nil); err != nil {
		return err
	}
	return nil
}

func GetIp(url string, jsonPath string) (netip.Addr, error) {
	req, err := http.Get(url)
	if err != nil {
		return netip.Addr{}, err
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return netip.Addr{}, err
	}

	// raw format
	if jsonPath == "" {
		return netip.ParseAddr(string(body))
	}

	// json
	root, err := ajson.Unmarshal(body)
	if err != nil {
		return netip.Addr{}, err
	}
	nodes, err := root.JSONPath(jsonPath)
	for _, node := range nodes {
		value, err := node.GetString()
		if err != nil {
			return netip.Addr{}, err
		}
		return netip.ParseAddr(value)
	}

	return netip.Addr{}, fmt.Errorf("There is no IP (%s) in result: %s", jsonPath, body)
}

func GetClient() (*ovh.Client, error) {
	endpoint := os.Getenv("ENDPOINT")
	appKey := os.Getenv("APPLICATION_KEY")
	appSecret := os.Getenv("APPLICATION_SECRET")
	consumerKey := os.Getenv("CUSTOMER_KEY")

	return ovh.NewClient(endpoint, appKey, appSecret, consumerKey)
}

func CheckAndUpdate(publicIP netip.Addr, zone string, subDomain string) (bool, error) {
	client, err := GetClient()
	if err != nil {
		return false, err
	}

	dynHostRecord, err := findDynHostRecord(client, zone, subDomain)
	if err != nil {
		return false, err
	}

	log.WithFields(logrus.Fields{
		"ip":        publicIP,
		"zone":      zone,
		"subdomain": subDomain,
		"dynhostid": dynHostRecord.ID,
	}).Debug("Found DynHost record")

	if dynHostRecord.IP != publicIP.String() {
		if err := updateDynHostIP(client, zone, dynHostRecord.ID, publicIP); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func main() {
	var args struct {
		Zone      string `arg:"required"`
		SubDomain string `arg:"required"`
		IPUrl     string
		JsonPath  string
		IP        string
		Loop      bool
	}
	arg.MustParse(&args)

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	for {
		publicIP := netip.Addr{}
		var err error
		if args.IP != "" {
			publicIP, err = netip.ParseAddr(args.IP)
			if err != nil {
				log.Error(err)
			}
		} else {
			publicIP, err = GetIp(args.IPUrl, args.JsonPath)
			if err != nil {
				log.Error(err)
			}
			log.WithFields(logrus.Fields{
				"ip": publicIP,
			}).Debug("Found public IP")
		}

		ok, err := CheckAndUpdate(publicIP, args.Zone, args.SubDomain)
		if err != nil {
			log.Error("Error: %v\n", err)
		}
		if ok {
			log.Info("Updated\n")
		}

		if !args.Loop {
			break
		}

		time.Sleep(60 * time.Second)
	}
}
