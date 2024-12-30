package vanity

import "errors"

var (
	ErrNoRepository      = errors.New("a repository url is required")
	ErrInvalidRepository = errors.New("expected repository url in the form vcsScheme://vcsHost/user/repo")
	ErrInvalidProtocol   = errors.New("protocol must be git, github, or gogs")
)
