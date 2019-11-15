package main

import (
	"fmt"
	"net"
	"sync"
)

func connectLoop(port string, otherConnectReceive <-chan []byte, otherConnectSend chan<- []byte) {
	for {
		listen, err := net.Listen("tcp", port)
		if err != nil {
			panic(err)
		}
		fmt.Printf("wait port[%s]\n", port)

		// コネクションをchanで待てるようにする
		var conn net.Conn
		connWait := make(chan net.Conn)
		go func() {
			for {
				c, err := listen.Accept()
				if err != nil {
					fmt.Printf("accept error:%s\n", err.Error())
					continue
				}
				connWait <- c
				break
			}
		}()

	L1: // コネクション待ち時に別の接続から来たデータは捨てる
		for {
			select {
			case conn = <-connWait:
				break L1
			case _ = <-otherConnectReceive:
			}
		}

		fmt.Printf("connect port[%s]\n", port)
		// 接続されたらlistenを辞める(1接続のみにする)
		listen.Close()

		// 受信ループ
		receive := make(chan []byte)
		go func() {
			buf := make([]byte, 4*1024) // TODO: define buf size
			for {
				readCount, err := conn.Read(buf)
				if err != nil {
					fmt.Printf("read error:%s\n", err.Error())
					break
				}
				receiveData := make([]byte, readCount)
				copy(receiveData, buf[0:readCount-1])
				receive <- receiveData
			}
			conn.Close()
			close(receive)
		}()

	L2: //
		for {
			select {
			case v, ok := <-receive:
				if !ok {
					break L2
				}
				fmt.Println(string(v))
				otherConnectSend <- v
			case v := <-otherConnectReceive:
				for writeCount := 0; writeCount < len(v); {
					c, err := conn.Write(v[writeCount:])
					if err != nil {
						fmt.Printf("write error:%s\n", err.Error())
						conn.Close()
						break L2
					}
					writeCount += c
				}
			}
		}
		fmt.Printf("disconnected port[%s]\n", port)
	}
}

func main() {
	a := make(chan []byte)
	b := make(chan []byte)

	go connectLoop(":8000", a, b)
	go connectLoop(":8001", b, a)

	// mainを止めておく
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
