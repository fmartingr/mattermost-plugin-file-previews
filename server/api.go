package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/mattermost/mattermost/server/public/plugin"
)

// ServeHTTP handles HTTP requests to the plugin.
// Requests will be `/plugins/com.mattermost.google-calendar/`...
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	// config := p.getConfiguration()

	// if err := config.IsValid(); err != nil {
	// 	http.Error(w, "This plugin is not configured.", http.StatusNotImplemented)
	// 	return
	// }

	switch path := r.URL.Path; path {
	case "/preview/":
		p.handlePreviewFile(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		http.NotFound(w, r)
	}
}

func (p *Plugin) handlePreviewFile(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.NotFound(w, r)
		return
	}

	fileID := r.URL.Query().Get("fileID")
	if fileID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contents, appErr := p.API.GetFile(fileID)
	if appErr != nil {
		p.API.LogError("error getting file contents", "file_id", fileID)
		return
	}

	fileInfo, appErr := p.API.GetFileInfo(fileID)
	if appErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpFile, err := os.CreateTemp("", "*."+fileInfo.Extension)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer tmpFile.Close()

	if _, errWrite := tmpFile.Write(contents); errWrite != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	input := &BackendFile{Path: tmpFile.Name(), File: tmpFile, FileInfo: fileInfo}

	output, err := p.backend.Convert(input, fileInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.API.LogError(err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/pdf")

	if _, err := io.Copy(w, output.File); err != nil {
		log.Printf("[Error] %s", err)
		return
	}
}
