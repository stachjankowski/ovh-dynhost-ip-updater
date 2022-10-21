package ovhdynhost

import (
	"net/netip"
	"os"
	"testing"

	"github.com/ovh/go-ovh/ovh"
	"github.com/stretchr/testify/assert"
)

func Environments(endpoint string, f func()) {
	os.Setenv("OVH_ENDPOINT", endpoint)
	os.Setenv("OVH_CONSUMER_KEY", "a")
	os.Setenv("OVH_APPLICATION_KEY", "b")
	os.Setenv("OVH_APPLICATION_SECRET", "c")
	f()
	defer os.Unsetenv("OVH_ENDPOINT")
	defer os.Unsetenv("OVH_CONSUMER_KEY")
	defer os.Unsetenv("OVH_APPLICATION_KEY")
	defer os.Unsetenv("OVH_APPLICATION_SECRET")
}

func TestGetClientByEnv(t *testing.T) {
	Environments("ovh-eu", func() {
		client, err := GetClient()
		assert.Equal(t, nil, err)
		assert.NotEqual(t, (*ovh.Client)(nil), client)
	})
}

func TestGetClientInvalidEndpoint(t *testing.T) {
	Environments("xyz", func() {
		client, err := GetClient()
		assert.NotEqual(t, nil, err)
		assert.Equal(t, (*ovh.Client)(nil), client)
	})
}

func TestFindDynHostRecordInvalidAppKey(t *testing.T) {
	Environments("ovh-eu", func() {
		client, _ := GetClient()

		_, err := FindDynHostRecord(client, "example.com", "subdomain")

		assert.NotEqual(t, nil, err)

		apierror, ok := err.(*ovh.APIError)
		assert.Equal(t, true, ok)
		assert.Equal(t, 403, apierror.Code)
	})
}

func TestUpdateDynHostIPInvalidAppKey(t *testing.T) {
	Environments("ovh-eu", func() {
		client, _ := GetClient()

		ip, _ := netip.ParseAddr("127.0.0.1")
		err := UpdateDynHostIP(client, "example.com", 123, ip)

		assert.NotEqual(t, nil, err)

		apierror, ok := err.(*ovh.APIError)
		assert.Equal(t, true, ok)
		assert.Equal(t, 403, apierror.Code)
	})
}
