package emails

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func (e Email) GetEmailsData(queryData map[string]interface{}) (int, []Email, error) {

	// Configuración del host
	index := os.Getenv("ZINC_INDEX")
	zincHost := os.Getenv("ZINC_HOST")
	zincUrl := fmt.Sprintf("%s/api/%s/_search", zincHost, index)

	queryJSON, err := json.Marshal(queryData)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to serialize query: %w", err)
	}

	// Hacer la solicitud
	response, err := makeRequestZinc(http.MethodPost, zincUrl, queryJSON)
	if err != nil {
		return 0, nil, err
	}

	// Parsear la respuesta
	var apiResponse APIResponse
	if err := json.Unmarshal(response, &apiResponse); err != nil {
		return 0, nil, fmt.Errorf("error parsing API response: %w", err)
	}

	// Procesar los emails
	emails := []Email{}
	for _, hit := range apiResponse.Hits.Hits {
		email := Email{
			ID:         hit.Source.ID,
			Message_ID: hit.Source.Message_ID,
			Date:       hit.Source.Date,
			From:       hit.Source.From,
			To:         hit.Source.To,
			Subject:    hit.Source.Subject,
			Body:       hit.Source.Body,
		}
		emails = append(emails, email)
	}

	return apiResponse.Hits.Total.Value, emails, nil
}
