package znet

import (
	"errors"
	"fmt"
	"net"
)

// iServer的接口实现，定义了一个Server的服务器模块
type Server struct {
	// 服务器的名称
	Name string
	// 服务器绑定的ip版本
	IPVersion string
	// 服务器监听的ip
	IP string
	// 服务器监听的端口
	Port int
}

// 这里写死一个业务方法，目前写死了handle，以后修改
//定义当前客户端连接的所绑定的handle api
func CallBackClient(conn *net.TCPConn, data []byte, cnt int) error {
	//	回显业务
	fmt.Println("[Conn Handle]CallBackToClient...")
	_, err := conn.Write(data[:cnt])
	if err != nil {
		fmt.Println("write back buf err", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[start]Server Listenner at IP:%s,port:%d is starting\n", s.IP, s.Port)
	go func() {
		// 1.获取一个TCP的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp address error:", err)
			return
		}

		// 2.监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
		}
		fmt.Println("Start Zinx server succ, ", s.Name, "success,Listening")
		// 分配一个connID
		var cid uint32 = 0

		// 3.阻塞的等待客户端链接，处理客户端链接业务(读写)
		for {
			//	如果有客户端链接，处理客户端链接业务(读写)
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			//	--已经与客户端建立了链接，做一些业务，做一个最基本的最大512字节长度的回显业务--
			//	使用connection模块
			// 将处理新连接的业务方法和conn进行绑定，得到连接模块
			dealConn := NewConnection(conn, cid, CallBackClient)
			cid++
			//	启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	// TODO 将服务器的资源、状态和资源停止回收
}

//  运行服务器
func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()
	// TODO 做一些启动服务器之后的额外业务

	// 阻塞
	select {}
}

// 初始化Server模块的方法
func NewServer(name string) *Server {
	//由于方法为指针类型
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
