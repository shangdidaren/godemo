package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// 创建服务器
func NewServer(ip string, port int) *Server {
	server := Server{ip, port}
	return &server
}

func (re *Server) Handle(conn net.Conn) {
	// 处理当前业务
	fmt.Println("链接成功...")
}

// 启动服务器的接口
func (re *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", re.Ip, re.Port))
	if err != nil {
		fmt.Println("net.Listen err", err)
		return
	}
	defer listener.Close() // close litster socket

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Listen accept err:", err)
			continue
		}

		// do handle
		go re.Handle(conn)
	}

}
