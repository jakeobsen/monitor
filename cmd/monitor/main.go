package main

import "monitor"

func main() {
	d := monitor.NewDaemon()
	d.Run()
}
