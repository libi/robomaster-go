package connection

import (
	"fmt"
	"sync"
	"testing"
)

func TestNewRoboMasterConn(t *testing.T) {
	wg := new(sync.WaitGroup)
	roboMasterConn, err := NewRoboMasterConn(&Option{EnableVideo: true})
	r, err := roboMasterConn.RunCmd("robot battery ?")

	fmt.Println(r, err)
	wg.Add(1)
	go func() {
		for {
			buff := make([]byte, 1024)
			n, err := roboMasterConn.VideoConn.Read(buff)
			fmt.Println("rec video stream", buff[0:n], err)
		}
		wg.Done()
	}()

	wg.Wait()
}
