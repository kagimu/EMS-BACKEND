package responses

type DecodedAdmissionLetter struct {
	AdmissionLetterNumber string `json:"admission_letter_number" binding:"required"`
	Phrase                string `json:"phrase" binding:"required"`
}
