package iptools

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIP(t *testing.T) {
	var server *httptest.Server

	tests := map[string]struct {
		response string
		jsonPath string
		ip       netip.Addr
		ok       bool
	}{
		"raw":                   {"127.0.0.2", "", netip.AddrFrom4([4]byte{127, 0, 0, 2}), true},
		"raw-newline":           {"127.0.0.2\n", "", netip.AddrFrom4([4]byte{127, 0, 0, 2}), true},
		"json":                  {"{\"ip\": \"127.0.0.2\"}", "$.ip", netip.AddrFrom4([4]byte{127, 0, 0, 2}), true},
		"invalid-json":          {"{\"ip\": \"127.0.0.2\"", "$.ip", netip.Addr{}, false},
		"json-invalid-jsonPath": {"{\"ip\": \"127.0.0.2\"}", "$.ip2", netip.Addr{}, false},
		"json-no-jsonPath":      {"{\"ip\": \"127.0.0.2\"}", "", netip.Addr{}, false},
		"json-invalid-type":     {"{\"ip\": {}}", "$.ip", netip.Addr{}, false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, tt.response)
			}))

			ip, err := GetIP(server.URL, tt.jsonPath)
			if tt.ok {
				assert.Empty(t, err)
			} else {
				assert.NotEmpty(t, err)
			}
			assert.Equal(t, tt.ip, ip)

			defer server.Close()
		})
	}
}
