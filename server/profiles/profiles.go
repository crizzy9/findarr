package profiles

// Profile represents a custom Findarr content profile (movies, shows, youtube, etc.)
type Profile struct {
	ID          int64  `db:"id"`
	Name        string `db:"name"`
	Type        string `db:"type"`
	Description string `db:"description"`
	Config      string `db:"config"` // JSON config for this profile
}

// ProfileManager manages profiles
type ProfileManager interface {
	CreateProfile(profile *Profile) error
	GetProfile(id int64) (*Profile, error)
	ListProfiles() ([]*Profile, error)
	DeleteProfile(id int64) error
}
