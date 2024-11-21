package api

type Response struct {
	Status  string      `json:"status,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func ResponseJson(data interface{}, message string) Response {
	var response Response
	if message == "" {
		response.Status = "SUCCESS"
		response.Data = data
	} else {
		response.Message = message
	}
	return response
}

type GetAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type LoanResponse struct {
	LoanID          int64   `json:"loan_id"`
	PrincipalAmount int64   `json:"principal_amount"`
	RateInPercent   float64 `json:"rate_in_percent"`
	RoiInPercent    float64 `json:"roi_in_percent"`
	Status          string  `json:"status"`
	AgreementLetter string  `json:"agreement_letter,omitempty"`
}
