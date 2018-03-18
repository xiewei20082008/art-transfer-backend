package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"os/exec"
	"log"
)

type Task struct {
	Taskid  string `form:"task_id" binding:"required"`
	PicURL  string `form:"pic_url" binding:"required"`
	PicHash string `form:"pic_hash" binding:"required"`
	Style   string `form:"style" binding:"required"`
}

// TaskQueue two task queue
type TaskQueue struct {
	downloadQ chan Task
	transferQ chan Task
}

var taskQueue *TaskQueue
var workDir = "/home/weix/work-dir/"

//InitTaskQueue before use
func InitTaskQueue() {
	taskQueue = &TaskQueue{downloadQ: make(chan Task, 100), transferQ: make(chan Task, 100)}
}

func copy(src string, dst string) {
	srcFile, _ := os.Open(src)
	defer srcFile.Close()
	destFile, _ := os.Create(dst) // creates if file doesn't exist
	defer destFile.Close()
	io.Copy(destFile, srcFile) // check first var for number of bytes copied
	destFile.Sync()
}

func transfer(src string, dst string, style string) {
	execution := "/home/ubuntu/src-code/fast-style-transfer/evaluate.py"
	checkpoint := path.Join("/home/ubuntu/src-code/style",style,"style.ckpt")
	cmd := exec.Command("/usr/bin/python3", execution,
		"--checkpoint", checkpoint,
		"--in-path", src,
		"--out-path", dst)
	cmd.Dir = "/home/ubuntu/src-code/fast-style-transfer"
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func downloadFile(picURL string, picHash string) (bool, error) {
	downloadDir := path.Join(workDir, picHash)
	downloadPath := path.Join(downloadDir, "src.jpg")
	fmt.Println("downloading file to ", downloadPath)
	err := os.MkdirAll(downloadDir, os.ModePerm)
	if err != nil {
		fmt.Println("cannot create dir")
		return false, err
	}
	if _, err := os.Stat(downloadPath); err == nil {
		fmt.Println("file already existed.")
		return true, nil
	}
	out, err := os.Create(downloadPath)
	if err != nil {
		fmt.Println("cannot create target file")
		return false, err
	}
	defer out.Close()
	resp, err := http.Get(picURL)
	if err != nil {
		fmt.Println("Failed to get file")
		return false, err
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Failed to copy file")
		return false, err
	}
	return true, nil
}

func (taskQ TaskQueue) addNewTask(taskid string, picURL string, picHash string, style string) {
	fmt.Printf("Received new task %v || %s || %s || %s\n", taskid, picURL, picHash, style)
	newTask := Task{taskid, picURL, picHash, style}
	taskQ.downloadQ <- newTask
}

func (taskQ TaskQueue) tDownload() {
	for {
		downloadTask := <-taskQ.downloadQ
		fmt.Printf("downloading %v\n", downloadTask.PicURL)
		_, err := downloadFile(downloadTask.PicURL, downloadTask.PicHash)
		if err != nil {
			fmt.Printf("Failed to download pic %v\n", downloadTask.PicURL)
			continue
			// taskQ.downloadQ <- downloadTask
		}
		taskQ.transferQ <- downloadTask
	}
}

func (taskQ TaskQueue) tTransferArt() {
	for {
		transferTask := <-taskQ.transferQ
		fmt.Printf("transferring %v with style %v\n", transferTask.PicURL, transferTask.Style)
		artPath := path.Join(workDir, transferTask.PicHash, transferTask.Style)
		os.MkdirAll(artPath, os.ModePerm)
		downloadDir := path.Join(workDir, transferTask.PicHash)
		downloadPath := path.Join(downloadDir, "src.jpg")
		transfer(downloadPath, path.Join(artPath, "art.jpg"), transferTask.Style)
		fmt.Printf("transferQ solved\n")
		url := fmt.Sprintf("http://art.not.com.cn/open/changeTaskStatus?task_id=%s", transferTask.Taskid)
		http.Get(url)
	}
}
