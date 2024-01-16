package responses

import (
	"errors"
)

type (
	SchemeOfWorkSubtopicInfo struct {
		Number             string `json:"number"`
		Code               int64  `json:"code"`
		SortCode           int64  `json:"sort_code"`
		Name               string `json:"name"`
		TopicName          string `json:"topic_name"`
		StrandName         string `json:"strand_name"`
		EstimatedTime      int64  `json:"estimated_time"`
		SchemeOfWorkNumber string `json:"scheme_of_work_number"`
		SubTopicNumber     string `json:"sub_topic_number"`
		UserNumber         string `json:"user_number"`
		ConceptCount       int64  `json:"concept_count"`
	}

	SchemeOfWorkResponseInfo struct {
		Number                   string                     `json:"number"`
		Name                     string                     `json:"name"`
		Status                   string                     `json:"status"`
		ClassNumber              string                     `json:"class_number"`
		ClassName                string                     `json:"class_name"`
		SubjectNumber            string                     `json:"subject_number"`
		SubjectName              string                     `json:"subject_name"`
		UserNumber               string                     `json:"user_number"`
		AcademicYearNumber       string                     `json:"academic_year_number"`
		AcademicYearName         string                     `json:"academic_year_name"`
		SchemeOfWorkSubtopicInfo []SchemeOfWorkSubtopicInfo `json:"subtopic"`
	}
)

func (userData *SchemeOfWorkResponseInfo) ReturnSchemeOfWorkPdf() (string, error) {
	return "", errors.New("")
	/*// Read HTML template file
	htmlTemplate, err := os.ReadFile("template.html")
	if err != nil {
		return "", errors.New(err.Error())
	}

	// Parse HTML template
	tmpl, err := template.New("pdfTemplate").Parse(string(htmlTemplate))
	if err != nil {
		return "", errors.New(err.Error())
	}

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Execute HTML template and write content to PDF
	err = tmpl.Execute(pdf, userData)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Error executing HTML template", http.StatusInternalServerError)
		return
	}

	// Specify the directory to save the PDF
	outputDirectory := "./pdf_output"
	outputPath := fmt.Sprintf("%s/user_data.pdf", outputDirectory)

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDirectory, os.ModePerm); err != nil {
		log.Fatal(err)
		http.Error(w, "Error creating output directory", http.StatusInternalServerError)
		return
	}

	// Save PDF to the specified file
	if err := pdf.OutputFileAndClose(outputPath); err != nil {
		log.Fatal(err)
		http.Error(w, "Error saving PDF file", http.StatusInternalServerError)
		return
	}

	fmt.Printf("PDF saved to: %s\n", outputPath)
	w.Write([]byte("PDF saved successfully"))
	*/
}
