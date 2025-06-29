package downloaders

// Downloader is an interface for download agents (yt-dlp, torrent, etc.)
type Downloader interface {
	ID() string
	Name() string
	Type() string // e.g. 'yt-dlp', 'torrent', etc.
	Download(sourceURL string, options map[string]any) error
}

var RegisteredDownloaders = make(map[string]Downloader)

// Register a downloader
func RegisterDownloader(d Downloader) {
	RegisteredDownloaders[d.ID()] = d
}
