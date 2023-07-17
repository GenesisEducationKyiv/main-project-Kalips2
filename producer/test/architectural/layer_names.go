package architectural

var (
	domainLayer         = []string{domain}
	applicationLayer    = []string{application}
	infrastructureLayer = []string{provider, repository, sender, logger}
	presenterLayer      = []string{handler, presenter}
)
