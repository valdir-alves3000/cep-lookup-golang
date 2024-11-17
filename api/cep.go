package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type Address struct {
	CEP        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
}

type Response struct {
	Address   *Address `json:"address"`
	API       string   `json:"api"`
	Error     string   `json:"error,omitempty"`
	TimeoutMS int64   `json:"timeout_ms,omitempty"`
}

func handlerAddress(ctx context.Context, url string) (*Address, string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	var address Address
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		return nil, "", err
	}
	cleanURL := strings.TrimPrefix(strings.TrimPrefix(url, "http://"), "https://")
	parts := strings.SplitN(cleanURL, "/", 2)
	baseURL := parts[0]

	return &address, baseURL, nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Configurar CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Extrair CEP da query string
	cep := r.URL.Query().Get("cep")
	if cep == "" {
		http.Error(w, "CEP não fornecido", http.StatusBadRequest)
		return
	}

	brasilAPIURL := "https://brasilapi.com.br/api/cep/v1/" + cep
	viaCEPURL := "http://viacep.com.br/ws/" + cep + "/json/"
	maxTime := 1.0

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(maxTime)*time.Second)
	defer cancel()

	resultChannelAddress := make(chan struct {
		address *Address
		api     string
		err     error
	})

	go func() {
		address, api, err := handlerAddress(ctx, brasilAPIURL)
		resultChannelAddress <- struct {
			address *Address
			api     string
			err     error
		}{address, api, err}
	}()

	go func() {
		address, api, err := handlerAddress(ctx, viaCEPURL)
		resultChannelAddress <- struct {
			address *Address
			api     string
			err     error
		}{address, api, err}
	}()

	w.Header().Set("Content-Type", "application/json")

	select {
	case res := <-resultChannelAddress:
		response := Response{}
		if res.err != nil {
			response.Error = res.err.Error()
			json.NewEncoder(w).Encode(response)
			return
		}
		response.Address = res.address
		response.API = res.api
		json.NewEncoder(w).Encode(response)

	case <-ctx.Done():
		response := Response{
			Error:     "Timeout - Nenhuma das APIs respondeu",
			TimeoutMS: int64(maxTime * 1000),
		}
		json.NewEncoder(w).Encode(response)
	}
}
