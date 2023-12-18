package service

type (
	IInboxMarketingFpt interface{}
	InboxMarketingFpt  struct{}
)

func NewInboxMarketingFpt() IInboxMarketingFpt {
	return &InboxMarketingFpt{}
}
