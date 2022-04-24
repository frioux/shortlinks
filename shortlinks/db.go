package shortlinks

// Shortlink redirects a user from /From to To.
type Shortlink struct {
	From, To string
}

// History represents a given version of a Shortlink.
type History struct {
	From, To, When string
}

// DB is used by the Server to store shortlinks and related history.
type DB interface {
	// Shortlink loads data for from.
	Shortlink(from string) (Shortlink, error)

	// CreateShortlink inserts or updates Shortlink.
	CreateShortlink(Shortlink) error

	// AllShortlinks loads a list of shortlinks for use in the index.
	// Hardcoding a nil return value is supported.
	AllShortlinks() ([]Shortlink, error)

	// History loads history for a given shortlink.  Hardcoding a nil
	// return value is supported.
	History(from string) ([]History, error)

	// InsertHistory stores the history for a newly inserted/updated
	// shortlink.  Hardcoding a nil return value is supported.
	InsertHistory(History) error
}
