package main

import (
	datasetindex "backend/dataset-index"
	"backend/emails"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

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
		"page_size": 0,
		"query": map[string]string{
			"search_type": "matchall",
		},
	}

	total, _, err := e.GetEmailsData(query)
	if err != nil {
		log.Printf("Error obteniendo los hits: %v", err)
	}

	if total == 0 {
		log.Println("No se encontraron datos. Iniciando subida de datos ...")
		datasetindex.IndexEmailData()
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

	// Usar un canal para capturar señales y detener el servidor
	go func() {
		log.Println("Starting server on :3000...")
		if err := http.ListenAndServe(os.Getenv("PORT"), r); err != nil {
			log.Fatalf("Error iniciando servidor: %v", err)
		}
	}()

	// Capturar señales para detener el servidor y generar perfiles
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c // Espera la señal

	// Memory Profiling
	memFile, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer memFile.Close()
	pprof.WriteHeapProfile(memFile)
}
