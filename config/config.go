package config

type Config struct {
	Server struct {
		Address string `json:"address"`
		Port    string `json:"port"`
		Debug   bool   `json:"debug"`
	} `json:"server"`
	Goka struct {
		Brokers            []string `json:"brokers"`
		TopicManagerConfig struct {
			TableReplication  int `json:"tableReplication"`
			StreamReplication int `json:"streamReplication"`
		} `json:"topicManagerConfig"`
		Topics []struct {
			Name      string `json:"name"`
			Partition int    `json:"partition"`
		} `json:"topics"`
	} `json:"goka"`
	Deposit DepositConfig `json:"deposit"`
}

type DepositConfig struct {
	Threshold struct {
		Amount float64 `json:"amount"`
		Time   string  `json:"time"`
	} `json:"threshold"`
}
