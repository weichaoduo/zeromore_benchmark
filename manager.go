//
//  main
//

package main

import (
	
	"runtime"
	"zeromore_benchmark/connector"
)

  
/**
 * 程序入口
 */
func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
 
    go connector.Keeplive_connector("*", 7002)
 
	select {}

}



