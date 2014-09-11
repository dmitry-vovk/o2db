package config

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type ConfigType struct {
	ListenTCP  string `json:"listen_tcp"`
	ListenHTTP string `json:"listen_http"`
	DataDir    string `json:"data_dir"`
	User       User   `json:"user"`
}
