package http

import (
	"amartha/internal/pkg/helper"
	modelApi "amartha/model/api"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func (h *Handler) GetAccessToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// parse request
	var reqBody modelApi.GetAccessToken
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, err.Error()))
		return
	}

	// validate request
	if err := modelApi.Validate.Struct(reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, err.Error()))
		return
	}

	// get access token
	response, err := h.ucLoan.GetAccessToken(r.Context(), reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(modelApi.ResponseJson(response, ""))
	return
}

func (h *Handler) LoanSubmit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// parse request
	var reqBody modelApi.LoanSubmit
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, err.Error()))
		return
	}

	// validate request
	if err := modelApi.Validate.Struct(reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, err.Error()))
		return
	}

	response, err := h.ucLoan.LoanSubmit(r.Context(), reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(modelApi.ResponseJson(response, ""))
	return
}

func (h *Handler) LoanApproval(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// set limit file
	err := r.ParseMultipartForm(5)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, err.Error()))
		return
	}

	loanID := r.FormValue("loan_id")
	if loanID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, "loan_id is required"))
		return
	}
	loanIDInt, err := strconv.ParseInt(loanID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, err.Error()))
		return
	}

	file, header, err := r.FormFile("validator_photo")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, "validator_photo is required"))
		return
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, err.Error()))
		return
	}

	fileName := header.Filename
	if err := helper.ValidateFileExtension(fileName); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, err.Error()))
		return
	}

	response, err := h.ucLoan.LoanApproval(r.Context(), modelApi.LoanApproval{
		LoanID:         loanIDInt,
		ValidatorPhoto: fileContent,
	}, fileName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(modelApi.ResponseJson(response, ""))
	return
}

func (h *Handler) LoanInvestment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// parse request
	var reqBody modelApi.LoanInvestment
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, err.Error()))
		return
	}

	// validate request
	if err := modelApi.Validate.Struct(reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, err.Error()))
		return
	}

	response, err := h.ucLoan.LoanInvestment(r.Context(), reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(modelApi.ResponseJson(response, ""))
	return
}

func (h *Handler) LoanDisbursement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// set limit file
	err := r.ParseMultipartForm(5)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, err.Error()))
		return
	}

	loanID := r.FormValue("loan_id")
	if loanID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, "loan_id is required"))
		return
	}
	loanIDInt, err := strconv.ParseInt(loanID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, err.Error()))
		return
	}

	file, header, err := r.FormFile("agreement_letter")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, "agreement_letter is required"))
		return
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, err.Error()))
		return
	}

	fileName := header.Filename
	if err := helper.ValidateFileExtension(fileName); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, err.Error()))
		return
	}

	response, err := h.ucLoan.LoanDisbursement(r.Context(), modelApi.LoanDisbursement{
		LoanID:          loanIDInt,
		AgreementLetter: fileContent,
	}, fileName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(modelApi.ResponseJson(nil, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(modelApi.ResponseJson(response, ""))
	return
}
