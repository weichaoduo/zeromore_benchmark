package main

import (
	 
   "fmt"
   "crypto/rand"
   "math/big"
   "time" 
   "github.com/antonholmquist/jason"  
   //"io/ioutil"
   "net"
   "os"
   "os/signal"
   "strconv"
   "runtime" 
   "bufio"
  // "log"
)


var urls = []string{ }
 

type SocketResponse struct {
  data      string
  response  string
  err       error
}

func main() {

    runtime.GOMAXPROCS(runtime.NumCPU()) 
    start:=time.Now().Unix()
    go end_hook( start )
    num, _ := strconv.ParseInt( os.Args[2], 10, 32)  
	for i:=0;i<int(num);i++{
		 
		max := big.NewInt(1000)
		rand, _ := rand.Int(rand.Reader, max) 
		urls  = append( urls,fmt.Sprintf( `{"token":"%d","req_id":"%d", "cmd":"socket.user_login","params":{"user":"admin_xbd","password":"258369"}}`  ,rand) )		
	}
    //log.Println( "urls:" ,urls  )
	results := asyncReq( urls )
  
    i:=0
	for _, result := range results {
        i++
		fmt.Printf( "%d result length: %d\n",i,  len(result.data) )
	}
  
    fmt.Printf( " result num: %d\n",  len(results) )
	 
}
 
 

func end_hook( start int64  ) {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, os.Kill)

    s := <-c
    end:=time.Now().Unix()
    els_time := end-start
    fmt.Println("Got signal:", s )
    fmt.Println("need time :", els_time )
    os.Exit( 1 )
    
}
 

func asyncReq( datas []string ) []*SocketResponse {

	ch := make( chan *SocketResponse )
	responses  := []*SocketResponse{}
	
	for _, data:= range datas {
		go func(data string) {
		
			//fmt.Printf("Req %s \n", data)
			times, _ := strconv.ParseInt( os.Args[3], 10, 32)  
			if len(os.Args) < 2 {
				fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
				os.Exit(1)
			}
			service := os.Args[1]
			tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
			checkError(err)
			conn, err := net.DialTCP("tcp", nil, tcpAddr)
            defer conn.Close()
			checkError(err) 
            time.Sleep(10 * time.Millisecond)
			_, err = conn.Write([]byte( data+"\n" ))
            
			checkError(err)
			var i int64
            i = 0
            sid:=""
            str :=""
			reader := bufio.NewReader(conn)
            for {
				
                msg, err := reader.ReadBytes( '\n' )
                if  err != nil { 
                    conn.Close() 
                    fmt.Println( "reader.ReadBytes: ", err.Error()  )  
                    break 
                }
				 
				//fmt.Printf( " response: %s\n", response_str )
				msg_json, errjson := jason.NewObjectFromBytes( msg  ) 
				checkError(errjson)
				cmd,  _ := msg_json.GetString("cmd") 
				//fmt.Printf( " cmd: %s\n", cmd )                    
				if cmd=="socket.user_login" { 
					sid,  _ = msg_json.GetString("data","sid")   
					//fmt.Printf( " sid: %s\n", sid )  
					str=fmt.Sprintf( `{ "cmd":"socket.getSession","params":{ "sid":"%s"  } }` ,sid )
					str =  str+"\n" 
					//fmt.Println( " post : ", str )  
					_,err = conn.Write([]byte( str  )) 
					//time.Sleep(1 * time.Millisecond )
					checkError(err)  
				}
				if cmd=="socket.getSession" {  
				   
					_,err = conn.Write([]byte(  str  )) 
					checkError(err) 
					i++
					
					//time.Sleep(1 * time.Millisecond )
					if( i>times ){
						conn.Close()
                        //fmt.Println( " i : ", i )  
						fmt.Println( " conn close! qps:" ,i ,"\n" )
						break
					}
				}
				 
                //ch <- &SocketResponse{ data, string(response_str), err} 
            }
            conn.Close()
		   
			
		}(data)
	}

	for {
		select {
		case r := <-ch:
			//fmt.Printf("%s was recved \n", r.data)
			responses = append(responses, r)
			if len(responses) == len(urls) {
				return responses
			}
		case <-time.After(500 * time.Millisecond):
          fmt.Printf(".")
		}
	}
	return responses
}

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}

 