package api

import (
	"github.com/gin-gonic/gin"
	"github.com/optimism-java/interopbankend/internal/schema"
	config "github.com/optimism-java/interopbankend/internal/types"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Api struct {
	Config *config.Config
	DB     *gorm.DB
}

func NewAPIHandler(config *config.Config, db *gorm.DB) *Api {
	return &Api{Config: config, DB: db}
}

func (h Api) GetSyncBlocks(c *gin.Context) {
	// Get query parameters from URL
	pageStr := c.Query("page")
	pageSizeStr := c.Query("pageSize")

	// Convert string to int with default values
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var blocks []schema.SyncBlock
	var total int64

	// Replace with your actual DB logic
	h.DB.Model(&schema.SyncBlock{}).Count(&total)
	h.DB.Offset(offset).Limit(pageSize).Order("block_number desc").Find(&blocks)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": gin.H{
			"list":     blocks,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

func (h Api) GetExecutingMessageByBlockNumber(c *gin.Context) {
	blockNumberStr := c.Param("blockNumber")
	blockNumber, err := strconv.ParseInt(blockNumberStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Invalid block number",
			"data": nil,
		})
		return
	}

	var events []schema.SyncEvent

	h.DB.Where("block_number = ?", blockNumber).Order("block_number desc").Find(&events)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": events,
	})
}

func (h Api) GetSentMessageByHash(c *gin.Context) {
	hash := c.Param("hash")

	var event schema.SyncEvent
	h.DB.Table("sync_events").
		Where("JSON_UNQUOTE((data -> '$.msgHash') = ?", hash).
		Find(&event)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": event,
	})
}

func (h Api) GetRelayedMessageByHash(c *gin.Context) {
	hash := c.Param("hash")

	var event schema.SyncEvent
	// Replace with your actual DB logic
	h.DB.Where("relayed_msg_hash = ?", hash).First(&event)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": event,
	})
}
