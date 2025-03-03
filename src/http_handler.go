package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	HttpServeMux        *http.ServeMux
	LoadedFrontendFiles []os.DirEntry
)

func loadFrontendFiles() {
	files, err := os.ReadDir(LoadedConfig.FrontendPath)
	if err != nil {
		fmt.Println("Error reading frontend directory:", err)
		return
	}

	// Save the files to LoadedFrontendFiles
	LoadedFrontendFiles = files
}

func handleFrontendRequest(w http.ResponseWriter, r *http.Request, path string) bool {
	if path == "/" {
		serveFile(w, r, "index.html")
		return true
	}

	return serveStaticFile(w, path)
}

func serveFile(w http.ResponseWriter, r *http.Request, filename string) bool {
	filePath := filepath.Join(LoadedConfig.FrontendPath, filename)

	if _, err := os.Stat(filePath); err != nil {
		http.NotFound(w, nil)
		return false
	}

	http.ServeFile(w, r, filePath)
	return true
}

// serveStaticFile checks if the requested file exists and serves it
func serveStaticFile(w http.ResponseWriter, path string) bool {
	for _, file := range LoadedFrontendFiles {
		if file.Name() == path {
			filePath := filepath.Join(LoadedConfig.FrontendPath, file.Name())
			http.ServeFile(w, nil, filePath)
			return true
		}
	}

	return false
}

func getClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For") // Mainly for NGINX support
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0]) // First IP is the real client
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	return r.RemoteAddr
}

func logClient(apiHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		address := getClientIP(r)

		log.Printf("Request [%s] | '%s' from %s", r.Method, r.URL.String(), address)

		if handleFrontendRequest(w, r, r.URL.Path) {
			return
		}

		apiHandler.ServeHTTP(w, r)
	})
}

func StartHttpServer() {
	BuildApi()
	loadFrontendFiles()

	HttpServeMux = http.NewServeMux()
	apiRequestFunc := http.HandlerFunc(HandleApiRequest)
	HttpServeMux.Handle("/", logClient(apiRequestFunc))

	fmt.Println("Starting server at PORT:", LoadedConfig.Port)

	if err := http.ListenAndServe(":"+strconv.Itoa(int(LoadedConfig.Port)), HttpServeMux); err != nil {
		fmt.Println(err)
		return
	}
}
