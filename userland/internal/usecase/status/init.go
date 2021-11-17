package status

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/entity"
)

type statusRepo interface {
	CheckIn(ctx context.Context, req model.CheckInRequest) (entity.Status, entity.Profile, entity.Place, error)
	CheckOut(ctx context.Context, req model.CheckOutRequest) (entity.Status, error)
	Report(ctx context.Context, placeId string) (entity.Place, []entity.Report, error)
}

type UsecaseStatus struct {
	status statusRepo
}

func New(repo statusRepo) *UsecaseStatus {
	return &UsecaseStatus{
		status: repo,
	}
}
