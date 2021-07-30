package main

import (
	"fmt"
	"io"
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
	fmt.Println(conn.RemoteAddr().String() + "链接成功...")
	// 用户上线了
	// 将用户加入到OnlineMap中，然后广播

	user := NewUser(conn,re)

	user.Online() //处理用户上线


	// 将当前请求阻塞
	//接受客户端发送的信息
	go func ()  {
		buf := make([]byte,4096)
		for {
		
			n,err := conn.Read(buf)
			if n==0{
				// n 为0 说明客户端主动的关闭了套接字连接
				user.Offline()
				return
			}
			if err !=nil  && err != io.EOF{
				fmt.Println("Conn Read Err：",err)
				return
			}
			msg := string(buf[:n-1])
			user.DoMessage(msg)  // 用户处理消息
		}
	}()
	
	
}

// 启动服务器的接口
func (re *Server) Start() {
	fmt.Println("服务器已启动...")
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
