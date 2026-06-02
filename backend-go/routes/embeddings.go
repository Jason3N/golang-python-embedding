package routes

import (
	"api-handler/models"
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type embeddingModelResponse struct {
	Embedding []float64 `json:"embedding"`
}

func (h *Handler) PostEmbedding(w http.ResponseWriter, r *http.Request) {
	// 1. Prevents any method being used here besides a POST
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed, only POST", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	// You want to create a request now per EmbeddingRequest, which is just a JSON
	var req models.EmbeddingRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// The caller did not send valid JSON that matches models.EmbeddingRequest.
		http.Error(w, "Invalid JSON Body", http.StatusBadRequest)
		return
	}

	// if req is Empty, then throw an error
	if req.Text == "" {
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		// Turning a valid Go struct into JSON should rarely fail, so treat it as a server error.
		http.Error(w, "Failed to encode model request", http.StatusInternalServerError)
		return
	}

	// After we verified that req is okay, then we can go ahead and ping the EMBEDDING_MODEL_URL
	modelURL := os.Getenv("EMBEDDING_MODEL_URL")
	if modelURL == "" {
		http.Error(w, "EMBEDDING_MODEL_URL is not set", http.StatusInternalServerError)
		return
	}

	// Create a new request to the embedding-model_url
	modelReq, err := http.NewRequest(http.MethodPost, modelURL, bytes.NewBuffer(jsonData))
	if err != nil {
		// The outbound request could not be built, usually because the URL is invalid.
		http.Error(w, "Failed to create model request", http.StatusInternalServerError)
		return
	}
	modelReq.Header.Set("Content-Type", "application/json")

	// modelResp should be the response from the Embedding endpoint hosted by Python
	modelResp, err := http.DefaultClient.Do(modelReq)
	if err != nil {
		// Go could not reach the Python embedding service or the request failed at the network layer.
		http.Error(w, "Failed to call embedding model", http.StatusBadGateway)
		return
	}

	defer modelResp.Body.Close()
	if modelResp.StatusCode != http.StatusOK {
		http.Error(w, "Embedding model returned an error", http.StatusBadGateway)
		return
	}

	// embeddingResp is the model that we want it to return which is a list of floating integers
	var embeddingResp models.EmbeddingResponse
	err = json.NewDecoder(modelResp.Body).Decode(&embeddingResp)
	if err != nil {
		// The Python service responded, but its body was not valid JSON in the shape we expected.
		http.Error(w, "Failed to decode embedding response", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(embeddingResp)
	if err != nil {
		// Go failed while writing the final JSON response back to the original caller.
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetAllEmbeddings(w http.ResponseWriter, r *http.Request) {
	// 1. Prevents any method being used here besides a GET
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method not allowed, only GET", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	rows, err := h.DB.Query("SELECT id, content, embedding FROM chunks")
	if err != nil {
		http.Error(w, "Could not get from database", http.StatusInternalServerError)
	}

	var embeddingResp []models.StoredEmbeddings
	for rows.Next() {
		var embedding models.StoredEmbeddings
		err := rows.Scan(&embedding.ID,
			&embedding.Content,
			&embedding.Embedding)
		if err != nil {
			http.Error(w, "ERROR", http.StatusInternalServerError)
		}
		embeddingResp = append(embeddingResp, embedding)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(embeddingResp)
	if err != nil {
		// Go failed while writing the final JSON response back to the original caller.
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
