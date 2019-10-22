package discovery

// Asset is the main concepts behind components or other things within the service that are discoverable
type Asset struct {
	Name        string
	Description string
}

type Named interface {
	GetName() string
	GetID() string
}
