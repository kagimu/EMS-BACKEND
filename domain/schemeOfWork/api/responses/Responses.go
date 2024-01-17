package responses

import (
	"errors"
	"github.com/signintech/gopdf"
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

func (schemeData *SchemeOfWorkResponseInfo) ReturnSchemeOfWorkPdf() (string, error) {

	pdf := gopdf.GoPdf{}

	htmlContent := `<h1>Hello, PDF!</h1>
	<p>This is a sample PDF generated from HTML-like content.</p>`

	err, path := generatePDFWithDefaultFont(&pdf, htmlContent)
	if err != nil {
		return "", errors.New(err.Error())
	}

	return path, nil

	/*// Your HTML-like content
	htmlContent := `<h1>Hello, PDF!</h1>
	<p>This is a sample PDF generated from HTML-like content.</p>`

	err := generatePDFFromHTML(&pdf, htmlContent)
	if err != nil {
		return "", errors.New(err.Error())
	}

	path := "/fileRepo/schemeTemplatePdf/" + schemeData.Name + ".pdf"
	// Save the PDF to a file
	err = pdf.WritePdf(path)
	if err != nil {
		return "", errors.New("----2. " + err.Error() + "----")
	}

	filePath := "https://" + context.Request.Host + path

	return filePath, nil*/

}

///*func generatePDFFromHTML(pdf *gopdf.GoPdf, htmlContent string) error {
//	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
//	pdf.AddPage()*/

// Use a default system font (e.g., "Arial")
/*err := pdf.AddTTFFont("Arial", "")
if err != nil {
	return errors.New("---3. " + err.Error() + "----")
}*/
/*err := pdf.SetFont("Arial", "", 14)
if err != nil {
	return errors.New("--4. " + err.Error() + "----")
}*/

// Simulate simple HTML rendering
/*	pdf.SetX(20)
	pdf.SetY(20)
	pdf.SetTextColor(0, 0, 0)

	err = pdf.Cell(nil, htmlContent)
	if err != nil {
		return errors.New("----4. " + err.Error() + "----")
	}

	return nil*/
//}

func generatePDFWithDefaultFont(pdf *gopdf.GoPdf, htmlContent string) (error, string) {
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	// No need to add a custom font, it will use the default font

	// Write HTML-like content to the PDF
	pdf.SetY(20)
	pdf.SetTextColor(0, 0, 0)
	err := pdf.Cell(nil, htmlContent)
	if err != nil {
		return err, ""
	}

	// Specify the folder path and file name
	filePath := "/fileRepo/schemeTemplatePdf/output.pdf.pdf"

	// Save the PDF to the specified file path
	err = pdf.WritePdf(filePath)
	if err != nil {
		return err, ""
	}

	return nil, filePath
}
