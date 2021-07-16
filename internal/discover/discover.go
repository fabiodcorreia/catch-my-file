// Discover package is responsible for the components that allow
// the peer discovery on the local network.
//
// It contains a server that register the peers to be discovered
// and the client that will discover the peers
package discover

const (
	serviceName   = "_catchmyfile._tcp"
	serviceDomain = "local."
)
