package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hliangzhao/LearnGo/11-gin/models"
	"sync"
)

// VideoController 给外界开放的是名为VideoController的接口，内部返回具体实现：&controller{}
type VideoController interface {
	GetAll(ctx *gin.Context)
	Update(ctx *gin.Context)
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

func NewVideoController() VideoController {
	return &controller{videos: make([]models.Video, 0)}
}

type generator struct {
	counter int
	mtx     sync.Mutex // 加锁是因为考虑到了web站点的并发访问性
}

func (g *generator) getNextId() int {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	g.counter++
	return g.counter
}

// 创建全局的id生成器
var g = &generator{}

type controller struct {
	videos []models.Video
}

func (c *controller) GetAll(ctx *gin.Context) {
	// 使用gin提供的上下文在前后端之间传递数据
	ctx.JSON(200, c.videos)
}

func (c *controller) Update(ctx *gin.Context) {
	var videoForUpdate models.Video
	// 这里要更新的video的id是写进url里面的，id扮演的角色是uri
	if err := ctx.ShouldBindUri(&videoForUpdate); err != nil {
		ctx.String(400, "bad request %v", err)
		return
	}
	// 要更新的video的新信息来自用户的json输入，因此也要bind！
	if err := ctx.ShouldBindJSON(&videoForUpdate); err != nil {
		ctx.String(400, "bad request %v", err)
		return
	}

	for idx, video := range c.videos {
		if video.Id == videoForUpdate.Id {
			c.videos[idx] = videoForUpdate
			ctx.String(200, "success, video with id %d has been updated", videoForUpdate.Id)
			return
		}
	}
	ctx.String(400, "bad request, cannot find video with %d to update", videoForUpdate.Id)
}

func (c *controller) Create(ctx *gin.Context) {
	video := models.Video{Id: g.getNextId()}
	// models.Video的Title和Description都被tag为json
	if err := ctx.BindJSON(&video); err != nil {
		ctx.String(400, "bad request %v", err)
	}
	c.videos = append(c.videos, video)
	ctx.String(200, "success, new video id is %d", video.Id)
}

func (c *controller) Delete(ctx *gin.Context) {
	var videoToDelete models.Video
	// 这里要更新的video的id是写进url里面的，id扮演的角色是uri
	if err := ctx.ShouldBindUri(&videoToDelete); err != nil {
		ctx.String(400, "bad request %v", err)
		return
	}
	for idx, video := range c.videos {
		if video.Id == videoToDelete.Id {
			// slice的操作并不是thread safe的，并发场景时要进行同步
			c.videos = append(c.videos[0:idx], c.videos[idx+1:len(c.videos)]...)
			ctx.String(200, "success, video with id %d has been deleted", videoToDelete.Id)
			return
		}
	}
	ctx.String(400, "bad request, cannot delete video with %d to update", videoToDelete.Id)
}
