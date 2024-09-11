package responses

type AttachmentResponse struct {
	FilePath  string `json:"file_path"`
	FileName  string `json:"file_name"`
	Extension string `json:"extension"`
}
