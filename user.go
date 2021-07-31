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


// 返回给当前用户信息
func (re *User) SendMsg(msg string) {
	re.conn.Write([]byte(msg))
}

// 封装用户处理消息功能
func (re *User) DoMessage(msg string) {

	// 协程 ———————— 规则，如果发送who ，就返回给当前用户所有在线用户信息
	if msg == "who" {
		re.server.mapLock.Lock()
		for _, user := range re.server.OnlineMap {
			re.SendMsg("[" + user.Addr + "] " + user.Name + "在线")
		}
		re.server.mapLock.Unlock()

	} else {
		re.server.BroadCast(re, msg)
	}

}

// 创建一个用户
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
