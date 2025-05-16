package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Payload struct {
	OrdemDeVenda string `json:"ordemDeVenda"`
	EtapaAtual   string `json:"etapaAtual"`
}

func main() {
	http.HandleFunc("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload Payload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		log.Printf("Received PATCH request for order %s with status %s", payload.OrdemDeVenda, payload.EtapaAtual)
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Mock service listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
} 