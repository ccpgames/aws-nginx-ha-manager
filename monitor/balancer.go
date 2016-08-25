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
