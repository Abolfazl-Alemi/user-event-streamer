package udp

import (
	"fmt"
	"go.uber.org/zap"
	"net"
	"sync"
	"user-event-streamer/pkg/logger"
)

func (sr *Server) Run(wg *sync.WaitGroup) {
	listener, err := net.ListenPacket("udp", fmt.Sprintf("%s:%v", sr.cfg.Host, sr.cfg.Port))
	if err != nil {
		logger.Zap.Fatal("error in establishing udp server", zap.Error(err))
	}
	defer func(listener net.PacketConn) {
		err := listener.Close()
		if err != nil {
			logger.Zap.Error("error in closing udp server", zap.Error(err))
		}
	}(listener)

	defer close(sr.UdpListenerChan)

	for {
		buffer := make([]byte, 1024)
		n, _, err := listener.ReadFrom(buffer)
		if err != nil {
			logger.Zap.Error("error in listening udp server operation", zap.Error(err))
			continue
		}

		go sr.handleRequest(buffer[:n], wg)
	}
}

func (sr *Server) handleRequest(buffer []byte, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	sr.UdpListenerChan <- buffer
}
