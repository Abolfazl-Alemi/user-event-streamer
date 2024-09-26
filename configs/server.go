package configs

// Server keeps info for TCP, UDP, and HTTP connections
type Server struct {
	TCP  TCP  `yaml:"tcp"`  // TCP: for receiving leads data
	UDP  UDP  `yaml:"udp"`  // UDP: for receiving views data
	HTTP HTTP `yaml:"HTTP"` // HTTP: for serving monitoring data
}

// TCP  for receiving leads data
// listen to port 8080
type TCP struct {
	Host string `yaml:"host" envconfig:"TCP_HOST"`
	Port string `yaml:"port" envconfig:"TCP_PORT"`
}

// UDP for receiving views data
// listen to port 5000
type UDP struct {
	Host string `yaml:"host" envconfig:"UDP_HOST"`
	Port string `yaml:"port" envconfig:"UDP_PORT"` // 5000
}

// HTTP for serving monitoring data
type HTTP struct {
	Host        string `yaml:"host" envconfig:"HTTP_HOST"`
	Port        string `yaml:"port" envconfig:"HTTP_PORT"`                   // Port exposed on 8001
	MetricsPath string `yaml:"metrics_path" envconfig:"HTTP_METRICS_PATH"`   // MetricsPath: for serving modules metrics
	AppInfoPath string `yaml:"app_info_path" envconfig:"HTTP_APP_INFO_PATH"` // AppInfoPath: for serving service generic metrics
}
