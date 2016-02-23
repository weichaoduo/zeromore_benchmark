package connector

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"sync/atomic"
	"time"
	"zeromore_benchmark/global"
	"zeromore_benchmark/worker" 
)

/**
 * 监听客户端连接
 */
func Keeplive_connector(ip string, port int) {

	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(""), port, ""})
	if err != nil {
		fmt.Println("ListenTCP Exception:", err.Error())
		return
	}
	// 初始化 
    go statKick()
	listenAcceptTCP(listen)
}

/**
 *  处理客户端连接
 */
func listenAcceptTCP(listen *net.TCPListener) {

	//go stat_kick()

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			fmt.Println("AcceptTCP Exception::", err.Error())
			break
		}
		// 校验ip地址
		conn.SetKeepAlive(true)
		//log.Info("RemoteAddr:", conn.RemoteAddr().String())

		//remoteAddr :=conn.RemoteAddr()
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		sid := fmt.Sprintf("%d%d", r.Intn(99999), rand.Intn(999999))
        
        go handleConnWithBufferio( conn ,sid )

	} //end for {

}

func closeConn(conn *net.TCPConn, sid string) {

	//conn.Write([]byte{'E', 'O', 'F'})
	conn.Close()
	//log.Warn("Sid closing:", sid) 
	atomic.AddInt32(&global.SumConnections, -1)

}
 

func handleConnWithBufferio(conn *net.TCPConn, sid string ) {

	defer conn.Close()
	atomic.AddInt32(&global.SumConnections, 1)
	conn.SetNoDelay(false)
    fmt.Println("RemoteAddr:", conn.RemoteAddr().String() , "sid:", sid   )

	//声明一个管道用于接收解包的数据
	reader := bufio.NewReader(conn)
	for { 
		msg, err := reader.ReadBytes( '\n' )
		//fmt.Println(  "handleConn ReadString: ", string(msg) )
		if err != nil {
		    closeConn( conn , sid ) 
			fmt.Println( "HandleConn connection error: ", err.Error())
			break
		}
	    // 数据传递给worker
    		go worker.WorkerTask( conn, msg , sid )
  
    	}
}

 
func checkError(err error) {
	if err != nil {
		fmt.Println(os.Stderr, "Fatal error: %s", err.Error())
	}
}

func statKick() {

	timer := time.Tick(1000 * time.Millisecond)
	for _ = range timer {
		//ping := fmt.Sprintf(`{"cmd":"ping","ret":200,"time":%d }` , time.Now().Unix() );
		fmt.Println(time.Now().Unix(), " Connections: ", global.SumConnections)
		fmt.Println(time.Now().Unix(), " Qps: ", global.Qps)
	}
}

 
