package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"kafka-consumer/internal/config"
	"kafka-consumer/internal/consumer"

	"github.com/stretchr/testify/assert"
)

func TestNewProcessor(t *testing.T) {
	cfg := &config.Config{
		TargetServiceURL: "http://localhost:8080/api/v1",
	}

	processor := NewProcessor(cfg)
	assert.NotNil(t, processor)
	assert.Equal(t, cfg.TargetServiceURL, processor.targetURL)
	assert.NotNil(t, processor.client)
}

func TestProcessor_Process(t *testing.T) {
	tests := []struct {
		name           string
		payload        *consumer.Payload
		serverResponse int
		expectedError  bool
	}{
		{
			name: "successful request",
			payload: &consumer.Payload{
				OrdemDeVenda: "test123",
				EtapaAtual:   "FATURADO",
			},
			serverResponse: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "server error",
			payload: &consumer.Payload{
				OrdemDeVenda: "test123",
				EtapaAtual:   "FATURADO",
			},
			serverResponse: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request method
				assert.Equal(t, http.MethodPatch, r.Method)

				// Verify content type
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

				// Verify request body
				var receivedPayload consumer.Payload
				err := json.NewDecoder(r.Body).Decode(&receivedPayload)
				assert.NoError(t, err)
				assert.Equal(t, tt.payload.OrdemDeVenda, receivedPayload.OrdemDeVenda)
				assert.Equal(t, tt.payload.EtapaAtual, receivedPayload.EtapaAtual)

				// Send response
				w.WriteHeader(tt.serverResponse)
			}))
			defer server.Close()

			// Create processor with test server URL
			processor := &Processor{
				targetURL: server.URL,
				client:    &http.Client{},
			}

			// Process payload
			err := processor.Process(tt.payload)

			// Verify result
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProcessor_Process_InvalidURL(t *testing.T) {
	processor := &Processor{
		targetURL: "http://invalid-url",
		client:    &http.Client{},
	}

	payload := &consumer.Payload{
		OrdemDeVenda: "test123",
		EtapaAtual:   "FATURADO",
	}

	err := processor.Process(payload)
	assert.Error(t, err)
}
