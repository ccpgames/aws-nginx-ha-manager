package monitor

// Balancer encapsulates an ELB and contains functions to query it's state
type Balancer struct {
	elbName     string
	resolver Resolver
}

// NewBalancer returns a new balancer instance
func NewBalancer(resolver Resolver, elbName string) *Balancer {
	balancer := Balancer{
		elbName:     elbName,
		resolver: resolver,
	}
	return &balancer
}

// GetIPList returns a string array containing resolved ips
func (b *Balancer) GetIPList() ([]string, error) {
	list, err := b.resolver.Resolve(b.elbName)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// IsHealthy returns true if there are healthy hosts
func (b *Balancer) IsHealthy() (bool, error) {
	list, err := b.GetIPList()
	if err != nil {
		return false, err
	}
	return len(list) > 0, nil
}
