package ocrmymail

type OCRMyMail struct {
	config *Config
}

// New creates a new OCRMyMail with the given config
func New(config *Config) (*OCRMyMail, error) {

	return &OCRMyMail{config: config}, nil

}
