package udp

import (
	"user-event-streamer/configs"
)

// Server is a class for keeping UDP protocol related methods
// UdpListenerChan is for passing incoming data from UDP connectionServer ...
type Server struct {
	cfg             *configs.UDP
	UdpListenerChan chan []byte
}

func New(config *configs.UDP) *Server {
	ch := make(chan []byte)
	return &Server{
		cfg:             config,
		UdpListenerChan: ch,
	}
}
