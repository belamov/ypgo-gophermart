package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

type AccrualHTTPClient struct {
	client      *http.Client
	ratelimiter *rate.Limiter
	url         string
}

func (c *AccrualHTTPClient) GetAccrualForOrder(ctx context.Context, orderID int) (float64, error) {
	req, err := c.getRequest(orderID)
	if err != nil {
		log.Error().Err(err).Msg("received unexpected error while getting accrual info. cant initialize request")
		return 0, err
	}

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("received unexpected error while getting accrual info. cant make request")
		return 0, err
	}

	if resp.StatusCode == http.StatusNoContent {
		return 0, ErrOrderIsNotYetProceeded
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("received unexpected error while getting accrual info. cant read response body")
			return 0, err
		}
		log.Error().
			Err(err).
			Str("response", string(bodyBytes)).
			Int("status_code", resp.StatusCode).
			Msg("received unexpected response status from accrual service")

		return 0, errors.New("accrual service returned 500: " + string(bodyBytes))
	}

	accrualResp, err := c.parseResponse(resp)
	if err != nil {
		log.Error().Err(err).Msg("received unexpected error while getting accrual info. cant parse response")
		return 0, err
	}

	if accrualResp.isNotYetProceeded() {
		return 0, ErrOrderIsNotYetProceeded
	}

	if accrualResp.isInvalid() {
		return 0, ErrInvalidOrderForAccrual
	}

	if !accrualResp.isProcessed() {
		log.Error().Err(err).Msg("received unexpected error while getting accrual info. unknown order status")
		return 0, fmt.Errorf("unknown status of order: %s", accrualResp.Status)
	}

	return accrualResp.getAccrualAmount(), nil
}

func NewAccrualHTTPClient(client *http.Client, serviceURL string, maxRequestsPerSecond int) *AccrualHTTPClient {
	rl := rate.NewLimiter(rate.Every(time.Second), maxRequestsPerSecond)
	return &AccrualHTTPClient{
		ratelimiter: rl,
		client:      client,
		url:         serviceURL,
	}
}

func (c *AccrualHTTPClient) getRequest(orderID int) (*http.Request, error) {
	url := fmt.Sprintf("%s/api/orders/%v", c.url, orderID)
	request, err := http.NewRequest("GET", url, nil)
	request.Close = true
	return request, err
}

func (c *AccrualHTTPClient) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	err := c.ratelimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type accrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func (r *accrualResponse) isNotYetProceeded() bool {
	return r.Status == "REGISTERED" || r.Status == "PROCESSING"
}

func (r *accrualResponse) isInvalid() bool {
	return r.Status == "INVALID"
}

func (r *accrualResponse) isProcessed() bool {
	return r.Status == "PROCESSED"
}

func (r *accrualResponse) getAccrualAmount() float64 {
	return r.Accrual
}

func (c *AccrualHTTPClient) parseResponse(resp *http.Response) (accrualResponse, error) {
	var result accrualResponse
	err := json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return accrualResponse{}, err
	}

	return result, nil
}
