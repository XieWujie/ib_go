package main

import (
	_ "../db"
	"../imio"
	"log"
	"sync"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func main()  {
	wait :=sync.WaitGroup{}
	wait.Add(1)
	imio.RegisterHttpListener()
	wait.Wait()
}
