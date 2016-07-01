package bftraderclient

import "os"
import "sync/atomic"
import "log"
import "os/signal"

// ctrl+c monitor:
// http://studygolang.com/articles/2333
// http://sugarmanman.blog.163.com/blog/static/8107908020136713147504/
func monCtrlc(exitNow *int32) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	s := <-c
	log.Printf("Got signal: %v, exiting......", s)
	signal.Stop(c)
	atomic.StoreInt32(exitNow, 1)
}
