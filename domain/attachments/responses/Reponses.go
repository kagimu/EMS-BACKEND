package responses

type AttachmentResponse struct {
	Number    string `json:"number"`
	Title     string `json:"title"`
	Category  string `json:"category"`
	CreatedAt string `json:"created_at"`
	FilePath  string `json:"file_path"`
	FileName  string `json:"file_name"`
	Extension string `json:"extension"`
	URL       string `json:"url"`
}
