package main

import (
	"fmt"
	"github.com/aymerick/raymond"
	"io"
	_ "io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var listTemplate *raymond.Template

type listFileInfo struct {
	Name  string
	IsDir bool
}

func init() {
	listContent, _ := Asset("templates/list.handlebars")

	listTemplate = raymond.MustParse(string(listContent))
}

func handleHttp(w http.ResponseWriter, r *http.Request) {
	dir, _ := os.Getwd()
	reqFilePath := filepath.Join(dir, r.URL.Path)

	log.Printf("Got request for file : %s", reqFilePath)

	file, err := os.Open(reqFilePath)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to open requested file/directory : %s", r.URL.Path), 400)
		log.Printf("Failed to find file %s", reqFilePath)
		return
	}

	stats, _ := file.Stat()

	if stats.IsDir() {
		serveDir(w, r, file)
	} else {
		serveFile(w, file)
	}
}

func serveDir(w http.ResponseWriter, r *http.Request, file *os.File) {
	dir, _ := os.Getwd()
	reqFilePath := filepath.Join(dir, r.URL.Path)

	files, _ := file.Readdir(0)

	fileInfo := make([]map[string]interface{}, len(files))

	for i, dirFileInfo := range files {
		rel, _ := filepath.Rel(dir, reqFilePath+string(os.PathSeparator)+dirFileInfo.Name())

		fileInfo[i] = map[string]interface{}{
			"name":       dirFileInfo.Name(),
			"isDir":      dirFileInfo.IsDir(),
			"pathToFile": "/" + rel,
		}
	}

	result := listTemplate.MustExec(map[string]interface{}{
		"files":    fileInfo,
		"isTopDir": r.URL.Path == "/",
	})

	fmt.Fprint(w, result)
}

func serveFile(w http.ResponseWriter, file *os.File) {
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(file.Name()))
	_, err := io.Copy(w, file)

	if err != nil {
		http.Error(w, "Failed to read file content", http.StatusInternalServerError)
		return
	}

	log.Printf("File %s sent to client", file.Name())
}
