package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	log "github.com/golang/glog"
)

type HttpPathMapping struct {
	HttpPath string
	FilePath string
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("image")
	fileName := r.FormValue("name")

	if err != nil {
		log.Infof("Failed to handle upload request with error: %v\n", err)
	}
	defer file.Close()

	out, err := os.Create(fmt.Sprintf("/vagrant/%s", fileName))

	if err != nil {
		fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
		return
	}
	defer out.Close()

	// Write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		log.Infoln(w, err)
	}

	fmt.Fprintf(w, "File uploaded successfully : ")
	fmt.Fprintf(w, header.Filename)
}

func registerUploadHandler() {
	http.HandleFunc("/", uploadHandler)
}

func registerDownloadHandler(fileToServe HttpPathMapping) {
	log.Infof("httpPath: %v\n", fileToServe.HttpPath)
	log.Infof("filePath: %v\n", fileToServe.FilePath)

	http.HandleFunc(fileToServe.HttpPath, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fileToServe.FilePath)
	})
}

func registerDownloadHandlers(filesToServe []HttpPathMapping) {
	for _, m := range filesToServe {
		registerDownloadHandler(m)
	}
}

func GetHttpPath(path string) string {
	// Create base path (http://foobar:5000/<base>)
	pathSplit := strings.Split(path, "/")
	var base string
	if len(pathSplit) > 0 {
		base = pathSplit[len(pathSplit)-1]
	} else {
		base = path
	}

	return "/" + base
}

func GetDefaultMappings(filePaths []string) []HttpPathMapping {
	mappings := []HttpPathMapping{}

	for _, f := range filePaths {
		m := HttpPathMapping{
			HttpPath: GetHttpPath(f),
			FilePath: f,
		}

		mappings = append(mappings, m)
	}

	return mappings
}

func StartHttpServer(address string, filesToServe []HttpPathMapping) {
	registerDownloadHandlers(filesToServe)
	registerUploadHandler()
	go http.ListenAndServe(address, nil)
}
