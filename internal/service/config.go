package service

type Config struct {
	Repositories []string `yaml:"git_repositories"` //absolute paths
	QSName       string
}
