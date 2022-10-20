package iptools

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/netip"
	"strings"

	"github.com/spyzhov/ajson"
)

func GetIP(url string, jsonPath string) (netip.Addr, error) {
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
		return netip.ParseAddr(strings.TrimRight(string(body), "\n"))
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
