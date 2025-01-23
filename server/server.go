package server

import (
	"bytes"
	"errors"
	"io"
	"log"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/balagrivine/nexus/server/http"
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
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			//TODO
		}

		data := buffer[:bytesRead]
		response := http.GetResponseWriter(conn)
		request, err := http.Decode(data)
		if err != nil {
			response.AddHeader("Content-Type", "text/plain")
			response.Send([]byte("Invalid Method\n"), http.HTTP_405_NOT_ALLOWED)
			break
		}

		srv.processStaticPath(response, request)
		return
	}
}

func (srv *HTTPServer) processStaticPath(response *http.Response, request *http.Request) {
	file, contentType, err := srv.processFilePath(request.Path)
	if err != nil {
		response.Send([]byte("404 Not Found\n"), http.HTTP_404_NOT_FOUND)
		return
	}
	defer file.Close()

	dataBuf := bytes.Buffer{}
	
	n, err := io.Copy(&dataBuf, file)
	if err != nil {
		response.Send([]byte("500 Internal Server Error\n"), http.HTTP_500_SERVER_ERROR)
		return
	}

	response.AddHeader("Content-Type", contentType)
	response.Send(dataBuf.Bytes()[:n], http.HTTP_200_OK)

	return
}

func (srv *HTTPServer) processFilePath(path string) (*os.File, string, error) {
	// Root directory from which to serve files
	// Any directory access out of the rootDir should not be permitted
	rootDir := "static"

	// Serve the default HTML file if only directory name is provided
	if strings.HasSuffix(path, "/") {
		path = filepath.Join(path, "index.html")
	}

	cleanPath := filepath.Clean(path)
	path = filepath.Join(rootDir, cleanPath)

	file, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}

	info, err := file.Stat()
	if err != nil || info.IsDir() {
		return nil, "", errors.New("not found")
	}

	contentType := getContentType(path)

	return file, contentType, nil
}

// getContentType determines the content type of a file based
// on its extension
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
