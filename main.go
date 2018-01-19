package main

import (
	"fmt"
	"path"

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
	picHash := c.Param("pic-hash")
	fmt.Println(picHash)
	styleID := c.Param("style-id")
	fmt.Println(styleID)
	artPath := path.Join(workDir, picHash, styleID, "art.jpg")
	fmt.Println(artPath)
	c.File(artPath)
}

func main() {
	InitTaskQueue()
	go taskQueue.tDownload()
	go taskQueue.tTransferArt()
	r := gin.Default()
	r.POST("/transfer-to-art/add-task", receiveTask)
	r.GET("/transfer-to-art/:pic-hash/:style-id", getFile)
	r.Run() // listen and serve on 0.0.0.0:8080
}
