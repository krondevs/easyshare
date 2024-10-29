package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

//go:embed tori
//go:embed torrc
var embeddedFiles embed.FS

func extractEmbeddedFiles() error {
	// Extraer tori
	torExeData, err := embeddedFiles.ReadFile("tori")
	if err != nil {
		return fmt.Errorf("error al leer tori embebido: %w", err)
	}
	err = os.WriteFile("tori", torExeData, 0755)
	if err != nil {
		return fmt.Errorf("error al extraer tori: %w", err)
	}

	// Extraer torrc
	torrcData, err := embeddedFiles.ReadFile("torrc")
	if err != nil {
		return fmt.Errorf("error al leer torrc embebido: %w", err)
	}
	err = os.WriteFile("torrc", torrcData, 0644)
	if err != nil {
		return fmt.Errorf("error al extraer torrc: %w", err)
	}

	return nil
}

func cleanup() {
	// Eliminar archivos extraídos al salir
	os.Remove("tori")
	os.Remove("torrc")
}

func listFiles(dir string) ([]string, error) {
	var files []string
	fileInfo, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range fileInfo {
		if !file.IsDir() {
			files = append(files, file.Name())
		}
	}
	return files, nil
}

func ensureHiddenServiceDir() error {
	if err := os.MkdirAll("hidden_service", 0700); err != nil {
		return fmt.Errorf("error al crear directorio hidden_service: %w", err)
	}
	return nil
}

// Middleware para logging y manejo de timeouts
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Iniciando request: %s %s", r.Method, r.URL.Path)

		// Envolver ResponseWriter para capturar el código de estado
		wrapped := wrapResponseWriter(w)

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		log.Printf("Completado %s %s con código %d en %v",
			r.Method, r.URL.Path, wrapped.status, duration)
	})
}

// ResponseWriter personalizado para capturar el código de estado
type responseWriter struct {
	http.ResponseWriter
	status  int
	written bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.status = http.StatusOK
		rw.written = true
	}
	return rw.ResponseWriter.Write(b)
}

func main() {
	// Asegurar limpieza al salir
	defer cleanup()

	// Extraer archivos embebidos
	if err := extractEmbeddedFiles(); err != nil {
		fmt.Println(err)
	}

	// Crear directorio hidden_service si no existe
	if err := ensureHiddenServiceDir(); err != nil {
		fmt.Println(err)
	}

	go executeTor()
	fmt.Println("Starting server...")
	time.Sleep(5 * time.Second)

	// Esperar a que se cree el archivo hostname
	filePath := filepath.Join("hidden_service", "hostname")
	maxAttempts := 30
	var address string
	var err error

	for i := 0; i < maxAttempts; i++ {
		if data, err := os.ReadFile(filePath); err == nil {
			address = strings.TrimSpace(string(data))
			break
		}
		time.Sleep(1 * time.Second)
	}

	if address == "" {
		log.Fatal("No se pudo leer el archivo hostname después de varios intentos")
	}

	// Listar archivos en el directorio actual
	files, err := listFiles(".")
	if err != nil {
		log.Printf("Error al listar archivos: %v", err)
	} else {
		fmt.Println("\nArchivos disponibles en el servicio:")
		fmt.Printf("URL base: http://%s/\n", address)
		fmt.Println("Lista de archivos:")
		for _, file := range files {
			if !strings.HasPrefix(file, ".") && file != "tori" && file != "torrc" {
				fmt.Printf("- http://%s/%s\n", address, file)
			}
		}
		fmt.Println()
	}

	// Configurar el servidor con timeouts extendidos
	server := &http.Server{
		Addr:    ":8085",
		Handler: loggingMiddleware(http.FileServer(http.Dir("."))),
		// Timeout para leer el cuerpo completo de la request
		ReadTimeout: 10 * time.Minute,
		// Timeout para escribir la respuesta
		WriteTimeout: 4 * time.Hour,
		// Timeout para mantener conexiones keepalive
		IdleTimeout: 10 * time.Hour,
		// Tiempo máximo para leer la request completa incluyendo el body
		ReadHeaderTimeout: 1 * time.Minute,
	}

	fmt.Println("Servidor local en http://localhost:8085/")
	fmt.Printf("Servidor Tor en http://%s/\n", address)

	// Iniciar el servidor con manejo de errores mejorado
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Error en el servidor:", err)
	}
}

func executeTor() {
	cmd := exec.Command("./tori", "-f", "torrc")

	// Capturar salida de tor
	cmd.Stdout = nil //os.Stdout
	cmd.Stderr = nil //os.Stderr

	// Ejecutar el comando
	if err := cmd.Start(); err != nil {
		fmt.Println("Error al iniciar tor:", err)
	}

	// Esperar a que el comando termine
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error durante la ejecución de tor:", err)
	}

	log.Println("tori ejecutado exitosamente")
}
