package tencent

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"kano/dto"
	"kano/internal/config"
	"kano/internal/repository"
	"kano/internal/repository/model"
	"mime/multipart"
	"net/url"
	"sort"
	"strings"
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
						"qcs::cos:" + config.Config.TencentYun.Region + ":uid/" + config.Config.TencentYun.AppId + ":" + config.Config.TencentYun.Bucket + "/application/" + application_code + "/*",
					},
					Condition: map[string]map[string]interface{}{},
				},
			},
		},
	}
	fmt.Println("qcs::cos:" + config.Config.TencentYun.Region + ":uid/" + config.Config.TencentYun.AppId + ":" + config.Config.TencentYun.Bucket + "/application/" + application_code + "/*")
	res, err := c.GetCredential(opt)
	if err != nil {
		return nil, err
	}

	// secretKey := "YourSecretKey"
	// httpMethod := "put"
	// uri := "/exampleobject"
	// params := map[string]string{} // GET/PUT参数
	// headers := map[string]string{
	// 	"content-length":   "13",
	// 	"content-md5":      "mQ/fVh815F3k6TAUm8m0eg==",
	// 	"content-type":     "text/plain",
	// 	"date":             "Thu, 16 May 2019 06:45:51 GMT",
	// 	"host":             "examplebucket-1250000000.cos.ap-beijing.myqcloud.com",
	// 	"x-cos-acl":        "private",
	// 	"x-cos-grant-read": "uin=\"100000000011\"",
	// }
	// startTime := int64(1557989151)
	// endTime := int64(1557996351)
	// signature, signKey, keyTime, _ := GenCosSignature(secretKey, httpMethod, uri, params, headers, startTime, endTime)
	// fmt.Println(signKey, keyTime, signature)
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
		// Authorization: fmt.Sprintf("q-sign-algorithm=sha1&q-ak=%s&q-sign-time=%d;%d&q-key-time=%d;%d&q-header-list=&q-url-param-list=&q-signature=%s", res.Credentials.TmpSecretKey, keyTime, res.ExpiredTime, res.StartTime, res.ExpiredTime, signature),
	}, nil

}

// UploadRecord 保存上传记录
func (u *TencentUploader) UploadRecord(ctx context.Context, uploadRecord []*model.UploadRecord) error {
	return repository.NewUploadRepo().SaveUploadRecord(ctx, uploadRecord)
}

// GenCosSignature 生成腾讯云COS标准签名
func GenCosSignature(secretKey, httpMethod, uri string, params, headers map[string]string, startTime, endTime int64) (signature, signKey, keyTime string, err error) {
	// 1. 计算 KeyTime
	keyTime = fmt.Sprintf("%d;%d", startTime, endTime)
	// 2. 计算 SignKey
	h := hmac.New(sha1.New, []byte(secretKey))
	h.Write([]byte(keyTime))
	signKey = fmt.Sprintf("%x", h.Sum(nil))

	// 3. 构造 HttpString
	// 参数、头部按 key 升序
	paramList, paramStr := buildSortedKV(params)
	headList, headStr := buildSortedKV(headers)
	httpString := fmt.Sprintf("%s\n%s\n%s\n%s\n", strings.ToLower(httpMethod), uri, paramStr, headStr)
	fmt.Println(paramList)
	fmt.Println(headList)
	// 4. 计算 StringToSign
	sha1HttpString := sha1Hex(httpString)
	stringToSign := fmt.Sprintf("sha1\n%s\n%s\n", keyTime, sha1HttpString)

	// 5. 计算 Signature
	h2 := hmac.New(sha1.New, []byte(signKey))
	h2.Write([]byte(stringToSign))
	signature = fmt.Sprintf("%x", h2.Sum(nil))
	return signature, signKey, keyTime, nil
}

// buildSortedKV 按 key 升序拼接，返回 key1;key2, key1=val1&key2=val2
func buildSortedKV(m map[string]string) (string, string) {
	if len(m) == 0 {
		return "", ""
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, strings.ToLower(k))
	}
	sort.Strings(keys)
	list := strings.Join(keys, ";")
	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		v := m[k]
		pairs = append(pairs, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(v)))
	}
	return list, strings.Join(pairs, "&")
}

// sha1Hex 计算字符串的sha1并返回16进制
func sha1Hex(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}
