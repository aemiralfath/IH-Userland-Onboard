package model

import "github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/entity"

type CheckInRequest struct {
	Profile entity.Profile `json:"user"`
	Place   entity.Place   `json:"place"`
}

type CheckInResponse struct {
	Status  entity.Status  `json:"status"`
	Profile entity.Profile `json:"user"`
	Place   entity.Place   `json:"place"`
}

type CheckOutRequest struct {
	StatusID string `json:"status_id"`
}

type CheckOutResponse struct {
	Status entity.Status `json:"status"`
}

type ReportResponse struct {
	Place   entity.Place    `json:"place"`
	Reports []entity.Report `json:"reports"`
}
