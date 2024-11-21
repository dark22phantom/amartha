package http

import (
	"github.com/go-chi/chi/v5"
)

func (h *Handler) RouteHandler(r chi.Router) {
	// public
	r.Post("/token/request", h.GetAccessToken)

	// private
	r.Group(func(r chi.Router) {
		r.Use(JWTAuthMiddleware)
		r.Post("/loan/submit", h.LoanSubmit)
		r.Post("/loan/approval", h.LoanApproval)
		r.Post("/loan/investment", h.LoanInvestment)
		r.Post("/loan/disbursement", h.LoanDisbursement)
	})
}
