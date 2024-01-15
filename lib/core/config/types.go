package config

type Config struct {
	Project struct {
		Environment string `yaml:"environment"`
	} `yaml:"project"`
	Nats struct {
		Url      string `yaml:"url"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"nats"`
}
