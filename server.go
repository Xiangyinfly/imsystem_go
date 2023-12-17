package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

//服务端的构建

type Server struct {
	Ip string
	Port int
	OnlineMap map[string]*User
	maplock sync.RWMutex
	Message chan string
}
func NewServer(ip string,port int) *Server {
	server := &Server{
		Ip: ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}
	return server
}

//广播方法
func (s *Server) BroadCast(user *User,msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

//监听Message chan，一有消息发送给所有的user
func (s *Server) ListenMessager() {
	for {
		msg := <- s.Message
		s.maplock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.maplock.Unlock()
	}
}

func (s *Server) Handler(conn net.Conn) {
	//当前连接的业务
	//fmt.Println("连接成功！")
	user := NewUser(conn)
	//用户加入onlineMap
	s.maplock.Lock()
	s.OnlineMap[user.Name] = user
	s.maplock.Unlock()

	//广播该用户上线
	s.BroadCast(user,"已上线")

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				s.BroadCast(user,"下线")
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn read err:",err)
				return
			}
			msg := string(buf[:n-1])//去除\n
			s.BroadCast(user,msg)
		}
	}()

	//阻塞
	select{}
}

//启动服务的接口
func (s *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp",fmt.Sprintf("%s:%d",s.Ip,s.Port))
	if  err != nil {
		fmt.Println("net.Listen err:",err)
		return
	}
	defer listener.Close()

	//启动监听Message的go程
	go s.ListenMessager()

	//accept&do handler
	for {
		conn, err := listener.Accept()	
		if err != nil {
			fmt.Println("listener accept err:",err)
			continue
		}

		go s.Handler(conn)
	}


}