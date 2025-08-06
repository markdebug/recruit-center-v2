package utils

import (
	"os"
	"testing"
)

func TestPDFParser(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "Valid PDF",
			filename: "testdata/test.pdf",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.filename)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Failed to open file: %v", err)
				}
				return
			}
			defer file.Close()

			parser := &PDFParser{}
			_, err = parser.Parse(file)
			if (err != nil) != tt.wantErr {
				t.Errorf("PDFParser.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWordParser(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "Valid DOCX",
			filename: "testdata/test.docx",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.filename)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Failed to open file: %v", err)
				}
				return
			}
			defer file.Close()

			parser := &WordParser{}
			_, err = parser.Parse(file)
			if (err != nil) != tt.wantErr {
				t.Errorf("WordParser.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
