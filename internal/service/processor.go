package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"kafka-consumer/internal/config"
	"kafka-consumer/internal/consumer"
)

type Processor struct {
	targetURL string
	client    *http.Client
}

func NewProcessor(cfg *config.Config) *Processor {
	return &Processor{
		targetURL: cfg.TargetServiceURL,
		client:    &http.Client{},
	}
}

func (p *Processor) Process(payload *consumer.Payload) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPatch, p.targetURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	log.Printf("Sending PATCH request to %s with payload: %s", p.targetURL, string(jsonData))
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	log.Printf("Successfully processed order %s with status %s", payload.OrdemDeVenda, payload.EtapaAtual)
	return nil
} 