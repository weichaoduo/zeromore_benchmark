//
//  worker
//

package worker

import (
	"github.com/antonholmquist/jason"
	"sync/atomic"
	"zeromore_benchmark/global"
	"fmt"
	"net"
)


func WorkerTask( conn *net.TCPConn, task_data []byte , sid string  ) {

    worker_json, _ := jason.NewObjectFromBytes( task_data )
 
	//  do some thing
	cmd, _ := worker_json.GetString("cmd")
	//fmt.Printf(" worker_task logic cmd: %s\n", cmd)

	json := fmt.Sprintf(`{"cmd":"%s","data":"%s"}`, cmd, sid )
	//fmt.Println(" WorkerTask : ", json)
	conn.Write([]byte(json + "\n"))
	atomic.AddInt64(&global.Qps, 1)
	
}

 








