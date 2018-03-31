## 功能

[x] 创建任务

[x] 删除任务

[x] 设定任务运行时间，间隔

[x] 友好的设定时间方法，包括：间隔，秒，分，时，日

## 使用

```
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

entries := scheduler.Entries()
for idx := range entries {
	fmt.Println(entries[idx].GetName())
}
fmt.Println(scheduler.Entries())
```