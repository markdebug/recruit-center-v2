package util

import (
	"bytes"
	"io"
	"path/filepath"
	"strings"

	pdf "github.com/ledongthuc/pdf"
	"github.com/unidoc/unioffice/document"
)

// DocumentParser 定义文档解析接口
type DocumentParser interface {
	Parse(reader io.Reader) (string, error)
}

// PDFParser PDF解析器
type PDFParser struct{}

// WordParser Word解析器
type WordParser struct{}

// NewDocumentParser 根据文件类型创建对应的解析器
func NewDocumentParser(filename string) DocumentParser {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".pdf":
		return &PDFParser{}
	case ".doc", ".docx":
		return &WordParser{}
	default:
		return nil
	}
}

// Parse 实现PDF文档解析
func (p *PDFParser) Parse(reader io.Reader) (string, error) {
	// 读取内容到buffer
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(reader); err != nil {
		return "", err
	}

	// 解析PDF
	pdfBytes := buf.Bytes()
	pdfReader, err := pdf.NewReader(bytes.NewReader(pdfBytes), int64(len(pdfBytes)))
	if err != nil {
		return "", err
	}

	var content strings.Builder
	numPages := pdfReader.NumPage()

	for pageIndex := 1; pageIndex <= numPages; pageIndex++ {
		page := pdfReader.Page(pageIndex)
		text, err := page.GetPlainText(make(map[string]*pdf.Font))
		if err != nil {
			return "", err
		}
		content.WriteString(text)
	}

	return content.String(), nil
}

// Parse 实现Word文档解析
func (w *WordParser) Parse(reader io.Reader) (string, error) {
	// 读取内容到buffer
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(reader); err != nil {
		return "", err
	}

	// 解析Word文档
	readerAt := bytes.NewReader(buf.Bytes())
	doc, err := document.Read(readerAt, int64(readerAt.Len()))
	if err != nil {
		return "", err
	}

	var content strings.Builder
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			content.WriteString(run.Text())
		}
		content.WriteString("\n")
	}

	return content.String(), nil
}
