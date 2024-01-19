package responses

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/phpdave11/gofpdf"
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

	path, err := ExampleFpdf_HTMLBasicNew(schemeData, context)
	if err != nil {
		log.Fatal(err)
	}

	//filePath := "https://" + context.Request.Host + path

	return path, nil

}

func ExampleFpdf_HTMLBasicNew(schemeData *SchemeOfWorkResponseInfo, context *gin.Context) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	// First page: manual local link
	pdf.AddPage()

	/*  pdf.AddTTFFont("Arial", "/fileRepo/ttf/arial.ttf")
	if err != nil {
		return errors.New("---3. " + err.Error() + "----")
	}*/
	//pdf.AddFont("Arial", "","/fileRepo/ttf/arial.ttf")
	pdf.SetFont("Arial", "", 12)
	//pdf.SetFont("Helvetica", "", 20)
	_, lineHt := pdf.GetFontSize()
	//pdf.Write(lineHt, "To find out what's new in this tutorial, click ")
	//pdf.SetFont("", "U", 0)
	//link := pdf.AddLink()
	//pdf.WriteLinkID(lineHt, "here", link)
	//pdf.SetFont("", "", 0)
	// Second page: image link and basic HTML with link
	//pdf.AddPage()
	//pdf.SetLink(link, 0, -1)
	//pdf.Image(example.ImageFile("logo.png"), 10, 12, 30, 0, false, "", 0, "http://www.fpdf.org")
	pdf.SetLeftMargin(45)
	pdf.SetFontSize(14)
	_, lineHt = pdf.GetFontSize()
	htmlStr := `You can now easily print text mixing different styles: <b>bold</b>, ` +
		`<i>italic</i>, <u>underlined</u>, or <b><i><u>all at once</u></i></b>!<br><br>` +
		`<center>You can also center text.</center>` +
		`<right>Or align it to the right.</right>` +
		`You can also insert links on text, such as ` +
		`<a href="http://www.fpdf.org">www.fpdf.org</a>, or on an image: click on the logo.`
	html := pdf.HTMLBasicNew()
	html.Write(lineHt, htmlStr)
	path := "/fileRepo/schemeTemplatePdf/" + schemeData.Name + ".pdf"
	err := pdf.OutputFileAndClose(path)
	if err != nil {
		return "", errors.New(err.Error())
	}

	filePath := "https://" + context.Request.Host + path

	return filePath, nil
}

//func doThings(schemeData *SchemeOfWorkResponseInfo) (string, error) {
//	path := "/fileRepo/schemeTemplatePdf/" + schemeData.Name + ".pdf"
//	htmlData := generateHtml(schemeData)
//	d, err := document.New(path)
//	if err != nil {
//		return "", errors.New(err.Error())
//	}
//	/*if d.AddCSS(css); err != nil {
//		return "", errors.New(err.Error())
//	}*/
//	// loads the font families serif, sans and monospace
//	if err = d.Frontend.LoadIncludedFonts(); err != nil {
//		return "", errors.New(err.Error())
//	}
//	wd := bag.MustSp("280pt")
//	colText := bag.MustSp("140pt")
//	//colImage := bag.MustSp("20pt")
//	rowText := bag.MustSp("23cm")
//	if err = d.OutputAt(htmlData, wd, colText, rowText); err != nil {
//		return "", errors.New(err.Error())
//	}
//
//	/*if err = d.OutputAt(`<img src="img/frogking-a.pdf" width="4cm" height="6cm">`, wd, colImage, rowText); err != nil {
//		return "", errors.New(err.Error())
//	}*/
//	err = d.Finish()
//	if err != nil {
//		return "", errors.New(err.Error())
//	}
//	return path, nil
//
//}

//func generateHtml(schemeData *SchemeOfWorkResponseInfo) string {
//	var html = `<h1>The frog king</h1><br /> <br />`
//	html += `<p>In olden times when wishing still <em>helped</em> one,
//		   there lived a king whose daughters were all beautiful,
//		   but the <span class="green">youngest</span> was so beautiful that the sun itself,
//		   which has seen so much, was
//		   <span style="font-weight: bold">astonished</span> whenever it
//		   shone in her face.</p><br /><br />`
//
//	html += `<p>Close by the king's castle lay a great dark forest,
//			and under an old lime-tree in the forest was a well,
//			and when the day was very warm, the king's child
//			went out into the forest and sat down by the side of
//			the cool <span id="important">fountain</span>, and when she was bored she took a
//			golden ball, and threw it up on high and caught it,
//			and this ball was her favorite plaything.
//		</p>`
//	html += `
//			<h1 style="text-align: center;">CSS-styled Inline Div Table with Two Columns</h1>
//			<div style="width: 100%; margin-bottom: 20px;">
//				<div style="display: flex; font-weight: bold; background-color: #f2f2f2; margin-bottom: 8px;">
//					<div style="flex: 1; border: 1px solid #dddddd; padding: 8px;">Header 1</div>
//					<div style="flex: 1; border: 1px solid #dddddd; padding: 8px;">Header 2</div>
//				</div>
//				<div style="display: flex; margin-bottom: 8px;">
//					<div style="flex: 1; border: 1px solid #dddddd; padding: 8px;">Data 1</div>
//					<div style="flex: 1; border: 1px solid #dddddd; padding: 8px;">Data 2</div>
//				</div>
//				<div style="display: flex;">
//					<div style="flex: 1; border: 1px solid #dddddd; padding: 8px;">Data 3</div>
//					<div style="flex: 1; border: 1px solid #dddddd; padding: 8px;">Data 4</div>
//				</div>
//			</div>`
//	return html
//}

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
