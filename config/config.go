package config

type Config struct {
	App   App
	Log   Log
	Mysql Mysql
	Redis Redis
}

type App struct {
	Host string
	Port int
	Mode string
}

type Log struct {
	Path  string
	Name  string
	Level string
}

type Mysql struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Charset  string
}

type Redis struct {
	Host     string
	Port     int
	Password string
	Database int
}
