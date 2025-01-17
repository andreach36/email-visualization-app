package emails

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Función para listar todos los emails
func (e Email) GetAllEmails(w http.ResponseWriter, r *http.Request) {

	// Parsear parámetro de paginación
	pageNum, err := strconv.Atoi(r.URL.Query().Get("page_num"))
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSize := 10

	from := (pageNum - 1) * pageSize

	query := map[string]interface{}{
		"from":      from,
		"page_size": pageSize,
		"query": map[string]string{
			"search_type": "matchall",
		},
		"sort_fields": []string{"-Date"},
	}

	//Obtener emails data
	total, emails, err := e.GetEmailsData(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Crear una respuesta paginada
	totalPages := (total + pageSize - 1) / pageSize
	responsePayload := map[string]interface{}{
		"data": emails,
		"pagination": map[string]interface{}{
			"page_num":    pageNum,
			"page_size":   pageSize,
			"total_pages": totalPages,
			"total_data":  total,
		},
	}

	// Enviar la respuesta al cliente
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responsePayload)
}

// Función para manejar la búsqueda por palabra clave
func (e Email) SearchEmails(w http.ResponseWriter, r *http.Request) {

	// Obtener los parámetros de búsqueda
	keyword := r.URL.Query().Get("q")
	if keyword == "" {
		http.Error(w, "Missing 'q' search parameter", http.StatusBadRequest)
		return
	}

	// Parsear los parámetros de paginación
	pageNum, err := strconv.Atoi(r.URL.Query().Get("page_num"))
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSize := 10

	from := (pageNum - 1) * pageSize

	query := map[string]interface{}{
		"search_type": "match",
		"from":        from,
		"page_size":   pageSize,
		"query": map[string]interface{}{
			"term":  keyword,
			"field": "_all",
		},
		"sort_fields": []string{"-Date"},
	}

	// Obtener emails data
	total, emails, err := e.GetEmailsData(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Crear una respuesta paginada
	totalPages := (total + pageSize - 1) / pageSize
	responsePayload := map[string]interface{}{
		"data": emails,
		"pagination": map[string]interface{}{
			"page_num":    pageNum,
			"page_size":   pageSize,
			"total_pages": totalPages,
			"total_data":  total,
		},
	}

	// Enviar la respuesta al cliente
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responsePayload)
}
