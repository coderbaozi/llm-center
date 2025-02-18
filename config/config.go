package config

type Config struct {
	Database struct {
		Driver    string `yaml:"driver"`
		Host      string `yaml:"host"`
		Port      int    `yaml:"port"`
		Username  string `yaml:"username"`
		Password  string `yaml:"password"`
		DBName    string `yaml:"dbname"`
		Charset   string `yaml:"charset"`
		ParseTime bool   `yaml:"parseTime"`
		Loc       string `yaml:"loc"`
	} `yaml:"database"`
	Server struct {
		Port int    `yaml:"port"`
		Mode string `yaml:"mode"`
	} `yaml:"server"`
}

