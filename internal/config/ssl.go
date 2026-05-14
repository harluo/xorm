package config

type SSL struct {
	CA string `json:"ca,omitempty"`
	Key string `json:"key,omitempty"`
	Cert string `json:"cert,omitempty"`
}
