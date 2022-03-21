package server

// Config for the service
type Config struct {
	// Service specific config like name, ip, etc
	Service struct {
		// Name of the service
		Name string `yaml:"name"`

		// FqdnOrIP of the service
		FqdnOrIP string `yaml:"fqdnOrIP"`

		// ApiPort where service hosts the API server
		ApiPort string `yaml:"apiPort"`

		// RpcPort where service hosts the RPC server
		RpcPort string `yaml:"rpcPort"`
	} `yaml:"service"`

	// Logging details for the service
	Logging struct {
		// LogDir where logs are hosted
		LogDir string `yaml:"logDir"`

		// LogFile where logs are written
		LogFile string `yaml:"logFile"`

		// LoggingLevel for this service instance (info, error, fatal, etc)
		LoggingLevel string `yaml:"loggingLevel"`
	} `yaml:"logging"`

	// Datastore configuration for persisting service data
	Datastore struct {
		// FqdnOrIP of datastore
		FqdnOrIP string `yaml:"fqdnOrIP"`

		// Port where datastore is listening for incoming connections
		Port string `yaml:"port"`

		// Username of datastore
		Username string `yaml:"username"`

		// Password of datastore
		Password string `yaml:"password"`

		// DBName represents the database name where data is stored
		DbName string `yaml:"dbName"`
	} `yaml:"datastore"`

	// KVStore configuration
	KVStore struct {
		// FqdnOrIP of the kv store
		FqdnOrIP string `yaml:"fqdnOrIP"`

		// Port where kv store is listening for incoming connections
		Port string `yaml:"port"`
	} `yaml:"kvstore"`
}
