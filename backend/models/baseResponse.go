package models

type Response interface {
	IsSuccess() bool
	GetMessage() string
}

type BaseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type HelloResponse struct {
	BaseResponse
	Data string `json:"data"`
}
func NewHelloResponse(success bool, message string) *HelloResponse {
	return &HelloResponse{
		BaseResponse: BaseResponse{
			Success: success,
			Message: message,
		},
		Data: "Hello, World!",
	}
}

func (r BaseResponse) IsSuccess() bool {
	return r.Success
}

func (r BaseResponse) GetMessage() string {
	return r.Message
}

type BaseStegoRequest struct {
	Key            string `json:"key" binding:"required"`
	UseEncryption  bool   `json:"use_encryption"`
	UseRandomStart bool   `json:"use_random_start"`
	LSBBits        int    `json:"lsb_bits" binding:"required,min=1,max=4"`
}

type StegoRequest struct {
	BaseStegoRequest
	SecretFilename string `json:"secret_filename"`
}

type StegoResponse struct {
	BaseResponse
	PSNR         float64 `json:"psnr,omitempty"`
	StegoFileURL string  `json:"stego_file_url,omitempty"`
}

func NewStegoResponse(success bool, message string, psnr float64, url string) *StegoResponse {
	return &StegoResponse{
		BaseResponse: BaseResponse{
			Success: success,
			Message: message,
		},
		PSNR:         psnr,
		StegoFileURL: url,
	}
}

type ExtractRequest struct {
	BaseStegoRequest
}

type ExtractResponse struct {
	BaseResponse
	SecretFileURL  string `json:"secret_file_url,omitempty"`
	SecretFilename string `json:"secret_filename,omitempty"`
}

func NewExtractResponse(success bool, message, fileURL, filename string) *ExtractResponse {
	return &ExtractResponse{
		BaseResponse: BaseResponse{
			Success: success,
			Message: message,
		},
		SecretFileURL:  fileURL,
		SecretFilename: filename,
	}
}
