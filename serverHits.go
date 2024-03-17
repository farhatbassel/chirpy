package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	db "farhatbassel/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	db             db.DB
}

func (config *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		config.fileserverHits += 1
		next.ServeHTTP(writer, request)
	})
}

func (config *apiConfig) displayNumberOfHits(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/html")
	writer.WriteHeader(http.StatusOK)

	htmlBody := fmt.Sprintf(`
    <html>

        <body>
            <h1>Welcome, Chirpy Admin</h1>
            <p>Chirpy has been visited %d times!</p>
        </body>

    </html>
    `, config.fileserverHits)

	writer.Write([]byte(htmlBody))
}

func (config *apiConfig) resetNumberOfHits(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	config.fileserverHits = 0
}

func (config *apiConfig) getChirps(writer http.ResponseWriter, request *http.Request) {
	chirps, err := config.db.GetChirps()

	if err != nil {
		sendJSONResponse(writer, err, http.StatusInternalServerError)
	}

	sendJSONResponse(writer, chirps, http.StatusOK)
}

func (config *apiConfig) createChirp(writer http.ResponseWriter, request *http.Request) {
	type chirpResponse struct {
		Valid       bool   `json:"valid"`
		Error       string `json:"error"`
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(request.Body)
	params := db.Chirp{}

	if err := decoder.Decode(&params); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody := chirpResponse{}

	if len(params.Body) > 140 {
		responseBody.Error = "Chirp is too long"
		sendJSONResponse(writer, responseBody, http.StatusBadRequest)
		return
	}

	responseBody.CleanedBody = cleanBody(params.Body)

	chirp, err := config.db.CreateChirp(responseBody.CleanedBody)

	if err != nil {
		responseBody.Error = fmt.Sprintf("%v", err)
		sendJSONResponse(writer, responseBody, http.StatusInternalServerError)
		return
	}

	sendJSONResponse(writer, chirp, http.StatusCreated)
}

func sendJSONResponse(writer http.ResponseWriter, data interface{}, status int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

	if err := json.NewEncoder(writer).Encode(data); err != nil {
		http.Error(writer, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}

func cleanBody(input string) string {
	words := strings.Fields(input)
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	output := make([]string, 0)
	for _, word := range words {
		lowerWord := strings.ToLower(word)
		if _, isBadWord := badWords[lowerWord]; isBadWord {
			output = append(output, "****")
		} else {
			output = append(output, word)
		}
	}

	return strings.Join(output, " ")
}
