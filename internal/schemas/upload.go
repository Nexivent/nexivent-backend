package schemas

// PresignUploadRequest represents the payload to request a presigned upload URL.
type PresignUploadRequest struct {
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
}

// PresignUploadResponse is returned to the front with data to upload directly to S3.
type PresignUploadResponse struct {
	UploadURL string `json:"uploadUrl"`
	Key       string `json:"key"`
}
