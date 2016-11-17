package requesthandler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubdomainDetection(t *testing.T) {

	hostnames := []string{"localhost", "192.168.0.103", "mydomain.manydomains.superdomain.co.uk"}

	assert.Equal(t, "", getSubdomainFromHost("localhost:3030", hostnames))
	assert.Equal(t, "sub", getSubdomainFromHost("sub.localhost:3030", hostnames))
	assert.Equal(t, "sub", getSubdomainFromHost("sub.localhost", hostnames))
	assert.Equal(t, "sub-hyphen", getSubdomainFromHost("sub-hyphen.localhost:30", hostnames))
	assert.Equal(t, "sub.sub", getSubdomainFromHost("sub.sub.localhost:3030", hostnames))

	// These are not really valid subdomains, but the test should still pass.
	assert.Equal(t, "sub", getSubdomainFromHost("sub.192.168.0.103:3030", hostnames))
	assert.Equal(t, "sub", getSubdomainFromHost("sub.192.168.0.103", hostnames), "sub")
	assert.Equal(t, "sub-domain", getSubdomainFromHost("sub-domain.mydomain.manydomains.superdomain.co.uk", hostnames))
}
