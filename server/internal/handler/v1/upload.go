package v1

import (
	"kano/internal/config"
	"kano/internal/repository/model"
	"kano/internal/service"
	"kano/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetUploadToken 获取上传凭证
// @Summary      获取上传凭证
// @Description  根据应用码获取上传凭证
// @Tags         上传服务
// @Accept       json
// @Produce      json
// @Param        application_code query string true "应用码"
// @Success      200 {object} dto.Response{data=interface{}}
// @Failure      2001 {object} dto.Response
// @Router       /v1/upload/token [get]
func GetUploadToken(c *gin.Context) {
	// Handle file upload
	application_code, res := c.GetQuery("application_code")
	if !res || application_code == "" {
		response.Error(c, response.CodeParameter, map[string]interface{}{
			"error": "application_code is required"})
		return
	}

	// Validate the application_code
	if !config.IsApplicationCodeAllowed(application_code) {
		response.Error(c, response.CodeInvalidApplicationCode)
		return
	}

	// Save the file to a specific location
	uploader, err := service.NewUploader(service.UploadTypeTencent).GetCredential(c, application_code)

	if err != nil {
		response.Error(c, response.CodeDefault, err.Error())
		return
	}
	response.Success(c, uploader)
}

// SaveUploadRecord 保存上传记录
// @Summary      保存上传记录
// @Description  将上传完成后的文件记录保存到数据库
// @Tags         上传服务
// @Accept       json
// @Produce      json
// @Param        records body []model.UploadRecord true "上传记录数组"
// @Success      200 {object} dto.Response
// @Failure      2001 {object} dto.Response
// @Router       /v1/upload/record [post]
func SaveUploadRecord(c *gin.Context) {
	// Handle saving upload record
	var pararms []*model.UploadRecord
	if err := c.ShouldBindJSON(&pararms); err != nil {
		response.Error(c, response.CodeParameter, err.Error())
		return
	}

	// Save the file to a specific location
	err := service.NewUploader(service.UploadTypeTencent).UploadRecord(c, pararms)

	if err != nil {
		response.Error(c, response.CodeDefault, err.Error())
		return
	}
	response.Success(c, nil)
}
