package services

type BalanceProcessorInterface interface {
	RegisterWithdraw(orderID int, userID int, amount float64) error
}

type BalanceProcessor struct{}

func (b *BalanceProcessor) RegisterWithdraw(orderID int, userID int, amount float64) error {
	// TODO implement me
	panic("implement me")
}

func NewBalanceProcessor() *BalanceProcessor {
	return &BalanceProcessor{}
}
