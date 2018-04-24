package main

import (
	"time"

	"github.com/xfstart07/gocrontab"
	log "github.com/xfstart07/logger"
)

func main() {
	scheduler := gocrontab.NewSchedule()
	scheduler.NewJob("test1").Every(10).Seconds().Do(func() {
		log.Println("test1......")
	})

	scheduler.Start()

	scheduler.NewJob("test2").Every(15).Seconds().Do(func() {
		log.Println("test2.....")
	})

	scheduler.NewJob("taskWithParams").Every(10).Seconds().Do(taskWithParams, 1, "hello")

	log.Println(scheduler.Entries())
	for {
		time.Sleep(15 * time.Second)

		scheduler.Stop()

		break
	}

	log.Println("删除", scheduler.RemoveJob("test2"))
	log.Println("删除", scheduler.RemoveJob("test2"))

	scheduler.NewJob("testc").Every(15).Minutes().Do(func() {
		log.Println("test2.....")
	})

	scheduler.NewJob("testb").Every(1).Days().At(12, 0).Do(func() {
		log.Println("testb...")
	})

	entries := scheduler.Entries()
	log.Println(scheduler.Entries())
	for idx := range entries {
		log.Printf("%+v\n", entries[idx])
	}
	log.Println(scheduler.Entries())
}

func taskWithParams(a int, b string) {
	log.Println(a, b)
}
