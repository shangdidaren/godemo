package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

//
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := User{
		userAddr,
		userAddr,
		make(chan string),
		conn,
	}
	go user.ListenMessage()

	return &user
}


// 监听信息 ,当用户启动时就开始监听，所以在创建时就调用此方法
func (re * User)ListenMessage(){
	for{
		msg := <-re.C // 一旦有信息就发送给client
		
		re.conn.Write([]byte(msg+"\n")) // send
	}
}
