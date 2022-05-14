package shortlinks

// Shortlink redirects a user from /From to To.
type Shortlink struct {
	From, To string
}

// History represents a given version of a Shortlink.
type History struct {
	From, To, When, Who string
}

// DB is used by the Server to store shortlinks and related history.  May
// optionally be a DBDeleted.
type DB interface {
	PublicDB

	// CreateShortlink inserts or updates Shortlink.
	CreateShortlink(Shortlink) error

	// DeleteShortlink deletes a shortlink from the database.
	DeleteShortlink(from, who string) error

	// History loads history for a given shortlink.  Hardcoding a nil
	// return value is supported.
	History(from string) ([]History, error)

	// InsertHistory stores the history for a newly inserted/updated
	// shortlink.  Hardcoding a nil return value is supported.
	InsertHistory(History) error
}

type PublicDB interface {
	// Shortlink loads data for from.
	Shortlink(from string) (Shortlink, error)

	// AllShortlinks loads a list of shortlinks for use in the index.
	// Hardcoding a nil return value is supported.
	AllShortlinks() ([]Shortlink, error)
}

type DBDeleted interface {
	// DeletedShortlinks returns all deleted shortlinks.
	DeletedShortlinks() ([]Shortlink, error)
}
