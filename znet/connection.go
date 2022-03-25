package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

//链接模块
type Connection struct {
	//当前连接的socket TCP套接字
	Coon *net.TCPConn
	//	连接的id
	CoonID uint32
	//	当前的连接状态
	isClosed bool
	// 当前连接所绑定的处理业务方法API
	handleAPI ziface.HandleFunc
	//	告知当前连接已经退出的/停止channel
	ExitChan chan bool
}

//初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, callback_api ziface.HandleFunc) *Connection {
	c := &Connection{
		Coon:      conn,
		CoonID:    connID,
		isClosed:  false,
		handleAPI: callback_api,
		ExitChan:  make(chan bool, 1),
	}
	return c
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println("connID= ", c.CoonID, "Reader is exit,remote addr is", c.RemoteAddr().String())
	defer c.Stop()
	for {
		//	读取客户端的数据到buf中，最大512
		buf := make([]byte, 512)
		cnt, err := c.Coon.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			continue
		}
		//	调用当前连接所绑定的handleAPI
		err = c.handleAPI(c.Coon, buf, cnt)
		if err != nil {
			fmt.Println("ConnID", c.CoonID, "handle is error", err)
			break
		}
	}

}

//	启动连接 让当前连接开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start(),CoonID = ", c.CoonID)
	//	启动从当前连接的读数据的业务
	go c.StartReader()
	// TODO 启动从当前连接写数据的业务

}

//	停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop(),CoonID = ", c.CoonID)
	//	如果当前连接已经关闭
	if c.isClosed {
		return
	}
	c.isClosed = false
	//关闭socket连接
	err := c.Coon.Close()
	if err != nil {
		fmt.Println("socket close err:", err)
		return
	}
	//关闭管道，回收资源
	close(c.ExitChan)
}

//	获取当前连接的绑定socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Coon
}

//	获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.CoonID
}

//	获取远程客户端的TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Coon.RemoteAddr()
}

//	发送数据，将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}
