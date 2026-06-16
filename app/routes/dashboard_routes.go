package routes

import (
	"encoding/json"
	repo "codigo/app/repository"
	"net/http"
)

func RotasDashboard() {

	http.HandleFunc("/api/lojas", func(w http.ResponseWriter, r *http.Request) {

		// Permite apenas GET
		if r.Method != http.MethodGet {

			// Retorna erro 405
			http.Error(
				w,
				"Método não permitido",
				http.StatusMethodNotAllowed,
			)

			return
		}

		// Busca lojas no banco
		lojas, err := repo.Read_lojas()

		// Tratamento de erro
		if err != nil {

			jsonError(
				w,
				"Erro ao buscar lojas",
				http.StatusInternalServerError,
			)

			return
		}

		// Retorna lista em JSO
		jsonResponse(w, lojas)
	})
}

func jsonError(w http.ResponseWriter, msg string, code int) {

	// Header JSON
	w.Header().Set("Content-Type", "application/json")

	// Status HTTP
	w.WriteHeader(code)

	// Corpo da resposta
	json.NewEncoder(w).Encode(
		map[string]string{
			"error": msg,
		},
	)
}

func jsonResponse(w http.ResponseWriter, data interface{}) {

	// Define header da resposta
	w.Header().Set("Content-Type", "application/json")

	// Converte struct/map/slice para JSON
	json.NewEncoder(w).Encode(data)
}