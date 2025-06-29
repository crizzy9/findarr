package metadata

// Provider is the interface for metadata resolvers.
type Provider interface {
	Search(query string) ([]MediaResult, error)
}

// MediaResult represents a found media item.
type MediaResult struct {
	Title string
	Type  string // movie, book, music, etc.
	ID    string
}

