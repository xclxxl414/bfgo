package main

import "fmt"

import _ "google.golang.org/grpc"
import _ "github.com/sunwangme/bfgo/api/bfcta"
import _ "github.com/sunwangme/bfgo/api/bfdatafeed"
import _ "github.com/sunwangme/bfgo/api/bfgateway"
import _ "github.com/sunwangme/bfgo/api/bfkv"

func main() {
	fmt.Println("bfgo")
}
