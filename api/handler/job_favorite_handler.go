package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"org.thinkinai.com/recruit-center/api/dto/response"
	"org.thinkinai.com/recruit-center/internal/service"
	"org.thinkinai.com/recruit-center/pkg/errors"
)

type JobFavoriteHandler struct {
	favoriteService *service.JobFavoriteService
}

func NewJobFavoriteHandler(service *service.JobFavoriteService) *JobFavoriteHandler {
	return &JobFavoriteHandler{favoriteService: service}
}

// AddFavorite 收藏职位
// @Summary 收藏职位
// @Description 添加职位到收藏夹
// @Tags 职位收藏
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param jobId path int true "职位ID"
// @Success 0000 {object} response.Response
// @Router /api/v1/jobs/{jobId}/favorite [post]
func (h *JobFavoriteHandler) AddFavorite(c *gin.Context) {
	jobID, _ := strconv.ParseUint(c.Param("jobId"), 10, 32)
	userID := c.GetUint("userId")

	if err := h.favoriteService.AddFavorite(userID, uint(jobID)); err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(nil))
}

// RemoveFavorite 取消收藏
// @Summary 取消收藏
// @Description 取消收藏指定职位
// @Tags 职位收藏
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param jobId path int true "职位ID"
// @Success 0000 {object} response.Response
// @Router /api/v1/jobs/{jobId}/favorite [delete]
func (h *JobFavoriteHandler) RemoveFavorite(c *gin.Context) {
	jobID, _ := strconv.ParseUint(c.Param("jobId"), 10, 32)
	userID := c.GetUint("userId")

	if err := h.favoriteService.RemoveFavorite(userID, uint(jobID)); err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(nil))
}

// ListFavorites 获取收藏列表
// @Summary 获取收藏列表
// @Description 分页获取用户收藏的职位列表
// @Tags 职位收藏
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int false "页码" minimum(1) default(1)
// @Param size query int false "每页数量" minimum(1) maximum(100) default(10)
// @Success 0000 {object} response.Response{data=response.JobListResponse}
// @Router /api/v1/users/favorites [get]
func (h *JobFavoriteHandler) ListFavorites(c *gin.Context) {
	userID := c.GetUint("userId")
	page, size := parsePageSize(c)

	favorites, err := h.favoriteService.ListFavorites(userID, page, size)
	if err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(favorites))
}

// GetUserStatistics 获取用户收藏统计
// @Summary 获取用户收藏统计
// @Description 统计用户收藏的职位信息
// @Tags 职位收藏
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {object} response.Response{data=response.JobFavoriteStatistics}
// @Router /api/v1/users/favorites/statistics [get]
func (h *JobFavoriteHandler) GetUserStatistics(c *gin.Context) {
	userID := c.GetUint("userId")

	stats, err := h.favoriteService.GetUserStatistics(userID)
	if err != nil {
		c.JSON(http.StatusOK, errors.Wrap(err, errors.InternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(stats))
}
