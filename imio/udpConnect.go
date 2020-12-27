package imio

import (
	"log"
	"net"
	"os"
)

func startUdpListen(port string )  {
	adr,err := net.ResolveUDPAddr("udp",port)
	if err != nil{
		log.Println(err)
		os.Exit(1)
	}
	conn,err :=net.ListenUDP("udp",adr)
	if err != nil{
		log.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	go udpRead(conn)
}

func udpRead(conn *net.UDPConn)  {

}