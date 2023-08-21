package config

type SshConfig struct {
	Hostname string `json:"hostname"`
	Password string `json:"password"`
	Address  string `json:"address"`
	Port     string `json:"port"`
	Rows     uint32 `json:"rows"`
	Cols     uint32 `json:"cols"`
}
