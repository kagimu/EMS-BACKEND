package pdfGenerator

import (
	"bytes"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"html/template"
)

// GeneratePdf **@description: ReturnNewData returns a new data object //**@params: templateName string, data interface{} //**@returns: *wkhtmltopdf.PDFGenerator, error
// **@returns: NewData
//
//	//to save the pdf use the returned pdfg object
//	err = pdfg.WriteFile("proof.pdf")
//	if err != nil {
//		panic(err)
//	}
//
//	// to convert pdf to base64
//	base := base64.StdEncoding.EncodeToString(pdfg.Bytes())
//	fmt.Println(base) /*
func GeneratePdf(templateName string, data interface{}) (*wkhtmltopdf.PDFGenerator, error) {
	//define the data
	if templateName == "" {
		templateName = "sample_template"
		data = ReturnNewData()
	}

	//read the template
	var templ *template.Template
	var err error
	classReportTemplateUrl := "./fileRepo/classReportTemplates/" + templateName + ".html"
	templ, err = template.ParseFiles(classReportTemplateUrl)
	if err != nil {
		return nil, err
	}

	// apply the parsed HTML template data and keep the result in a Buffer
	var doc bytes.Buffer
	err = templ.Execute(&doc, data)
	if err != nil {
		return nil, err
	}

	//convert html to pdf
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(&doc))
	err = pdfg.Create()
	if err != nil {
		return nil, err
	}

	return pdfg, nil
}
