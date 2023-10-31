package listen

import (
	"log"
	"net"
)

func HandleListener(listener net.Listener) {
	//blockChan := make(chan os.Signal, 1)
	//signal.Notify(blockChan, os.Interrupt, syscall.SIGUSR1)
	//listen, err := net.Listen("tcp", "0.0.0.0:9998")
	//if err != nil {
	//	panic(err)
	//}
	go func() {
		for {
			_, err := listener.Accept()
			if err != nil {
				log.Printf("accept error: %v\n", err)
				continue
			}

		}
	}()

}
