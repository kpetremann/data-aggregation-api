package netbox

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"

	"github.com/criteo/data-aggregation-api/internal/config"
)

type NetboxResponse[R any] struct {
	Next    string `json:"next"`
	Results []*R   `json:"results" validate:"dive"`
	Count   int    `json:"count"`
}

// NewGetRequest returns a prepared Netbox request with the authentication set.
func NewGetRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Token "+config.Cfg.NetBox.APIKey)

	return req, err
}

// Get fetches a Netbox endpoint.
func Get[R any](endpoint string, out *NetboxResponse[R], params url.Values) error {
	const endpointKey = "endpoint"
	client := http.Client{Timeout: 10 * 60 * time.Second}

	params.Set("limit", "0")
	params.Set("ordering", "id")

	baseURL, err := url.JoinPath(config.Cfg.NetBox.URL, endpoint)
	if err != nil {
		return fmt.Errorf("failed to assemble URL: %w", err)
	}
	url := baseURL + "?" + params.Encode()
	log.Info().Str(endpointKey, endpoint).Msgf("Get %s", url)

	// Get all pages
	for url != "" {
		req, err := NewGetRequest(url)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		data, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to query netbox: %w", err)
		}
		defer func() {
			if err := data.Body.Close(); err != nil {
				log.Error().Err(err).Msg("issue to close netbox query")
			}
		}()

		if data.StatusCode < 200 || data.StatusCode > 299 {
			return fmt.Errorf("netbox returned HTTP error code: '%s'", data.Status)
		}

		var buffer NetboxResponse[R]
		err = json.NewDecoder(data.Body).Decode(&buffer)
		if err != nil {
			return fmt.Errorf("failed to decode netbox response: %w", err)
		}

		out.Results = append(out.Results, buffer.Results...)

		// Print paging status
		offset := data.Request.URL.Query().Get("offset")
		if len(offset) == 0 {
			offset = "0"
		}
		log.Debug().Str(endpointKey, endpoint).Msgf("status: %s/%d", offset, buffer.Count)

		url = buffer.Next
		out.Count = buffer.Count
	}

	// Validate
	validate := validator.New()
	if err := validate.Struct(out); err != nil {
		return err
	}

	return nil
}
