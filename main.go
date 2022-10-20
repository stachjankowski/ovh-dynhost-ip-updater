package main

import (
	"net/netip"
	"time"

	arg "github.com/alexflint/go-arg"
	"github.com/sirupsen/logrus"
	"github.com/stachjankowski/ovh-dynhost-ip-updater/iptools"
	"github.com/stachjankowski/ovh-dynhost-ip-updater/ovhdynhost"
)

var log = logrus.New()

func CheckAndUpdate(publicIP netip.Addr, zone string, subDomain string) (bool, error) {
	client, err := ovhdynhost.GetClient()
	if err != nil {
		return false, err
	}

	dynHostRecord, err := ovhdynhost.FindDynHostRecord(client, zone, subDomain)
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
		if err := ovhdynhost.UpdateDynHostIP(client, zone, dynHostRecord.ID, publicIP); err != nil {
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
			publicIP, err = iptools.GetIP(args.IPUrl, args.JsonPath)
			if err != nil {
				log.Error(err)
			}
			log.WithFields(logrus.Fields{
				"ip": publicIP,
			}).Debug("Found public IP")
		}

		ok, err := CheckAndUpdate(publicIP, args.Zone, args.SubDomain)
		if err != nil {
			log.Error(err)
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
