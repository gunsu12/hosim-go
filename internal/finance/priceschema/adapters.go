package priceschema

import (
	"context"
	"errors"
	"time"

	"hosim-go/internal/catalog/item"
	"hosim-go/internal/finance/priceplan"
	"hosim-go/internal/organization/department"
	"hosim-go/internal/organization/serviceunit"

	"github.com/shopspring/decimal"
)

// DepartmentServiceAdapter mengadaptasi department.Service menjadi DepartmentReader.
type DepartmentServiceAdapter struct {
	svc department.Service
}

func NewDepartmentServiceAdapter(svc department.Service) DepartmentReader {
	return &DepartmentServiceAdapter{svc: svc}
}

func (a *DepartmentServiceAdapter) IsActiveDepartment(ctx context.Context, id string) (bool, error) {
	dept, err := a.svc.GetDepartementByID(ctx, id)
	if err != nil {
		return false, err
	}
	return dept.IsActive, nil
}

// ServiceUnitServiceAdapter mengadaptasi serviceunit.Service menjadi ServiceUnitReader.
type ServiceUnitServiceAdapter struct {
	svc serviceunit.Service
}

func NewServiceUnitServiceAdapter(svc serviceunit.Service) ServiceUnitReader {
	return &ServiceUnitServiceAdapter{svc: svc}
}

func (a *ServiceUnitServiceAdapter) GetServiceUnitInfo(ctx context.Context, id string) (deptID *string, isActive bool, err error) {
	unit, err := a.svc.GetServiceUnitByID(ctx, id)
	if err != nil {
		return nil, false, err
	}
	return unit.DepartementID, unit.IsActive, nil
}

// ItemUseCaseAdapter mengadaptasi item.GetItemUseCase menjadi ItemReader.
type ItemUseCaseAdapter struct {
	uc item.GetItemUseCase
}

func NewItemUseCaseAdapter(uc item.GetItemUseCase) ItemReader {
	return &ItemUseCaseAdapter{uc: uc}
}

func (a *ItemUseCaseAdapter) IsActiveItem(ctx context.Context, id string) (bool, error) {
	itemDetail, err := a.uc.GetByID(ctx, id)
	if err != nil {
		return false, err
	}
	return itemDetail.IsActive, nil
}

// PricePlanReaderAdapter mengadaptasi priceplan.GetPricePlanUseCase menjadi PricePlanReader.
type PricePlanReaderAdapter struct {
	uc priceplan.GetPricePlanUseCase
}

func NewPricePlanReaderAdapter(uc priceplan.GetPricePlanUseCase) PricePlanReader {
	return &PricePlanReaderAdapter{uc: uc}
}

func (a *PricePlanReaderAdapter) PlanExists(ctx context.Context, id string) (bool, error) {
	plan, err := a.uc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, priceplan.ErrPricePlanNotFound) {
			return false, nil
		}
		return false, err
	}
	return plan != nil, nil
}

// PricePlanLookupAdapter mengadaptasi priceplan.LookupTariffUseCase menjadi TariffBaseLookup.
type PricePlanLookupAdapter struct {
	uc priceplan.LookupTariffUseCase
}

func NewPricePlanLookupAdapter(uc priceplan.LookupTariffUseCase) TariffBaseLookup {
	return &PricePlanLookupAdapter{uc: uc}
}

func (a *PricePlanLookupAdapter) LookupTariffBase(ctx context.Context, pricePlanID *string, tariffClassID, itemID string, serviceDate time.Time, isCito bool) (*TariffBaseResult, error) {
	dateStr := serviceDate.Format("2006-01-02")

	// Panggil lookup reguler untuk mendapatkan base price dan base amount per komponen
	baseReq := priceplan.LookupTariffRequest{
		PricePlanID:     pricePlanID,
		TariffClassID:   tariffClassID,
		ItemID:          itemID,
		TransactionDate: dateStr,
		IsCito:          false,
	}
	baseResp, err := a.uc.Execute(ctx, baseReq)
	if err != nil {
		return nil, err
	}

	totalBase := decimal.NewFromFloat(baseResp.TotalPrice)
	totalCito := totalBase

	components := make([]ComponentInput, len(baseResp.Components))
	for i, c := range baseResp.Components {
		baseAmt := decimal.NewFromFloat(c.Amount)
		components[i] = ComponentInput{
			ComponentID:   c.ComponentID,
			ComponentName: c.ComponentName,
			ComponentType: c.ComponentType,
			COACode:       c.COACode,
			BaseAmount:    baseAmt,
			CitoAmount:    baseAmt,
		}
	}

	// Jika CITO, panggil lookup dengan IsCito: true untuk mendapatkan total CITO dan CITO amount per komponen
	if isCito {
		citoReq := priceplan.LookupTariffRequest{
			PricePlanID:     pricePlanID,
			TariffClassID:   tariffClassID,
			ItemID:          itemID,
			TransactionDate: dateStr,
			IsCito:          true,
		}
		citoResp, err := a.uc.Execute(ctx, citoReq)
		if err == nil && citoResp != nil {
			totalCito = decimal.NewFromFloat(citoResp.TotalPrice)
			citoMap := make(map[string]decimal.Decimal, len(citoResp.Components))
			for _, c := range citoResp.Components {
				citoMap[c.ComponentID] = decimal.NewFromFloat(c.Amount)
			}
			for i, comp := range components {
				if amt, ok := citoMap[comp.ComponentID]; ok {
					components[i].CitoAmount = amt
				} else if i < len(citoResp.Components) {
					components[i].CitoAmount = decimal.NewFromFloat(citoResp.Components[i].Amount)
				}
			}
		}
	}

	return &TariffBaseResult{
		PricePlanID:    baseResp.PricePlanID,
		PricePlanCode:  baseResp.PricePlanCode,
		TotalBasePrice: totalBase,
		TotalCitoPrice: totalCito,
		Components:     components,
	}, nil
}
