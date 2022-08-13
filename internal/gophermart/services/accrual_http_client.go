package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type AccrualHttpClient struct {
	client      *http.Client
	ratelimiter *rate.Limiter
	url         string
}

func (c *AccrualHttpClient) GetAccrualForOrder(ctx context.Context, orderID int) (float64, error) {
	req, err := c.getRequest(orderID)
	if err != nil {
		return 0, err
	}

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("accrual service returned 500")
	}

	accrualResp, err := c.parseResponse(resp)
	if err != nil {
		return 0, err
	}

	if accrualResp.isNotYetProceeded() {
		return 0, ErrOrderIsNotYetProceeded
	}

	if accrualResp.isInvalid() {
		return 0, ErrInvalidOrderForAccrual
	}

	if !accrualResp.isProcessed() {
		return 0, fmt.Errorf("unknown status of order: %s", accrualResp.Status)
	}

	return accrualResp.getAccrualAmount(), nil
}

func NewAccrualHttpClient(client *http.Client, serviceURL string, maxRequestsPerSecond int) *AccrualHttpClient {
	rl := rate.NewLimiter(rate.Every(time.Second), maxRequestsPerSecond)
	return &AccrualHttpClient{
		ratelimiter: rl,
		client:      client,
		url:         serviceURL,
	}
}

func (c *AccrualHttpClient) getRequest(orderID int) (*http.Request, error) {
	url := fmt.Sprintf("%s/%v", c.url, orderID)
	return http.NewRequest("GET", url, nil)
}

func (c *AccrualHttpClient) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
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

func (c *AccrualHttpClient) parseResponse(resp *http.Response) (accrualResponse, error) {
	var result accrualResponse
	err := json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return accrualResponse{}, err
	}

	return result, nil
}
