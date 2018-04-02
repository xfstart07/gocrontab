package main

import (
	"gocrontab"
	"fmt"
	"time"
)

func main() {
	scheduler := gocrontab.NewSchedule()
	scheduler.NewJob("test1").Every(10).Seconds().Do(func() {
		fmt.Println("test1......")
	})

	scheduler.Start()

	scheduler.NewJob("test2").Every(15).Seconds().Do(func() {
		fmt.Println("test2.....")
	})

	fmt.Println(scheduler.Entries())
	for {
		time.Sleep(15 * time.Second)

		scheduler.Stop()

		break
	}

	fmt.Println("删除", scheduler.RemoveJob("test2"))
	fmt.Println("删除", scheduler.RemoveJob("test2"))

	scheduler.NewJob("testc").Every(15).Minutes().Do(func() {
		fmt.Println("test2.....")
	})

	scheduler.NewJob("testb").Every(1).Days().At(12, 0).Do(func() {
		fmt.Println("testb...")
	})

	entries := scheduler.Entries()
	for idx := range entries {
		fmt.Printf("%+v\n", entries[idx])
	}
	fmt.Println(scheduler.Entries())
}

