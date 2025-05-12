package gen

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/C0d3-5t3w/go-ssg/inc/config"
	"github.com/russross/blackfriday/v2"
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

func GEN(cfg *config.Config) error {

	files, err := ioutil.ReadDir(cfg.ContentDir)
	if err != nil {
		return fmt.Errorf("failed to read content directory '%s': %v", cfg.ContentDir, err)
	}

	if _, err := os.Stat(cfg.OutputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory '%s': %v", cfg.OutputDir, err)
		}
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) != ".md" {
			fmt.Printf("Skipping non-markdown file: %s\n", file.Name())
			continue
		}

		filePath := filepath.Join(cfg.ContentDir, file.Name())
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %v", filePath, err)
		}

		htmlContent := blackfriday.Run(content)

		page := Page{
			Title: strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())),
			Body:  template.HTML(htmlContent),
		}

		outputFilePath := filepath.Join(cfg.OutputDir, strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))+".html")
		if err := generateHTML(page, outputFilePath); err != nil {
			return fmt.Errorf("failed to generate HTML for %s: %v", filePath, err)
		}
		fmt.Printf("Generated: %s\n", outputFilePath)
	}
	fmt.Printf("Site generation complete. Output in '%s'.\n", cfg.OutputDir)
	return nil
}
