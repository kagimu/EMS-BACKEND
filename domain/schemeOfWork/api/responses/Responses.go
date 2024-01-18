package responses

import (
	"errors"
	"github.com/gin-gonic/gin"
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

func (schemeData *SchemeOfWorkResponseInfo) ReturnSchemeOfWorkPdf(context *gin.Context) (string, error) {

	pdf := gopdf.GoPdf{}

	// Your HTML-like content

	/*	htmlContent := "<h1>Hello, PDF!</h1>"
		htmlContent += "<p>This is a sample PDF generated with HTML-like content.</p>"
	*/

	htmlContent := `
<html>
  <head>
    <title>My Mega Service</title>
    <link rel=\"stylesheet\" href=\"main.css\">
    <style type=\"text/css\">body { color: #fff; }</style>
  </head>

  <body>
    <div class="logo">
      <a href="http://jaytaylor.com/"><img src="/logo-image.jpg" alt="Mega Service"/></a>
    </div>

    <h1>Welcome to your new account on my service!</h1>

    <p>
      Here is some more information:

      <ul>
        <li>Link 1: <a href="https://example.com">Example.com</a></li>
        <li>Link 2: <a href="https://example2.com">Example2.com</a></li>
        <li>Something else</li>
      </ul>
    </p>

    <table>
      <thead>
        <tr><th>Header 1</th><th>Header 2</th></tr>
      </thead>
      <tfoot>
        <tr><td>Footer 1</td><td>Footer 2</td></tr>
      </tfoot>
      <tbody>
        <tr><td>Row 1 Col 1</td><td>Row 1 Col 2</td></tr>
        <tr><td>Row 2 Col 1</td><td>Row 2 Col 2</td></tr>
      </tbody>
    </table>
  </body>
</html>`

	/*htmlText, err := html2text.FromString(htmlContent, html2text.Options{PrettyTables: true})
	if err != nil {
		panic(err)
	}*/

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

	return filePath, nil

}

func generatePDFFromHTML(pdf *gopdf.GoPdf, htmlContent string) error {
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	err := pdf.AddTTFFont("Arial", "/fileRepo/ttf/arial.ttf")
	if err != nil {
		return errors.New("---3. " + err.Error() + "----")
	}
	err = pdf.SetFont("Arial", "", 12)
	if err != nil {
		return errors.New("--4. " + err.Error() + "----")
	}

	// Simulate simple HTML rendering
	pdf.SetX(20)
	pdf.SetY(20)
	pdf.SetTextColor(0, 0, 0)

	err = pdf.Cell(nil, htmlContent)
	if err != nil {
		return errors.New("----4. " + err.Error() + "----")
	}

	return nil
}
