package server

type Config struct {
	Service struct {
		Name     string `yaml:"name"`
		FqdnOrIP string `yaml:"fqdnOrIP"`
		ApiPort  string `yaml:"restPort"`
		RpcPort  string `yaml:"rpcPort"`
	} `yaml:"service"`
	Logging struct {
		LogDir       string `yaml:"log_dir"`
		LogFile      string `yaml:"log_file"`
		LoggingLevel string `yaml:"logging_level"`
	} `yaml:"logging"`
	SQLDb struct {
		FqdnOrIP string `yaml:"fqdnOrIP"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DBName   string `yaml:"db_name"`
	} `yaml:"sql_db"`
	KVStore struct {
		FqdnOrIP string `yaml:"fqdnOrIP"`
		Port     string `yaml:"port"`
	} `yaml:"kvstore"`
}
