package helper

import (
	"bytes"
	"log"
	"text/template"

	"github.com/jung-kurt/gofpdf"
)

type TemplateDataAgreementLetter struct {
	Title           string
	Date            string
	Admin           string
	BorrowerName    string
	BorrowerAddress string
	PrincipalAmount string
	Rate            string
}

func (data *TemplateDataAgreementLetter) GenerateAgreementLetter(templatePath string) ([]byte, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	// Execute the template with data
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	// Create a new PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16) // "B" for bold, 16 for size
	pdf.CellFormat(0, 10, data.Title, "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Convert the template's HTML to PDF
	pdf.SetFont("Arial", "", 12)
	html := pdf.HTMLBasicNew()
	html.Write(10, tpl.String())

	// Save the PDF to a file
	var pdfBytes bytes.Buffer
	err = pdf.Output(&pdfBytes)
	if err != nil {
		return nil, err
	}

	return pdfBytes.Bytes(), nil
}
