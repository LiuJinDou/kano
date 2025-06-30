package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

const TableNameUploadRecord = "upload_records"

// UploadRecord 上传记录表（包含成功与失败）
type UploadRecord struct {
	ID             int64                 `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:唯一标识每条上传记录，自增" json:"id"`                               // 唯一标识每条上传记录，自增
	UploadID       string                `gorm:"column:upload_id;type:varchar(64);not null;comment:上传任务ID" json:"upload_id"`                                        // 上传任务ID
	Status         int32                 `gorm:"column:status;type:tinyint(1);not null;comment:0success, 1failure 表示该次上传是成功还是失败，用于判断是否需要触发告警或后续处理逻辑" json:"status"` // 0success, 1failure 表示该次上传是成功还是失败，用于判断是否需要触发告警或后续处理逻辑
	FilePath       string                `gorm:"column:file_path;type:varchar(255);comment:文件在服务器上的存储路径，仅在上传成功时填写，用于后续访问文件或生成下载链接" json:"file_path"`                // 文件在服务器上的存储路径，仅在上传成功时填写，用于后续访问文件或生成下载链接
	ErrorMessage   string                `gorm:"column:error_message;type:varchar(255);comment:上传失败时的错误信息，用于定位问题原因，如“文件过大”、“格式不支持”、“网络中断”等" json:"error_message"`   // 上传失败时的错误信息，用于定位问题原因，如“文件过大”、“格式不支持”、“网络中断”等
	CredentialType bool                  `gorm:"column:credential_type;type:tinyint(1);comment:凭证类型（1普通凭证、2加密凭证）" json:"credential_type"`                           // 凭证类型（1普通凭证、2加密凭证）
	BucketSpec     string                `gorm:"column:bucket_spec;type:varchar(255);comment:bucket 规范" json:"bucket_spec"`                                         // bucket 规范
	BusinessType   int32                 `gorm:"column:business_type;type:int;comment:用于区分不同模块的上传操作（如头像上传、素材上传等）" json:"business_type"`                             // 用于区分不同模块的上传操作（如头像上传、素材上传等）
	Username       string                `gorm:"column:username;type:varchar(0);comment:用户名字" json:"username"`                                                      // 用户名字
	CreatedAt      time.Time             `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;comment:记录创建时间，默认当前时间" json:"created_at"`                 // 记录创建时间，默认当前时间
	UpdatedAt      time.Time             `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updated_at"`                                                    // 更新时间
	DeletedAt      soft_delete.DeletedAt `gorm:"softDelete:flag;column:deleted_at;type:int;comment:删除时间;" json:"deleted_at" swaggerignore:"true"`                   // 删除时间
}

// TableName UploadRecord's table name
func (*UploadRecord) TableName() string {
	return TableNameUploadRecord
}
