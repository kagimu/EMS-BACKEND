package responses

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/speedata/bagme/document"
	"github.com/speedata/boxesandglue/backend/bag"
	"log"
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

type Report struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (schemeData *SchemeOfWorkResponseInfo) ReturnSchemeOfWorkPdf(context *gin.Context) (string, error) {

	path, err := doThings(schemeData)
	if err != nil {
		log.Fatal(err)
	}

	filePath := "https://" + context.Request.Host + path

	return filePath, nil

	// Create a new PDF document
	//pdf := gopdf.GoPdf{}
	//pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	//pdf.AddPage()
	//
	//// Set the HTML content for the PDF
	//htmlContent := `<html><body><h1 style="color: blue;">Hello, this is a PDF with HTML and CSS content</h1></body></html>`
	///*err := pdf.SetHtml(htmlContent)
	//if err != nil {
	//	return "", errors.New(err)
	//}*/
	//
	//err := pdf.Cell(nil, htmlContent)
	//if err != nil {
	//	return "", errors.New(err.Error())
	//}
	//
	//// Save the PDF to a file
	//path := "/fileRepo/schemeTemplatePdf/" + schemeData.Name + ".pdf"
	//err = pdf.WritePdf(path)
	//if err != nil {
	//	return "", errors.New(err.Error())
	//}
	//
	//filePath := "https://" + context.Request.Host + path
	//
	//return filePath, nil

}

var html = `<h1>The frog king</h1>

<p>In olden times when wishing still <em>helped</em> one,
   there lived a king whose daughters were all beautiful,
   but the <span class="green">youngest</span> was so beautiful that the sun itself,
   which has seen so much, was
   <span style="font-weight: bold">astonished</span> whenever it
   shone in her face.</p>

<p>Close by the king's castle lay a great dark forest,
	and under an old lime-tree in the forest was a well,
	and when the day was very warm, the king's child
	went out into the forest and sat down by the side of
	the cool <span id="important">fountain</span>, and when she was bored she took a
	golden ball, and threw it up on high and caught it,
	and this ball was her favorite plaything.
</p>`

var css = `
body {
	font-family: serif;
    font-size: 12pt;
    line-height: 14pt;
}

p {
    margin-top: 8pt;
    margin-bottom: 2pt;
}

.green {
    color: green;
}

#important {
    color: rebeccapurple;
    font-weight: bolder;
    font-style: italic;
}`

func doThings(schemeData *SchemeOfWorkResponseInfo) (string, error) {
	path := "/fileRepo/schemeTemplatePdf/" + schemeData.Name + ".pdf"
	d, err := document.New(path)
	if err != nil {
		return "", errors.New(err.Error())
	}
	if d.AddCSS(css); err != nil {
		return "", errors.New(err.Error())
	}
	// loads the font families serif, sans and monospace
	if err = d.Frontend.LoadIncludedFonts(); err != nil {
		return "", errors.New(err.Error())
	}
	wd := bag.MustSp("280pt")
	colText := bag.MustSp("140pt")
	colImage := bag.MustSp("20pt")
	rowText := bag.MustSp("23cm")
	if err = d.OutputAt(html, wd, colText, rowText); err != nil {
		return "", errors.New(err.Error())
	}

	if err = d.OutputAt(`<img src="img/frogking-a.pdf" width="4cm" height="6cm">`, wd, colImage, rowText); err != nil {
		return "", errors.New(err.Error())
	}
	err = d.Finish()
	if err != nil {
		return "", errors.New(err.Error())
	}
	return path, nil

}

/*
func (schemeData *SchemeOfWorkResponseInfo) ReturnSchemeOfWorkPdf(context *gin.Context) (string, error) {
	// Sample JSON data
	reportData := Report{
		Title:   "Sample Report",
		Content: "<h1>This is the content of the report</h1>",
	}

	// Marshal the JSON data
	jsonData, err := json.Marshal(reportData)
	if err != nil {
		return "", errors.New(err.Error())
	}

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Add the JSON data including HTML content to the PDF
	pdf.SetFont("Arial", "", 16)
	pdf.MultiCell(0, 10, string(jsonData), "", "0", false)

	// Output the PDF to a file
	path := "/fileRepo/schemeTemplatePdf/" + schemeData.Name + ".pdf"
	err = pdf.OutputFileAndClose(path)
	if err != nil {
		return "", errors.New(err.Error())
	}

	filePath := "https://" + context.Request.Host + path

	return filePath, nil

}*/

///
//func (schemeData *SchemeOfWorkResponseInfo) ReturnSchemeOfWorkPdf(context *gin.Context) (string, error) {
//
//	pdf := gopdf.GoPdf{}
//
//	// Your HTML-like content
//
//	/*	htmlContent := "<h1>Hello, PDF!</h1>"
//		htmlContent += "<p>This is a sample PDF generated with HTML-like content.</p>"
//	*/
//
//	htmlContent := `
//<html>
//  <head>
//    <title>My Mega Service</title>
//    <link rel=\"stylesheet\" href=\"main.css\">
//    <style type=\"text/css\">body { color: #fff; }</style>
//  </head>
//
//  <body>
//    <div class="logo">
//      <a href="http://jaytaylor.com/"><img src="/logo-image.jpg" alt="Mega Service"/></a>
//    </div>
//
//    <h1>Welcome to your new account on my service!</h1>
//
//    <p>
//      Here is some more information:
//
//      <ul>
//        <li>Link 1: <a href="https://example.com">Example.com</a></li>
//        <li>Link 2: <a href="https://example2.com">Example2.com</a></li>
//        <li>Something else</li>
//      </ul>
//    </p>
//
//    <table>
//      <thead>
//        <tr><th>Header 1</th><th>Header 2</th></tr>
//      </thead>
//      <tfoot>
//        <tr><td>Footer 1</td><td>Footer 2</td></tr>
//      </tfoot>
//      <tbody>
//        <tr><td>Row 1 Col 1</td><td>Row 1 Col 2</td></tr>
//        <tr><td>Row 2 Col 1</td><td>Row 2 Col 2</td></tr>
//      </tbody>
//    </table>
//  </body>
//</html>`
//
//	/*	htmlText, err := html2text.FromString(htmlContent, html2text.Options{PrettyTables: true})
//		if err != nil {
//			panic(err)
//		}
//	*/
//	err := generatePDFFromHTML(&pdf, htmlContent)
//	if err != nil {
//		return "", errors.New(err.Error())
//	}
//
//	path := "/fileRepo/schemeTemplatePdf/" + schemeData.Name + ".pdf"
//	// Save the PDF to a file
//	err = pdf.WritePdf(path)
//	if err != nil {
//		return "", errors.New("----2. " + err.Error() + "----")
//	}
//
//	filePath := "https://" + context.Request.Host + path
//
//	return filePath, nil
//
//}
//
//func generatePDFFromHTML(pdf *gopdf.GoPdf, htmlContent string) error {
//	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
//	pdf.AddPage()
//
//	err := pdf.AddTTFFont("Arial", "/fileRepo/ttf/arial.ttf")
//	if err != nil {
//		return errors.New("---3. " + err.Error() + "----")
//	}
//	err = pdf.SetFont("Arial", "", 12)
//	if err != nil {
//		return errors.New("--4. " + err.Error() + "----")
//	}
//
//	// Simulate simple HTML rendering
//	pdf.SetX(20)
//	pdf.SetY(20)
//	pdf.SetTextColor(0, 0, 0)
//
//	err = pdf.Cell(nil, htmlContent)
//	if err != nil {
//		return errors.New("----4. " + err.Error() + "----")
//	}
//
//	return nil
//}
