package server

type Config struct {
	Service struct {
		Name     string `yaml:"name"`
		FqdnOrIP string `yaml:"fqdnOrIP"`
		ApiPort  string `yaml:"apiPort"`
		RpcPort  string `yaml:"rpcPort"`
	} `yaml:"service"`
	Logging struct {
		LogDir       string `yaml:"logDir"`
		LogFile      string `yaml:"logFile"`
		LoggingLevel string `yaml:"loggingLevel"`
	} `yaml:"logging"`
	SQLDb struct {
		FqdnOrIP string `yaml:"fqdnOrIP"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbName"`
	} `yaml:"sql_db"`
	KVStore struct {
		FqdnOrIP string `yaml:"fqdnOrIP"`
		Port     string `yaml:"port"`
	} `yaml:"kvstore"`
}
