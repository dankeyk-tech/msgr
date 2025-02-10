package model

type Config struct {
	DB   DB   `json:"db"`
	HTTP Http `json:"http"`
	JWT  JWT  `json:"jwt"`
}

type DB struct {
	Scheme   string `json:"scheme"`
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Path     string `json:"path"`
}

type Http struct {
	Port string `json:"port"`
}

type JWT struct {
	Domain string `json:"domain"`
	Key    string `json:"key"`
}
