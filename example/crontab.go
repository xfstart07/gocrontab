package main

import (
	"gocrontab"
	"fmt"
	"time"
)

func main() {
	scheduler := gocrontab.NewSchedule()
	scheduler.NewJob("test1", printTest1).Every(10).Seconds().Go()

	scheduler.Start()

	scheduler.NewJob("test2", printTest2).Every(15).Seconds().Go()

	fmt.Println(scheduler.Entries())
	for {
		time.Sleep(15 * time.Second)

		scheduler.Stop()

		break
	}

	fmt.Println("删除", scheduler.RemoveJob("test2"))
	fmt.Println("删除", scheduler.RemoveJob("test2"))

	scheduler.NewJob("testc", printTest2).Every(15).Minutes().Go()

	entries := scheduler.Entries()
	for idx := range entries {
		fmt.Println(entries[idx].GetName())
	}
	fmt.Println(scheduler.Entries())
}

func printTest1() {
	fmt.Println("test1......")
}

func printTest2() {
	fmt.Println("test2.....")
	time.Sleep(10 * time.Second)
}
