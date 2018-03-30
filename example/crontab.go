package example

import (
	"gocrontab"
	"fmt"
	"time"
)

func main() {

	scheduler := gocrontab.NewSchedule()
	scheduler.AddJob(10, "test1", printTest1)

	scheduler.Start()

	scheduler.AddJob(15, "test2", printTest2)

	fmt.Println(scheduler.Entries())
	for {
		time.Sleep(15 * time.Second)

		scheduler.Stop()

		break
	}

	fmt.Println("删除", scheduler.RemoveJob("test2"))

	scheduler.AddJob(15, "testc", printTest2)

	fmt.Println(scheduler.Entries())
}

func printTest1() {
	fmt.Println("test1......")
}

func printTest2() {
	fmt.Println("test2.....")
	time.Sleep(10 * time.Second)
}
