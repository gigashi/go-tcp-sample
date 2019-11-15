package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("listen :8000")
	listen, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Printf("accept %s\n", conn.RemoteAddr().String())

		buf := make([]byte, 4*1024) // TODO: define buf size
		for {
			readCount, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("read error:%s\n", err.Error())
				break
			}
			fmt.Println("__________________________________________________")
			fmt.Println(string(buf[:readCount]))
		}
		conn.Close()
		fmt.Println("__________________________________________________")
		fmt.Println("disconnect")
	}
}
