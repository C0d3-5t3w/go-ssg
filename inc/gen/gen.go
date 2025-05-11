package gen

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday/v2"
)

const (
	contentDir = "content"
	outputDir  = "output"
)

const htmlTemplate = `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>{{ .Title }}</title>
</head>
<body>
{{ .Body }}
</body>
</html>`

type Page struct {
	Title string
	Body  template.HTML
}

func generateHTML(page Page, outputFilePath string) error {
	tmpl, err := template.New("page").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	file, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %v", outputFilePath, err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, page); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	return nil
}

func GEN() error {
	// Read the content directory
	files, err := ioutil.ReadDir(contentDir)
	if err != nil {
		return fmt.Errorf("failed to read content directory: %v", err)
	}

	// Create the output directory if it doesn't exist
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.Mkdir(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %v", err)
		}
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(contentDir, file.Name())
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %v", filePath, err)
		}

		htmlContent := blackfriday.Run(content)

		page := Page{
			Title: strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())),
			Body:  template.HTML(htmlContent),
		}

		outputFilePath := filepath.Join(outputDir, strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))+".html")
		if err := generateHTML(page, outputFilePath); err != nil {
			return fmt.Errorf("failed to generate HTML for %s: %v", filePath, err)
		}
	}

	return nil
}
