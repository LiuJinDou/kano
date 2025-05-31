package tencent

import (
	"context"
	"fmt"
	"kano/api/dto"
	"kano/internal/config"
	"kano/internal/repository"
	"kano/internal/repository/model"
	"mime/multipart"
	"time"

	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

type TencentUploader struct {
}

func New() *TencentUploader {
	return &TencentUploader{}
}

// Upload 上传文件到腾讯云COS
func (u *TencentUploader) Upload(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	// 文件上传逻辑
	return "/path/file.png", nil
}

// GetCredential 获取腾讯云临时密钥
func (u *TencentUploader) GetCredential(ctx context.Context, application_code string) (interface{}, error) {
	fmt.Println(config.Config.TencentYun)
	c := sts.NewClient(config.Config.TencentYun.SecretID, config.Config.TencentYun.SecretKey, nil)
	opt := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          config.Config.TencentYun.Region,
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					Action: []string{
						// 简单上传
						"name/cos:PostObject",
						"name/cos:PutObject",
						// 分片上传
						"name/cos:InitiateMultipartUpload",
						"name/cos:ListMultipartUploads",
						"name/cos:ListParts",
						"name/cos:UploadPart",
						"name/cos:CompleteMultipartUpload",
					},
					Effect: "allow",
					Resource: []string{
						"qcs::cos:ap-shanghai:uid/1257262858:shadow-1257262858/tuling",
					},
					Condition: map[string]map[string]interface{}{
						"ip_equal": map[string]interface{}{
							"qcs:ip": []string{
								"*",
							},
						},
					},
				},
			},
		},
	}
	res, err := c.GetCredential(opt)
	if err != nil {
		return nil, err
	}
	return &dto.GetUploadTokenResponse{
		Credentials: dto.Credentials{
			TmpSecretId:  res.Credentials.TmpSecretID,
			TmpSecretKey: res.Credentials.TmpSecretKey,
			Token:        res.Credentials.SessionToken,
		},
		ExpiredTime: res.ExpiredTime,
		Expiration:  res.Expiration,
		StartTime:   res.StartTime,
		RequestId:   res.RequestId,
		BucketName:  config.Config.TencentYun.Bucket,
		RegionName:  config.Config.TencentYun.Region,
	}, nil

}

// UploadRecord 保存上传记录
func (u *TencentUploader) UploadRecord(ctx context.Context, uploadRecord []*model.UploadRecord) error {
	return repository.NewUploadRepo().SaveUploadRecord(ctx, uploadRecord)
}
