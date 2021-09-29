package lib

type Option interface {
	apply(cfg *config) error
}

type config struct {
	alphaSortTypes bool
	alphaSortFuncs bool
}

func AlphaSortTypes() Option { return alphaSortTypes{} }

type alphaSortTypes struct{}

func (o alphaSortTypes) apply(cfg *config) error {
	cfg.alphaSortTypes = true
	return nil
}

func AlphaSortFuncs() Option { return alphaSortFuncs{} }

type alphaSortFuncs struct{}

func (o alphaSortFuncs) apply(cfg *config) error {
	cfg.alphaSortFuncs = true
	return nil
}

func parseOptions(opts []Option) (config, error) {
	cfg := config{}
	for _, opt := range opts {
		if err := opt.apply(&cfg); err != nil {
			return cfg, err
		}
	}
	return cfg, nil
}
