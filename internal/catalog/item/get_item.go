package item

import (
	"context"
	"math"

	"hosim-go/pkg/response"
)

type GetItemUseCase interface {
	GetByID(ctx context.Context, id string) (*ItemDetailResponse, error)
	List(ctx context.Context, req QueryItemsRequest) ([]ItemSummaryResponse, response.PaginationMeta, error)
}

type getItemUseCase struct {
	repo Repository
}

func NewGetItemUseCase(repo Repository) GetItemUseCase {
	return &getItemUseCase{repo: repo}
}

func (uc *getItemUseCase) GetByID(ctx context.Context, id string) (*ItemDetailResponse, error) {
	item, err := uc.repo.FindItemByID(ctx, id)
	if err != nil {
		return nil, err
	}

	res := ToItemDetailResponse(item)
	return &res, nil
}

func (uc *getItemUseCase) List(ctx context.Context, req QueryItemsRequest) ([]ItemSummaryResponse, response.PaginationMeta, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	params := ListParams{
		Page:          page,
		Limit:         limit,
		Search:        req.Search,
		ItemType:      req.ItemType,
		CategoryID:    req.CategoryID,
		ProductLineID: req.ProductLineID,
		IsActive:      req.IsActive,
	}

	items, total, err := uc.repo.FindAllItems(ctx, params)
	if err != nil {
		return nil, response.PaginationMeta{}, err
	}

	res := make([]ItemSummaryResponse, len(items))
	for i := range items {
		res[i] = ToItemSummaryResponse(&items[i])
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	meta := response.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		TotalItems:  total,
		TotalPages:  totalPages,
	}

	return res, meta, nil
}
