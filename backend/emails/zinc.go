package emails

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func makeRequestZinc(method, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(os.Getenv("ZINC_USER"), os.Getenv("ZINC_PASSWORD"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	// Leer la respuesta
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo la respuesta: %w", err)
	}

	// Manejar diferentes códigos de estado
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(
			"solicitud fallida con código %d, respuesta: %s",
			resp.StatusCode, string(responseBody),
		)
	}

	return responseBody, nil
}
