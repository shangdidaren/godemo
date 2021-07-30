package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// 在用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息管道
	Message chan string
}

// 创建服务器
func NewServer(ip string, port int) *Server {
	server := Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: map[string]*User{},
		mapLock:   sync.RWMutex{},
		Message:   make(chan string),
	}
	return &server
}

//监听Message的go程,也是Server启动时就开始监听
func (re *Server) ListenMessage() {
	for {
		msg := <-re.Message
		re.mapLock.Lock() // 此时上线的client
		for _, cli := range re.OnlineMap {
			cli.C <- msg
		}
		re.mapLock.Unlock()
	}
}

// 广播消息
func (re *Server) BroadCast(user *User, msg string) {
	sendText := "[" + user.Addr + "] " + user.Name + ":" + msg
	re.Message <- sendText // 这里网Message发信息了，还需要一个监听Message的go程
}

func (re *Server) Handle(conn net.Conn) {
	// 处理当前业务
	fmt.Println("链接成功...")
	// 用户上线了
	// 将用户加入到OnlineMap中，然后广播

	user := NewUser(conn)
	re.mapLock.Lock()
	re.OnlineMap[user.Name] = user
	re.mapLock.Unlock()

	// 广播上线信息
	re.BroadCast(user, "已上线")
	select {}

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


	go re.ListenMessage()
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
