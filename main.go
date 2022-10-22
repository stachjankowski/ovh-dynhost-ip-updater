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

// CheckAndUpdate checks if an IP address update is needed and updates it if so
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
		Zone      string `arg:"required" help:"domain"`
		SubDomain string `arg:"required" help:"subdomain"`
		IPUrl     string `arg:"--ipurl" help:"address to the service that returns your public ip address" default:"https://api.ipify.org"`
		JSONPath  string
		IP        string `arg:"--ip" help:"IPv4 address"`
		Loop      bool   `arg:"-l,--loop" help:"work in an infinite loop"`
		Verbose   bool   `arg:"-v,--verbose" help:"verbosity level"`
	}
	arg.MustParse(&args)

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	if args.Verbose {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}

	for {
		var publicIP netip.Addr
		var err error
		if args.IP != "" {
			publicIP, err = netip.ParseAddr(args.IP)
			if err != nil {
				log.Error(err)
			}
		} else {
			publicIP, err = iptools.GetIP(args.IPUrl, args.JSONPath)
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
