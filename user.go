package main

import "net"

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// 封装用户上线功能
func (re *User) Online() {
	re.server.mapLock.Lock()
	re.server.OnlineMap[re.Name] = re
	re.server.mapLock.Unlock()
	// 广播上线信息
	re.server.BroadCast(re, "已上线")
}

// 封装用户下线功能

func (re *User) Offline() {
	re.server.mapLock.Lock()
	delete(re.server.OnlineMap, re.Name) // 删除字典中的元素:--> 上线列表中删除
	re.server.mapLock.Unlock()
	re.server.BroadCast(re, "下线")
}

// 封装用户处理消息功能
func (re *User) DoMessage(msg string) {
	// 这里的处理消息是借助server的广播  暂时！！
	re.server.BroadCast(re,msg)
}

//
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := User{
		userAddr,
		userAddr,
		make(chan string),
		conn,
		server,
	}
	go user.ListenMessage()

	return &user
}

// 监听信息 ,当用户启动时就开始监听，所以在创建时就调用此方法
func (re *User) ListenMessage() {
	for {
		msg := <-re.C // 一旦有信息就发送给client

		re.conn.Write([]byte(msg + "\n")) // send
	}
}
