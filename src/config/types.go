package config

type ConfigType struct {
	ListenTCP  string `json:"listen_tcp"`
	ListenHTTP string `json:"listen_http"`
	DataDir    string `json:"data_dir"`
}
