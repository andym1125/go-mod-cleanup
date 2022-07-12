package main

/* TRUNCATE_THRESHOLD is the threshold after which nodes are truncated for being too long. */
var TRUNCATE_THRESHOLD int = 10

/* DIR_NAME is the name of the directory to put all client-side files */
var DIR_NAME string = "gomodgraph"

type Config struct {
	TruncateThreshold int    `json:"truncateThreshold"`
	DirName           string `json:"dirName"`
}

func GetConfig() Config {
	return Config{
		TruncateThreshold: TRUNCATE_THRESHOLD,
		DirName:           DIR_NAME,
	}
}

func LoadConfig(config *Config) {
	TRUNCATE_THRESHOLD = config.TruncateThreshold
	DIR_NAME = config.DirName
}
