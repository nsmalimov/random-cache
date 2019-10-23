package config

// read from yaml if need
type Config struct {
	BindAddress            string
	BindPort               int
	LenStringForAddToCache int
	FrequencyAddToCacheSec int
	EndPointStr            string
}
