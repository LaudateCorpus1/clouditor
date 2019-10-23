package discovery

// Asset is the main concepts behind components or other things within the service that are discoverable
type Asset struct {
	Name        string
	Description string
}

type DiscoveryResult struct {
}

type Discoverable interface {
	Name() string
	ID() string
}

type Scanner interface {
	Scan() DiscoveryResult
}
