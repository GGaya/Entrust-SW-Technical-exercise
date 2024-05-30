package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/jung-kurt/gofpdf"
)

// Prints the Usage message
func usage() {
	fmt.Println("Usage: ./pagination document_name [txt | pdf]")

}

// Input: file of a document (.txt)
// Output: list of strings with each page paginated (each index it's a page)
func process_text_file(file *os.File) []string {

	// Read, scan and divide the text in words
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)

	// Initialization of vars to count chars and lines
	var count_chars int = 0
	var page string
	var text []string
	var count int = 0
	// Code of all the types of '-'
	re := regexp.MustCompile("[\u2010-\u2015\u2212\u2E3A\u2E3B]")

	// We read word for word until we finish
	for scanner.Scan() {
		word := scanner.Text()
		word_size := len(word)
		// Check we can add another word to the line
		if count_chars+word_size < 80 {
			count_chars += word_size + 1
			page += word + " "

			// Same without the ' '
		} else if count_chars+word_size == 80 {
			count_chars += word_size
			page += word

			// We can't add another word so we restart counters for the next line
		} else {

			count_chars = word_size + 1
			count++
			// Checks if we need a new page
			if count >= 25 {
				// Replace invalid types of '-'
				page = re.ReplaceAllString(page, "-")
				text = append(text, page)
				// We start preparing the next page and reset count of lines
				page = word + " "
				count = 0

				// Preparing the next line with a \n
			} else {
				page += "\n" + word + " "
			}

		}
	}
	// Replace invalid types of '-'
	page = re.ReplaceAllString(page, "-")
	// Save the page paginated
	text = append(text, page)

	// Checking error
	if err1 := scanner.Err(); err1 != nil {
		fmt.Println(err1)
	}

	return text
}

// Configurates the pdf parameters like page size, font, etc.
func pdf_config_ini() *gofpdf.Fpdf {

	pdf := gofpdf.New("P", "mm", "A4", "") // Type of page
	pdf.SetFont("Arial", "", 11)           // Font
	pdf.SetLeftMargin(30)                  // Margin
	pdf.SetFooterFunc(func() {             // Footnote
		pdf.SetY(-30)
		pdf.SetDrawColor(0, 0, 0)
		pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
		pdf.Ln(10)
		pdf.SetFont("Arial", "B", 9)
		pdf.CellFormat(150, 10, fmt.Sprintf("Page %d", pdf.PageNo()), "", 0, "C",
			false, 0, "")
	})
	return pdf
}

// Input: lists of strigns with the pages paginated
// Output: pdf object ready to be generated
func generate_pdf_file(text []string) *gofpdf.Fpdf {
	pdf := pdf_config_ini()
	// We add a pdf page and write on it
	for _, page := range text {
		pdf.AddPage()
		pdf.MultiCell(190, 5, page, "", "", false)

	}
	return pdf
}

// Input: lists of strigns with the pages paginated and the document name
// Output: -
func generate_txt_file(text []string, doc_name string) {

	// Creation of the .txt file in the specified directory and error check
	file, err := os.Create("./pg_docs/" + doc_name + ".txt")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	// Write page by page making the footnote
	datawriter := bufio.NewWriter(file)
	for i, page := range text {
		page += fmt.Sprintf("\n\n\n------------------------------ Page %d "+
			"------------------------------\n\n\n", i+1)
		_, _ = datawriter.WriteString(page)

	}
	datawriter.Flush()
	file.Close()
}

func main() {
	// Usage control
	if (len(os.Args) < 2) || (len(os.Args) > 3) {
		usage()
		return
	}

	// Check type of document for the paginated solution
	var ext_type bool = false
	if len(os.Args) == 3 {
		if os.Args[2] == "pdf" {
			ext_type = true
		}
	}

	// Open the document and error check
	var doc_name string = os.Args[1]
	file, err1 := os.Open(doc_name + ".txt")
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	defer file.Close()

	// Save the list of pages in text
	text := process_text_file(file)

	// Prepares the directory for the paginated docs
	if _, err := os.Stat("./pg_docs"); os.IsNotExist(err) {
		os.Mkdir("./pg_docs", 0755)
	}

	if !ext_type {
		// txt output generation
		generate_txt_file(text, doc_name)
	} else {
		// pdf output generation
		pdf := generate_pdf_file(text)
		// Generate the pdf file and error check
		err2 := pdf.OutputFileAndClose("./pg_docs/" + doc_name + ".pdf")
		if err2 != nil {
			panic(err2)
		}
	}

}
