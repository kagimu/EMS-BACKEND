package requests

type LanguageItemRequest struct {
	Language string `json:"language"`
	Speaking string `json:"speaking"`
	Reading  string `json:"reading"`
	Writing  string `json:"writing"`
}

type StaffBioDataLanguageRequest struct {
	Languages []LanguageItemRequest `json:"languages"`
}

func (r *StaffBioDataLanguageRequest) Validate() error {
	return nil
}
