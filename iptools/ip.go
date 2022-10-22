package iptools

import (
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"strings"

	"github.com/spyzhov/ajson"
)

// GetIP gets your public address from an external service.
func GetIP(url string, jsonpath string) (netip.Addr, error) {
	req, err := http.Get(url)
	if err != nil {
		return netip.Addr{}, err
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return netip.Addr{}, err
	}

	// raw format
	if jsonpath == "" {
		return netip.ParseAddr(strings.TrimRight(string(body), "\n"))
	}

	// json
	root, err := ajson.Unmarshal(body)
	if err != nil {
		return netip.Addr{}, err
	}
	nodes, err := root.JSONPath(jsonpath)
	if err != nil {
		return netip.Addr{}, err
	}

	if len(nodes) > 0 {
		value, err := nodes[0].GetString()
		if err != nil {
			return netip.Addr{}, err
		}
		return netip.ParseAddr(value)
	}

	return netip.Addr{},
		fmt.Errorf("There is no IP (%s) in result: %s", jsonpath, body)
}
