package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/pdftron/pdftron-go/v2"
)

const (
	BackendPDFTron   = "pdftron"
	BackendGotenberg = "gotenberg"
)

type BackendFile struct {
	Path     string
	File     *os.File
	FileInfo *model.FileInfo
}

type Backend interface {
	Init(c configuration) error
	Convert(*BackendFile, *model.FileInfo) (*BackendFile, error)
}

type PDFTronBackend struct{}

func (pb *PDFTronBackend) Init(c configuration) error {
	pdftron.PDFNetInitialize(c.PDFTronLicenseKey)
	return nil
}

func (pb *PDFTronBackend) simpleConvert(input, output *BackendFile, extension string) error {
	pdfdoc := pdftron.NewPDFDoc()
	options := pdftron.NewConversionOptions()
	options.SetFileExtension(extension)
	pdftron.ConvertOfficeToPDF(pdfdoc, input.Path, options)
	pdfdoc.Save(output.Path, uint(pdftron.SDFDocE_linearized))
	return nil
}

func (pb *PDFTronBackend) Convert(input *BackendFile, fileInfo *model.FileInfo) (*BackendFile, error) {
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		return nil, fmt.Errorf("error creating temp file: %w", err)
	}
	output := &BackendFile{
		Path: tmpFile.Name(),
		File: tmpFile,
	}
	if err := pb.simpleConvert(input, output, fileInfo.Extension); err != nil {
		log.Println(err)
		return nil, err
	}
	return output, nil
}

func NewPDFTronBackend() Backend {
	return &PDFTronBackend{}
}

type GotenbergBackend struct {
	serverURL string
}

func (gb *GotenbergBackend) healthCheck() error {
	gotenbergURL, err := url.JoinPath(gb.serverURL, "/health")
	if err != nil {
		return fmt.Errorf("error creating url: %w", err)
	}
	resp, err := http.Get(gotenbergURL)
	if err != nil {
		return fmt.Errorf("error sending health check request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gotenberg is not ready")
	}

	return nil
}

func (gb *GotenbergBackend) Init(c configuration) error {
	gb.serverURL = c.GotenbergServerURL

	if c.GotenbergServerURL == "" {
		return fmt.Errorf("gotenberg backend requires setting up the server URL")
	}

	return gb.healthCheck()
}

func (gb *GotenbergBackend) Convert(input *BackendFile, _ *model.FileInfo) (*BackendFile, error) {
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		return nil, fmt.Errorf("error creating temp file: %w", err)
	}
	output := &BackendFile{
		Path: tmpFile.Name(),
		File: tmpFile,
	}

	log.Println(output.Path)

	url, err := url.JoinPath(gb.serverURL, "/forms/libreoffice/convert")
	if err != nil {
		return nil, fmt.Errorf("error creating url: %w", err)
	}
	if err := uploadFile(url, input, output); err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	return output, nil
}

func NewGotenbergBackend() Backend {
	return &GotenbergBackend{}
}

func uploadFile(url string, input, output *BackendFile) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	inputFile, err := os.Open(input.Path)
	if err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}

	part, _ := writer.CreateFormFile("file", input.FileInfo.Name)
	if _, err := io.Copy(part, inputFile); err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}
	writer.Close()

	r, _ := http.NewRequest("POST", url, body)
	r.Header.Add("Content-Type", writer.FormDataContentType())

	debugReq, _ := httputil.DumpRequest(r, true)
	log.Println(string(debugReq))

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}

	defer resp.Body.Close()

	debugRes, _ := httputil.DumpResponse(resp, true)
	log.Println(string(debugRes))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Error: %d", resp.StatusCode)
	}

	if _, err := io.Copy(output.File, resp.Body); err != nil {
		return fmt.Errorf("error copying gotenberg output to file: %w", err)
	}

	output.File.Seek(0, 0)

	return nil
}
