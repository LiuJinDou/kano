package dto

type GetUploadTokenRequest struct {
	ApplicationCode string `json:"application_code" binding:"required"`
}

type GetUploadTokenResponse struct {
	Credentials   Credentials `json:"Credentials"`
	ExpiredTime   int         `json:"ExpiredTime"`
	Expiration    string      `json:"Expiration"`
	StartTime     int         `json:"StartTime"`
	RequestId     string      `json:"RequestId"`
	BucketName    string      `json:"bucket_name"`
	RegionName    string      `json:"region_name"`
	Authorization string      `json:"authorization"` // Authorization token for the upload
}
type Credentials struct {
	TmpSecretId  string `json:"TmpSecretId"`
	TmpSecretKey string `json:"TmpSecretKey"`
	Token        string `json:"Token"`
}

type SaveUploadRecordRequest struct {
}
