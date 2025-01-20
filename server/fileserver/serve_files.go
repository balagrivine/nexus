package fileserver

import (
	"io"
	"net"
	"os"
	"path/filepath"
)

func ServeStaticFile(conn net.Conn, path string) {
	if path == "/" {
		path = "/index.html"
	}

	staticDir := "www"
	fullFilePath := filepath.Join(staticDir, filepath.Clean(path))

	file, err := os.Open(fullFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			conn.Write([]byte("HTTP/1.1 404 Not found \r\n\r\n"))
		} else {
			conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
		}

		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil || info.IsDir() {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		return
	}

	contentType := getContentType(fullFilePath)
	conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
	conn.Write([]byte("Content-Type: " + contentType + "\r\n"))
	conn.Write([]byte("Content-Length: " + fmt.Sprint(info.Size()) + "\r\n"))
	conn.Write([]byte("\r\n"))

	// Write the file content
	io.Copy(conn, file)
}

func getContentType(filePath string) string {
	ext := filepath.Ext(filePath)

	switch ext {
	case ".html":
		return "text/html"
	case ".text":
		return "text/txt"
	case ".css":
		return "text/css"
	case ".js":
		"application/javascript"
	case ".png":
		return "image/png"
	case ".jpg". ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}
