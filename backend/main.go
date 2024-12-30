package main

import (
	datasetindex "backend/dataset-index"
	"backend/emails"
	"log"
	"net/http"
	"os"
	"runtime/pprof"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func loadEnvVars() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	// CPU Profiling
	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer cpuFile.Close()
	pprof.StartCPUProfile(cpuFile)
	defer pprof.StopCPUProfile()

	// Cargar las variables de entorno
	loadEnvVars()

	// Crear instancia de Email
	e := emails.Email{}

	// Verificar si el índice tiene datos
	query := map[string]interface{}{
		"from":      0,
		"page_size": 0, // Solo interesa obtener el total, no los datos
		"query": map[string]string{
			"search_type": "matchall",
		},
	}

	total, _, err := e.GetEmailsData(query)
	if err != nil {
		log.Fatalf("Error obteniendo los hits: %v", err)
	}

	if total == 0 {
		log.Println("No se encontraron datos. Iniciando subida de datos ...")
		datasetindex.IndexAndCreateJson()
	} else {
		log.Println("Los datos ya se encuentran cargados en la base de datos.")
		log.Printf("Total datos cargados: %d", total)
	}

	// Configurar servidor
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedHeaders: []string{"Content-Type"},
		AllowedMethods: []string{"GET", "OPTIONS"},
	}))
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	r.Mount("/emails", emails.EmailsRoutes())

	// Configurar el puerto del servidor
	port := os.Getenv("PORT")
	log.Println("Starting server on :3000...")
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}

	// Memory Profiling
	memFile, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer memFile.Close()
	pprof.WriteHeapProfile(memFile)
}
