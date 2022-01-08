package config

type Server struct {
	HTTP HTTP `yaml:"http"`
}

type HTTP struct {
	Address string `yaml:"address"`
}

type Database struct {
	Driver          string `yaml:"driver"`
	MaxOpenConns    string `yaml:"max_open_conns"`
	MaxIdleConns    string `yaml:"max_idle_conns"`
	ConnMaxLifetime string `yaml:"conn_max_lifetime"`
	ConnMaxIdletime string `yaml:"conn_max_idletime"`

	MasterDsn         string
	FollowerDsn       string
	MasterDsnNoCred   string `yaml:"master_dsn_no_cred"`
	FollowerDsnNoCred string `yaml:"follower_dsn_no_cred"`
}

type Redis struct {
	Engine        string `yaml:"engine"`
	MaxActiveConn string `yaml:"max_active_conn"`
	MaxIdleConn   string `yaml:"max_idle_conn"`
	Timeout       string `yaml:"timeout"`
	Address       string `yaml:"address"`
}
