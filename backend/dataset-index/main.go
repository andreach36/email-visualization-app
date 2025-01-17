package datasetindex

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Leer los directorios
func folderList(path string) []string {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Printf("Error leyendo directorio %s: %v", path, err)
		return nil
	}

	var listFolders []string

	for _, file := range files {
		if file.IsDir() {
			listFolders = append(listFolders, file.Name())
		}
	}

	return listFolders
}

// Leer los archivos
func fileList(path string) []string {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Printf("Error leyendo archivo %s: %v", path, err)
		return nil
	}

	var listFiles []string

	for _, file := range files {
		listFiles = append(listFiles, file.Name())
	}

	return listFiles
}

// Convierte la data de cada archivo
func parseData(dataLines *bufio.Scanner, id int) Email {
	var data Email
	var lastField string

	fieldSetters := map[string]func(string){
		"Message-ID:":                func(value string) { data.Message_ID = strings.TrimSpace(value) },
		"Date:":                      func(value string) { data.Date = strings.TrimSpace(value) },
		"From:":                      func(value string) { data.From = strings.TrimSpace(value) },
		"To:":                        func(value string) { data.To = strings.TrimSpace(value) },
		"Subject:":                   func(value string) { data.Subject = strings.TrimSpace(value) },
		"Cc:":                        func(value string) { data.Cc = strings.TrimSpace(value) },
		"Mime-Version:":              func(value string) { data.Mime_Version = strings.TrimSpace(value) },
		"Content-Type:":              func(value string) { data.Content_Type = strings.TrimSpace(value) },
		"Content-Transfer-Encoding:": func(value string) { data.Content_Transfer_Encoding = strings.TrimSpace(value) },
		"Bcc:":                       func(value string) { data.Bcc = strings.TrimSpace(value) },
		"X-From:":                    func(value string) { data.X_From = strings.TrimSpace(value) },
		"X-To:":                      func(value string) { data.X_To = strings.TrimSpace(value) },
		"X-cc:":                      func(value string) { data.X_cc = strings.TrimSpace(value) },
		"X-bcc:":                     func(value string) { data.X_bcc = strings.TrimSpace(value) },
		"X-Folder:":                  func(value string) { data.X_Folder = strings.TrimSpace(value) },
		"X-Origin:":                  func(value string) { data.X_Origin = strings.TrimSpace(value) },
		"X-FileName:":                func(value string) { data.X_FileName = strings.TrimSpace(value) },
	}

	for dataLines.Scan() {
		line := dataLines.Text()
		data.ID = id

		// Manejo de líneas continuadas (con espacios al inicio)
		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			if lastField != "" {
				// concatenar al último campo detectado
				switch lastField {
				case "To:":
					data.To += " " + strings.TrimSpace(line)
				case "Cc:":
					data.Cc += " " + strings.TrimSpace(line)
				case "Bcc:":
					data.Bcc += " " + strings.TrimSpace(line)
				default:
					data.Body += line + "\n"
				}
			}
			continue
		}

		// manejo normal de encabezados
		matched := false
		for prefix, setter := range fieldSetters {
			if strings.HasPrefix(line, prefix) {
				lastField = prefix
				setter(line[len(prefix):])
				matched = true
				break
			}
		}

		// Si no coincide con níngun encabezado, agregar al body
		if !matched {
			lastField = ""
			data.Body += line + "\n"
		}
	}

	return data
}

// Se configura solicitudes HTTP
var httpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        10,               //máximo número de conexiones inactivas a nivel global
		MaxIdleConnsPerHost: 5,                // máximo número de conexiones inactivas por host
		IdleConnTimeout:     30 * time.Second, // tiempo máximo que una conexión inactiva puede estar abierta
	},
}

// indexa data en batches
func batchIndexData(batch []Email) {
	user := os.Getenv("ZINC_USER")
	password := os.Getenv("ZINC_PASSWORD")
	auth := user + ":" + password
	encodeCreds := base64.StdEncoding.EncodeToString([]byte(auth))
	index := os.Getenv("ZINC_INDEX")
	zincHost := os.Getenv("ZINC_HOST")
	zincUrl := zincHost + "/api/" + index + "/_bulk"

	var bulkData []byte
	for _, data := range batch {
		// creo la metadata
		meta := []byte(fmt.Sprintf(`{ "index" : { "_index" : "%s" } }%s`, index, "\n"))
		// convierto los datos del email a JSON
		jsonData, _ := json.Marshal(data)
		// agrego la meta data y datos convertidos a bulkdata
		bulkData = append(bulkData, meta...)
		bulkData = append(bulkData, jsonData...)
		bulkData = append(bulkData, []byte("\n")...)
	}

	// creo la solicitud POST con los datos en bruto
	req, err := http.NewRequest("POST", zincUrl, bytes.NewBuffer(bulkData))
	if err != nil {
		log.Printf("Error creando solicitud: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+encodeCreds)

	// envío la solicitud http
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error enviando solicitud: %v", err)
		return
	}
	defer resp.Body.Close()
}

func IndexEmailData() {
	path := os.Getenv("DATA_PATH")
	semaphore := make(chan struct{}, 6)
	defer close(semaphore)

	var wg sync.WaitGroup
	count := 0

	fmt.Println("Indexando...")

	userList := folderList(path)

	for _, user := range userList {
		fmt.Println("Usuario:", user)
		folders := folderList(filepath.Join(path, user))

		for _, folder := range folders {
			fmt.Println("Carpeta:", folder)
			mailFiles := fileList(filepath.Join(path, user, folder))

			for _, mailFile := range mailFiles {
				wg.Add(1)
				count++
				filePath := filepath.Join(path, user, folder, mailFile)
				fmt.Printf("Procesando archivo: %s, Ruta: %s \n", mailFile, filePath)

				go func(filePath string, id int) {
					defer func() { <-semaphore }()
					semaphore <- struct{}{}

					emailFile, err := os.Open(filePath)
					if err != nil {
						log.Printf("Error abriendo archivo %s: %v", filePath, err)
						return
					}

					defer emailFile.Close()

					lines := bufio.NewScanner(emailFile)
					emailData := parseData(lines, id)
					batchIndexData([]Email{emailData})
				}(filePath, count)
			}
		}

	}

	log.Println("Indexación completada. Todos los archivos han sido procesados.")
	fmt.Println("Finished!!!!")
	wg.Wait()

}
