package main

import (
	"github.com/gin-gonic/gin"
)

func receiveTask(c *gin.Context) {
	var task Task
	c.BindJSON(&task)
	taskQueue.addNewTask(task.Taskid, task.PicURL, task.PicHash, task.Style)
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func main() {
	InitTaskQueue()
	go taskQueue.tDownload()
	go taskQueue.tTransferArt()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/transfer-to-art/task", receiveTask)
	r.Run() // listen and serve on 0.0.0.0:8080
}
