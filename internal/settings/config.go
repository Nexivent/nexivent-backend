package settings

type Config struct {
	Port int
	Env  string
	Db struct {
		URL string
	}
}
