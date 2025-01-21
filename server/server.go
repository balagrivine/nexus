package server

import (
	"errors"
	"log"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/balagrivine/nexus/server/http"
)

var logger *slog.Logger = configLogger()

type HTTPServer struct {
	ListenAddr string
	Listener   net.Listener
	Ready      chan struct{} // Signals server's readiness to accept connection
	Quit       chan struct{} // Signals shutdown
}

// NewHTTPServer creates a new nexus server instance
// This functions accepts the address on which the server is to run
func NewHTTPServer(listenAddr string) *HTTPServer {
	return &HTTPServer{
		ListenAddr: listenAddr,
		Ready:      make(chan struct{}),
		Quit:       make(chan struct{}),
	}
}

// Starts an instance of the HTTP server
func (srv *HTTPServer) Start() error {

	listener, err := net.Listen("tcp", srv.ListenAddr)
	if err != nil {
		logger.Error("unable to listen to connection", slog.Any("error", err))
		return err
	}

	srv.Listener = listener
	close(srv.Ready) // Signal readiness to handle connection

	srv.acceptConnection()

	<-srv.Quit // Block until a shutdown signal is received
	log.Println("Shutting down server...")
	srv.Listener.Close()

	return nil
}

// Close signals that the server is shutting down
func (srv *HTTPServer) ShutDown() {
	close(srv.Quit)
}
func (srv *HTTPServer) acceptConnection() {
	log.Print("Ready to accept connection...")
	for {
		conn, err := srv.Listener.Accept()
		if err != nil {
			logger.Warn("could not accept connection", slog.Any("error", err))
			continue
		}

		go srv.handleConnection(conn)
	}
}

func (srv *HTTPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		buffer := make([]byte, 6048)
		_, err := conn.Read(buffer)
		if err != nil {
			//TODO
		}

		//data := buffer[:n]
		//fmt.Println(string(data))
		//response := http.GetResponseWriter(conn)
		//_, err = http.Decode(data)
		//if err != nil {
		//	response.Send([]byte("Invalid Method"), http.HTTP_405_NOT_ALLOWED)
		//}

		//respBody := []byte("Hello World!")
		//response.AddHeader("Content-Type", "text/plain")
		//response.AddHeader("Content-Length", string(len(respBody)))
		//response.AddHeader("Content-Type", "text/plain")
		//response.AddHeader("Accept-Ranges", "bytes")

		//response.Send(respBody, http.HTTP_200_OK)
		conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World\r\n"))
		return
	}
}

func processFilePath(path string) (*os.File, int64, string, error) {
	// Root directory from which to serve files
	// Any directory access out of the rootDir should not be permitted
	rootDir := "www"

	cleanPath := filepath.Clean(path)
	path = filepath.Join(rootDir, cleanPath)

	// Serve the default HTML file if only directory name is provided
	if strings.HasSuffix(path, "/") {
		path = filepath.Join(path, "index.html")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, 0, "", err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil || info.IsDir() {
		return nil, 0, "", errors.New("Not found")
	}

	contentType := getContentType(path)
	fileSize := info.Size()

	return file, fileSize, contentType, nil
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
		return "application/javascript"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}

// Creates and returns a custom Logger instance
func configLogger() *slog.Logger {
	logLevel := &slog.LevelVar{}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	return logger
}
