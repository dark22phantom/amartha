package http

type Handler struct {
	ucLoan ucLoanInterface
}

func New(ucLoan ucLoanInterface) *Handler {
	return &Handler{
		ucLoan: ucLoan,
	}
}
