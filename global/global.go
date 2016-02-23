/**
 *  全局变量
 *
 */

package global

import (
	"fmt"
)
 
 

var SumConnections int32

var Qps int64

 

func CheckError(err error) {
	if err != nil {
		fmt.Println("Fatal error: %s", err.Error())
	}
}
