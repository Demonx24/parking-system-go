package config

type Config struct {
	Mysql  Mysql  `yaml:"mysql" json:"mysql"`
	System System `yaml:"system" json:"system"`
	Zap    Zap    `yaml:"zap" json:"zap"`
}
