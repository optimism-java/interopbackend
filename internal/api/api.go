package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/optimism-java/interopbackend/internal/schema"
	config "github.com/optimism-java/interopbackend/internal/types"
	"gorm.io/gorm"
)

type Api struct {
	Config *config.Config
	DB     *gorm.DB
}

func NewAPIHandler(config *config.Config, db *gorm.DB) *Api {
	return &Api{Config: config, DB: db}
}

// @Summary Get sync blocks with pagination
// @Description Get a paginated list of sync blocks
// @Tags Blocks
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} map[string]interface{} "code": 200, "msg": "success", "data": {"list": []schema.SyncBlock, "total": int64, "page": int, "pageSize": int}
// @Router /blocks [get]
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

// @Summary Get executing messages by block number
// @Description Get all executing messages for a specific block number
// @Tags Messages
// @Accept json
// @Produce json
// @Param blockNumber path int true "Block Number"
// @Success 200 {object} map[string]interface{} "code": 200, "msg": "success", "data": []schema.SyncEvent
// @Failure 400 {object} map[string]interface{} "code": 400, "msg": "Invalid block number", "data": nil
// @Router /blocks/{blockNumber}/executingMessage [get]
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

// @Summary Get sent message by hash
// @Description Get sent message details using message hash
// @Tags Messages
// @Accept json
// @Produce json
// @Param hash path string true "Message Hash"
// @Success 200 {object} map[string]interface{} "code": 200, "msg": "success", "data": schema.SyncEvent
// @Router /blocks/sentMessage/{hash} [get]
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

// @Summary Get relayed message by hash
// @Description Get relayed message details using message hash
// @Tags Messages
// @Accept json
// @Produce json
// @Param hash path string true "Message Hash"
// @Success 200 {object} map[string]interface{} "code": 200, "msg": "success", "data": schema.SyncEvent
// @Router /blocks/relayedMessage/{hash} [get]
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
