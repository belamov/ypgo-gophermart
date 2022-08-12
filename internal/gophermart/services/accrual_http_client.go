package services

type AccrualHttpClient struct{}

func (a *AccrualHttpClient) GetAccrualForOrder(orderID int) (float64, error) {
	// TODO implement me
	panic("implement me")
}

func NewAccrualHttpClient() *AccrualHttpClient {
	return &AccrualHttpClient{}
}
