package status

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
)

type statusUC interface {
	CheckIn(ctx context.Context, body model.CheckInRequest) (model.CheckInResponse, error)
	CheckOut(ctx context.Context, body model.CheckOutRequest) (model.CheckOutResponse, error)
	Report(ctx context.Context, placeId string) (model.ReportResponse, error)
}

type DeliveryStatus struct {
	status statusUC
}

func NewStatus(status statusUC) *DeliveryStatus {
	return &DeliveryStatus{
		status: status,
	}
}
