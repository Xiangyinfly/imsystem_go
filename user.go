package main

import "net"

type User struct {
	Name string
	Addr string
	C chan string
	conn net.Conn
}
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C: make(chan string),
		conn: conn,
	}

	//启动一个go程，监听信息
	go user.ListenMessage()
	return user
}

func (u *User) ListenMessage() {
	for {
		//不断循环去读取chan中的信息并通过conn传递
		msg := <- u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}