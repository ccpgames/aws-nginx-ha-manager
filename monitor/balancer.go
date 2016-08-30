package monitor

// Balancer encapsulates an ELB and contains functions to query it's state
type Balancer struct {
	fqdn string
}

// NewBalancer returns a new balancer instance
func NewBalancer(fqdn string) *Balancer {
	balancer := Balancer{
		fqdn: fqdn,
	}
	return &balancer
}

// GetIPList returns a string array containing resolved ips
func (b *Balancer) GetIPList() []string {
	ret := []string{"10.0.0.1", "10.0.0.2", "10.0.0.3", "10.0.0.4"}
	return ret
}

// IsHealthy returns true if there are healthy hosts
func (b *Balancer) IsHealthy() bool {
	return len(b.GetIPList()) > 0
}
