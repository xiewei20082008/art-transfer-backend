package main

import (
	"github.com/gin-gonic/gin"
)

func receiveTask(c *gin.Context) {
	var task Task
	c.Bind(&task)
	taskQueue.addNewTask(task.Taskid, task.PicURL, task.PicHash, task.Style)
	c.JSON(400, gin.H{
		"message": "Add Task OK",
	})
}

func getFile(c *gin.Context) {
	c.File("/home/weix/work-dir/md5/src.jpg")
}

func main() {
	InitTaskQueue()
	go taskQueue.tDownload()
	go taskQueue.tTransferArt()
	r := gin.Default()
	r.POST("/transfer-to-art/task", receiveTask)
	r.GET("/transfer-to-art/pic-hash/style", getFile)
	r.Run() // listen and serve on 0.0.0.0:8080
}
