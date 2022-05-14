package fileio

type Provider struct {
	FIOPath string
}

// NewFIOProvider creates a new FIO provider.
func NewFIOProvider(path string) *Provider {
	return &Provider{
		FIOPath: path,
	}
}
