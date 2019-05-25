package main

import (
	"fmt"
	"github.com/kr/pretty"
	"io"
	"log"
	"net"
	"runtime/debug"
	"time"
)

var gLocalConn net.Listener
var sourcePort  = "0.0.0.0:10000"
var targetPort  =  "0.0.0.0:3213"


func forever(fn func()) {
	f := func() {
		defer func() {
			if r := recover(); r != nil {
				debug.PrintStack()
				pretty.Println("Recover from error:", r)
			}
		}()
		fn()
	}
	for {
		f()
	}
}
func handle(){
	for {
		sourceConn, err := gLocalConn.Accept()

		if err != nil {
			log.Println("server err:", err.Error())
		}
		targetConn, err := net.DialTimeout("tcp", targetPort, 30*time.Second)

		go func() {
			_, err = io.Copy(targetConn, sourceConn)
			if err != nil {
				fmt.Println("io.Copy 1 failed：", err.Error())
			}
		}()

		go func() {
			_, err = io.Copy(sourceConn, targetConn)
			if err != nil {
				fmt.Println("io.Copy 2 failed：", err.Error())
			}
		}()

	}
}
func main(){
	fmt.Println("sourcePort：", sourcePort, "targetPort：", targetPort)
	localConn, err := net.Listen("tcp", sourcePort)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	gLocalConn = localConn
	forever(handle)
}

