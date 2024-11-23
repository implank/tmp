package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

type transferReq struct {
	Filename string
	Size     int64  // 文件大小 B
	Src      string // 源地址
	Dst      string // 目标地址
}

type transferTask struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	// 文件大小 B
	TransferredSize int64 `json:"transferredSize"`
	// 已传输大小 B
	Progress float64 `json:"progress"`
	// 浮点数
	Src string `json:"src"`
	// 源地址
	Dst string `json:"dst"`
	// 目标地址
	TransferSpeed int64 `json:"transferSpeed"`
	// 传输速度 B/s
	NetBandwidth int64 `json:"netBandwidth"`
	// 网络带宽 B/s
}

func transfer(transferTask *transferTask) {
	var refreshFreq = 5 // 刷新频率
	for {
		if transferTask.TransferredSize >= transferTask.Size {
			fmt.Print("传输完成", transferTask.Filename)
			break
		}
		transferTask.NetBandwidth = baseNetBandwidth
		r := 0.8 + 0.01*rand.Float64()*7
		transferTask.TransferSpeed = int64(r * float64(baseNetBandwidth))
		transferTask.TransferredSize += transferTask.TransferSpeed / int64(refreshFreq)
		transferTask.Progress = float64(transferTask.TransferredSize) / float64(transferTask.Size)
		time.Sleep(time.Second / time.Duration(refreshFreq))
	}
}

func emulateNetBandwidth() {
	var refreshFreq = 20 // 刷新频率
	for {
		r := rand.Float64()*0.2 - 0.1
		var maxNetBandwidth int64 = 1024 * 1024 * 12
		var minNetBandwidth int64 = 1024 * 1024 * 8
		baseNetBandwidth = max(minNetBandwidth, min(maxNetBandwidth, baseNetBandwidth+int64(r*1024*1024)))
		time.Sleep(time.Second / time.Duration(refreshFreq))
	}
}

var tasks []*transferTask = make([]*transferTask, 0)
var baseNetBandwidth int64 = 1024 * 1024 * 10 // 10MB/s

func main() {
	go emulateNetBandwidth()

	r := gin.Default()
	r.Use(Cors())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.StaticFS("/static", http.Dir("static"))

	r.GET("netBandwidth", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"netBandwidth": baseNetBandwidth,
		})
	})

	r.POST("transfer", func(c *gin.Context) {
		var req transferReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "参数错误",
			})
			return
		}
		fmt.Print(req)
		task := transferTask{
			Filename: req.Filename,
			Size:     req.Size,
			Src:      req.Src,
			Dst:      req.Dst,
		}
		tasks = append(tasks, &task)
		fmt.Print(task)
		fmt.Print(tasks)
		go transfer(&task)
		c.JSON(http.StatusOK, gin.H{
			"message": "任务已添加",
		})
	})

	r.GET("tasks", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"tasks": tasks,
		})
	})
	r.Run()
}
