package main

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
	"golang.org/x/net/proxy"
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

func downloadFile(url string) error {
	// Configurar el proxy Tor (SOCKS5)
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9050", nil, proxy.Direct)
	if err != nil {
		return fmt.Errorf("error al configurar el proxy: %w", err)
	}

	// Crear un cliente HTTP que use el proxy Tor con timeout extendido
	httpTransport := &http.Transport{
		Dial: dialer.Dial,
	}
	client := &http.Client{
		Transport: httpTransport,
		Timeout:   4 * time.Hour, // Timeout extendido para archivos grandes
	}

	// Hacer la solicitud HTTP
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Comprobar el código de estado
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: %s", resp.Status)
	}

	// Extraer el nombre del archivo de la URL
	filename := filepath.Base(url)

	// Crear el archivo local con el nombre extraído
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// Inicializar la barra de progreso
	bar := progressbar.NewOptions(
		int(resp.ContentLength),
		progressbar.OptionSetDescription("Descargando..."),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionShowCount(),
	)

	// Copiar el contenido al archivo local y actualizar la barra de progreso
	_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	return err
}

func main() {
	// Asegurar limpieza al salir
	defer cleanup()

	// Extraer archivos embebidos
	if err := extractEmbeddedFiles(); err != nil {
		fmt.Println(err)
	}

	// Iniciar Tor
	go executeTor()
	fmt.Println("Iniciando cliente...")
	time.Sleep(5 * time.Second)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Introduce la URL del archivo .onion: ")
	url, _ := reader.ReadString('\n')

	// Eliminar caracteres no deseados
	url = strings.TrimSpace(url)

	err := downloadFile(url)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Archivo descargado exitosamente como:", filepath.Base(url))
	reader = bufio.NewReader(os.Stdin)
	fmt.Print("Presione enter para salir...")
	url, _ = reader.ReadString('\n')
	fmt.Println(url)
}

func executeTor() {
	cmd := exec.Command("./tori", "-f", "torrc")

	// Redirigir la salida a ioutil.Discard para no mostrar mensajes
	cmd.Stdout = nil //io.Discard
	cmd.Stderr = nil //io.Discard

	// Ejecutar el comando
	if err := cmd.Start(); err != nil {
		fmt.Println("Error al iniciar tor posiblemente ya esta ejecutandose")
	}

	// Esperar a que el comando termine
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error al iniciar tor posiblemente ya esta ejecutandose")
	}
}
