package config

// mimics the tmpl package templater
type templater interface {
	Apply(string) (string, error)
}
