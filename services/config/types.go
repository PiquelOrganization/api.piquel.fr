package config

type EnvsConfig struct {
	PublicHost         string
	OrgDomain          string
	Host               string
	Port               string
	SSL                string
	DB_URL             string
	CookiesAuthSecret  string
	GoogleClientID     string
	GoogleClientSecret string
	GithubClientID     string
	GithubClientSecret string
	ConfigPath         string
}

type CORSConfig struct {
    AllowedOrigins []string
    // Duration of preflight caching
    MaxAge int
}

type Route struct {
	Name       string `yaml:"route"`
	Method     string
	Auth       bool
	Slug       bool
	Permission string
}
