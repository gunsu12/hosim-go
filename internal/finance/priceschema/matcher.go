package priceschema

import (
	"context"
	"strings"
)

type ruleMatcher struct {
	repo       Repository
	unitReader ServiceUnitReader
}

// NewRuleMatcher membuat instance RuleMatcher yang mengimplementasikan algoritma override hirarki.
func NewRuleMatcher(repo Repository, unitReader ServiceUnitReader) RuleMatcher {
	return &ruleMatcher{
		repo:       repo,
		unitReader: unitReader,
	}
}

func (m *ruleMatcher) FindMatchingRule(ctx context.Context, q MatchRuleQuery) (MatchedRule, error) {
	q.PriceSchemaID = strings.TrimSpace(q.PriceSchemaID)
	q.ItemID = strings.TrimSpace(q.ItemID)
	q.DepartementID = strings.TrimSpace(q.DepartementID)
	q.ServiceUnitID = strings.TrimSpace(q.ServiceUnitID)

	if q.PriceSchemaID == "" {
		return MatchedRule{}, ErrNoMatchingRule
	}

	// Turunkan departement_id dari service unit jika belum diisi (PRD §3.1, §6.1)
	resolvedDeptID := q.DepartementID
	if resolvedDeptID == "" && q.ServiceUnitID != "" && m.unitReader != nil {
		if deptID, _, err := m.unitReader.GetServiceUnitInfo(ctx, q.ServiceUnitID); err == nil && deptID != nil {
			resolvedDeptID = *deptID
		}
	}

	rule, err := m.repo.FindBestMatchingRule(ctx, q.PriceSchemaID, q.ItemID, q.ServiceUnitID, resolvedDeptID)
	if err != nil {
		return MatchedRule{}, err
	}

	return MatchedRule{
		RuleID:                rule.ID,
		PriceSchemaID:         rule.PriceSchemaID,
		Level:                 rule.Level(),
		TariffAdjustmentType:  rule.TariffAdjustmentType,
		TariffAdjustmentValue: rule.TariffAdjustmentValue,
		PayByPatientPercent:   rule.PayByPatientPercent,
		PayByCustomerPercent:  rule.PayByCustomerPercent,
		TaxPercent:            rule.TaxPercent,
	}, nil
}
