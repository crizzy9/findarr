package indexers

// Indexer is an interface for indexer/search providers (torrent, usenet, etc.)
type Indexer interface {
	ID() string
	Name() string
	Type() string // e.g. 'torrent', 'youtube', etc.
	Search(query string) ([]Result, error)
}

type Result struct {
	Title    string
	Type     string
	SourceID string
	Meta     map[string]any
}

var RegisteredIndexers = make(map[string]Indexer)

// Register an indexer
func RegisterIndexer(indexer Indexer) {
	RegisteredIndexers[indexer.ID()] = indexer
}

