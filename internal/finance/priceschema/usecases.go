package priceschema

import (
	"context"
	"strings"
	"time"

	"hosim-go/internal/audit"
	"hosim-go/internal/database"
	"hosim-go/pkg/response"

	"github.com/shopspring/decimal"
)

// OperatorContext membawa data actor yang terotentikasi dari middleware JWT.
type OperatorContext struct {
	UserID        string
	Username      string
	Role          string
	CorrelationID string
}

// =========================================================================
// 1. Create Price Schema Use Case
// =========================================================================

type CreatePriceSchemaUseCase interface {
	Execute(ctx context.Context, req CreatePriceSchemaRequest, actor OperatorContext) (*PriceSchemaResponse, error)
}

type createPriceSchemaUseCase struct {
	repo       Repository
	tx         database.TransactionManager
	audit      audit.AuditWriter
	deptReader DepartmentReader
	unitReader ServiceUnitReader
	itemReader ItemReader
}

func NewCreatePriceSchemaUseCase(
	repo Repository,
	tx database.TransactionManager,
	audit audit.AuditWriter,
	deptReader DepartmentReader,
	unitReader ServiceUnitReader,
	itemReader ItemReader,
) CreatePriceSchemaUseCase {
	return &createPriceSchemaUseCase{
		repo:       repo,
		tx:         tx,
		audit:      audit,
		deptReader: deptReader,
		unitReader: unitReader,
		itemReader: itemReader,
	}
}

func (uc *createPriceSchemaUseCase) Execute(ctx context.Context, req CreatePriceSchemaRequest, actor OperatorContext) (*PriceSchemaResponse, error) {
	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)

	if req.Code == "" {
		return nil, ErrSchemaCodeRequired
	}
	if req.Name == "" {
		return nil, ErrSchemaNameRequired
	}

	// Cek keunikan kode
	if existing, err := uc.repo.FindSchemaByCode(ctx, req.Code); err == nil && existing != nil {
		return nil, ErrSchemaCodeExists
	}

	schema := &PriceSchema{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
		CreatedBy:   actor.Username,
		UpdatedBy:   actor.Username,
	}

	// Konversi dan validasi aturan awal jika ada
	var rules []PriceSchemaRule
	hasDefaultRule := false
	seenScopes := make(map[string]bool)

	for _, rReq := range req.Rules {
		medAdjType := rReq.MedicationAdjustmentType
		if medAdjType == "" {
			medAdjType = AdjustmentTypeNone
		}
		tariffAdjType := rReq.TariffAdjustmentType
		if tariffAdjType == "" {
			tariffAdjType = AdjustmentTypeNone
		}

		rule := PriceSchemaRule{
			DepartementID:             cleanStringPtr(rReq.DepartementID),
			ServiceUnitID:             cleanStringPtr(rReq.ServiceUnitID),
			ItemID:                    cleanStringPtr(rReq.ItemID),
			TariffAdjustmentType:      tariffAdjType,
			TariffAdjustmentValue:     rReq.TariffAdjustmentValue,
			MedicationAdjustmentType:  medAdjType,
			MedicationAdjustmentValue: rReq.MedicationAdjustmentValue,
			PayByPatientPercent:       rReq.PayByPatientPercent,
			PayByCustomerPercent:      rReq.PayByCustomerPercent,
			TaxPercent:                rReq.TaxPercent,
			CreatedBy:                 actor.Username,
			UpdatedBy:                 actor.Username,
		}

		if err := rule.Validate(); err != nil {
			return nil, err
		}

		// Validasi scope ke modul pemilik (organization / catalog)
		if rule.DepartementID != nil && uc.deptReader != nil {
			active, err := uc.deptReader.IsActiveDepartment(ctx, *rule.DepartementID)
			if err != nil || !active {
				return nil, ErrDepartmentInactiveOrNotFound
			}
		}
		if rule.ServiceUnitID != nil && uc.unitReader != nil {
			_, active, err := uc.unitReader.GetServiceUnitInfo(ctx, *rule.ServiceUnitID)
			if err != nil || !active {
				return nil, ErrServiceUnitInactiveOrNotFound
			}
		}
		if rule.ItemID != nil && uc.itemReader != nil {
			active, err := uc.itemReader.IsActiveItem(ctx, *rule.ItemID)
			if err != nil || !active {
				return nil, ErrItemInactiveOrNotFound
			}
		}

		if rule.Level() == MatchLevelSchemaDefault {
			if hasDefaultRule {
				return nil, ErrMultipleDefaultRules
			}
			hasDefaultRule = true
		}

		scopeKey := stringFromPtr(rule.DepartementID) + "|" + stringFromPtr(rule.ServiceUnitID) + "|" + stringFromPtr(rule.ItemID)
		if seenScopes[scopeKey] {
			return nil, ErrDuplicateRuleScope
		}
		seenScopes[scopeKey] = true
		rules = append(rules, rule)
	}

	// Invariant PRD §11 Q3: Setiap Price Schema WAJIB memiliki tepat 1 rule SCHEMA_DEFAULT
	if !hasDefaultRule {
		defaultRule := PriceSchemaRule{
			TariffAdjustmentType:      AdjustmentTypeNone,
			TariffAdjustmentValue:     decimal.Zero,
			MedicationAdjustmentType:  AdjustmentTypeNone,
			MedicationAdjustmentValue: decimal.Zero,
			PayByPatientPercent:       decimal.Zero,
			PayByCustomerPercent:      decimal.NewFromInt(100),
			TaxPercent:                decimal.Zero,
			CreatedBy:                 actor.Username,
			UpdatedBy:                 actor.Username,
		}
		rules = append(rules, defaultRule)
	}

	// Eksekusi atomik: Header + Rules + Audit dalam 1 transaksi
	err := uc.tx.WithTx(ctx, func(txCtx context.Context) error {
		if err := uc.repo.CreateSchema(txCtx, schema); err != nil {
			return err
		}

		for i := range rules {
			rules[i].PriceSchemaID = schema.ID
			if err := uc.repo.CreateRule(txCtx, &rules[i]); err != nil {
				return err
			}
		}

		schema.Rules = rules

		// Catat jejak audit append-only
		return uc.audit.Write(txCtx, audit.AuditEntry{
			ActorID:       actor.UserID,
			ActorUsername: actor.Username,
			ActorRole:     actor.Role,
			Action:        "CREATE",
			ResourceType:  "price_schema",
			ResourceID:    schema.ID,
			After:         schema,
			CorrelationID: actor.CorrelationID,
		})
	})

	if err != nil {
		return nil, err
	}

	resp := ToPriceSchemaResponse(schema)
	return &resp, nil
}

// =========================================================================
// 2. Update Price Schema Use Case
// =========================================================================

type UpdatePriceSchemaUseCase interface {
	Execute(ctx context.Context, id string, req UpdatePriceSchemaRequest, actor OperatorContext) (*PriceSchemaResponse, error)
}

type updatePriceSchemaUseCase struct {
	repo  Repository
	tx    database.TransactionManager
	audit audit.AuditWriter
}

func NewUpdatePriceSchemaUseCase(repo Repository, tx database.TransactionManager, audit audit.AuditWriter) UpdatePriceSchemaUseCase {
	return &updatePriceSchemaUseCase{repo: repo, tx: tx, audit: audit}
}

func (uc *updatePriceSchemaUseCase) Execute(ctx context.Context, id string, req UpdatePriceSchemaRequest, actor OperatorContext) (*PriceSchemaResponse, error) {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, ErrSchemaNameRequired
	}

	schema, err := uc.repo.FindSchemaByID(ctx, id)
	if err != nil {
		return nil, err
	}

	beforeSnapshot := *schema

	schema.Name = req.Name
	schema.Description = req.Description
	if req.IsActive != nil {
		schema.IsActive = *req.IsActive
	}
	schema.UpdatedBy = actor.Username
	schema.UpdatedAt = time.Now().UTC()

	err = uc.tx.WithTx(ctx, func(txCtx context.Context) error {
		if err := uc.repo.UpdateSchema(txCtx, schema); err != nil {
			return err
		}

		return uc.audit.Write(txCtx, audit.AuditEntry{
			ActorID:       actor.UserID,
			ActorUsername: actor.Username,
			ActorRole:     actor.Role,
			Action:        "UPDATE",
			ResourceType:  "price_schema",
			ResourceID:    schema.ID,
			Before:        beforeSnapshot,
			After:         schema,
			CorrelationID: actor.CorrelationID,
		})
	})

	if err != nil {
		return nil, err
	}

	resp := ToPriceSchemaResponse(schema)
	return &resp, nil
}

// =========================================================================
// 3. Deactivate Price Schema Use Case
// =========================================================================

type DeactivatePriceSchemaUseCase interface {
	Execute(ctx context.Context, id string, actor OperatorContext) (*PriceSchemaResponse, error)
}

type deactivatePriceSchemaUseCase struct {
	repo  Repository
	tx    database.TransactionManager
	audit audit.AuditWriter
}

func NewDeactivatePriceSchemaUseCase(repo Repository, tx database.TransactionManager, audit audit.AuditWriter) DeactivatePriceSchemaUseCase {
	return &deactivatePriceSchemaUseCase{repo: repo, tx: tx, audit: audit}
}

func (uc *deactivatePriceSchemaUseCase) Execute(ctx context.Context, id string, actor OperatorContext) (*PriceSchemaResponse, error) {
	schema, err := uc.repo.FindSchemaByID(ctx, id)
	if err != nil {
		return nil, err
	}

	beforeSnapshot := *schema

	schema.IsActive = false
	schema.UpdatedBy = actor.Username
	schema.UpdatedAt = time.Now().UTC()

	err = uc.tx.WithTx(ctx, func(txCtx context.Context) error {
		if err := uc.repo.UpdateSchema(txCtx, schema); err != nil {
			return err
		}

		return uc.audit.Write(txCtx, audit.AuditEntry{
			ActorID:       actor.UserID,
			ActorUsername: actor.Username,
			ActorRole:     actor.Role,
			Action:        "DEACTIVATE",
			ResourceType:  "price_schema",
			ResourceID:    schema.ID,
			Before:        beforeSnapshot,
			After:         schema,
			CorrelationID: actor.CorrelationID,
		})
	})

	if err != nil {
		return nil, err
	}

	resp := ToPriceSchemaResponse(schema)
	return &resp, nil
}

// =========================================================================
// 4. Get & List Price Schema Use Cases (Query Path)
// =========================================================================

type GetPriceSchemaUseCase interface {
	Execute(ctx context.Context, id string) (*PriceSchemaResponse, error)
}

type getPriceSchemaUseCase struct {
	repo Repository
}

func NewGetPriceSchemaUseCase(repo Repository) GetPriceSchemaUseCase {
	return &getPriceSchemaUseCase{repo: repo}
}

func (uc *getPriceSchemaUseCase) Execute(ctx context.Context, id string) (*PriceSchemaResponse, error) {
	schema, err := uc.repo.FindSchemaByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := ToPriceSchemaResponse(schema)
	return &resp, nil
}

type ListPriceSchemasUseCase interface {
	Execute(ctx context.Context, params ListSchemaParams) ([]PriceSchemaResponse, response.PaginationMeta, error)
}

type listPriceSchemasUseCase struct {
	repo Repository
}

func NewListPriceSchemasUseCase(repo Repository) ListPriceSchemasUseCase {
	return &listPriceSchemasUseCase{repo: repo}
}

func (uc *listPriceSchemasUseCase) Execute(ctx context.Context, params ListSchemaParams) ([]PriceSchemaResponse, response.PaginationMeta, error) {
	schemas, total, err := uc.repo.ListSchemas(ctx, params)
	if err != nil {
		return nil, response.PaginationMeta{}, err
	}

	res := make([]PriceSchemaResponse, len(schemas))
	for i := range schemas {
		res[i] = ToPriceSchemaResponse(&schemas[i])
	}

	totalPages := 0
	if params.Limit > 0 && total > 0 {
		totalPages = int((total + int64(params.Limit) - 1) / int64(params.Limit))
	}

	meta := response.PaginationMeta{
		CurrentPage: params.Page,
		PerPage:     params.Limit,
		TotalItems:  total,
		TotalPages:  totalPages,
	}

	return res, meta, nil
}

// =========================================================================
// 5. Upsert Price Schema Rule Use Case
// =========================================================================

type UpsertPriceSchemaRuleUseCase interface {
	Execute(ctx context.Context, schemaID string, req UpsertRuleRequest, actor OperatorContext) (*PriceSchemaRuleResponse, error)
}

type upsertPriceSchemaRuleUseCase struct {
	repo       Repository
	tx         database.TransactionManager
	audit      audit.AuditWriter
	deptReader DepartmentReader
	unitReader ServiceUnitReader
	itemReader ItemReader
}

func NewUpsertPriceSchemaRuleUseCase(
	repo Repository,
	tx database.TransactionManager,
	audit audit.AuditWriter,
	deptReader DepartmentReader,
	unitReader ServiceUnitReader,
	itemReader ItemReader,
) UpsertPriceSchemaRuleUseCase {
	return &upsertPriceSchemaRuleUseCase{
		repo:       repo,
		tx:         tx,
		audit:      audit,
		deptReader: deptReader,
		unitReader: unitReader,
		itemReader: itemReader,
	}
}

func (uc *upsertPriceSchemaRuleUseCase) Execute(ctx context.Context, schemaID string, req UpsertRuleRequest, actor OperatorContext) (*PriceSchemaRuleResponse, error) {
	// Pastikan schema induk ada
	schema, err := uc.repo.FindSchemaByID(ctx, schemaID)
	if err != nil {
		return nil, err
	}

	medAdjType := req.MedicationAdjustmentType
	if medAdjType == "" {
		medAdjType = AdjustmentTypeNone
	}
	tariffAdjType := req.TariffAdjustmentType
	if tariffAdjType == "" {
		tariffAdjType = AdjustmentTypeNone
	}

	ruleCandidate := &PriceSchemaRule{
		PriceSchemaID:             schema.ID,
		DepartementID:             cleanStringPtr(req.DepartementID),
		ServiceUnitID:             cleanStringPtr(req.ServiceUnitID),
		ItemID:                    cleanStringPtr(req.ItemID),
		TariffAdjustmentType:      tariffAdjType,
		TariffAdjustmentValue:     req.TariffAdjustmentValue,
		MedicationAdjustmentType:  medAdjType,
		MedicationAdjustmentValue: req.MedicationAdjustmentValue,
		PayByPatientPercent:       req.PayByPatientPercent,
		PayByCustomerPercent:      req.PayByCustomerPercent,
		TaxPercent:                req.TaxPercent,
		CreatedBy:                 actor.Username,
		UpdatedBy:                 actor.Username,
	}

	if err := ruleCandidate.Validate(); err != nil {
		return nil, err
	}

	// Validasi scope ke owner modul masing-masing
	if ruleCandidate.DepartementID != nil && uc.deptReader != nil {
		active, err := uc.deptReader.IsActiveDepartment(ctx, *ruleCandidate.DepartementID)
		if err != nil || !active {
			return nil, ErrDepartmentInactiveOrNotFound
		}
	}
	if ruleCandidate.ServiceUnitID != nil && uc.unitReader != nil {
		_, active, err := uc.unitReader.GetServiceUnitInfo(ctx, *ruleCandidate.ServiceUnitID)
		if err != nil || !active {
			return nil, ErrServiceUnitInactiveOrNotFound
		}
	}
	if ruleCandidate.ItemID != nil && uc.itemReader != nil {
		active, err := uc.itemReader.IsActiveItem(ctx, *ruleCandidate.ItemID)
		if err != nil || !active {
			return nil, ErrItemInactiveOrNotFound
		}
	}

	var savedRule *PriceSchemaRule

	err = uc.tx.WithTx(ctx, func(txCtx context.Context) error {
		existing, err := uc.repo.FindRuleByScope(txCtx, schemaID, ruleCandidate.DepartementID, ruleCandidate.ServiceUnitID, ruleCandidate.ItemID)
		if err != nil {
			return err
		}

		if existing != nil {
			// Update rule yang sudah ada di scope ini
			beforeSnapshot := *existing

			existing.TariffAdjustmentType = ruleCandidate.TariffAdjustmentType
			existing.TariffAdjustmentValue = ruleCandidate.TariffAdjustmentValue
			existing.MedicationAdjustmentType = ruleCandidate.MedicationAdjustmentType
			existing.MedicationAdjustmentValue = ruleCandidate.MedicationAdjustmentValue
			existing.PayByPatientPercent = ruleCandidate.PayByPatientPercent
			existing.PayByCustomerPercent = ruleCandidate.PayByCustomerPercent
			existing.TaxPercent = ruleCandidate.TaxPercent
			existing.UpdatedBy = actor.Username
			existing.UpdatedAt = time.Now().UTC()

			if err := uc.repo.UpdateRule(txCtx, existing); err != nil {
				return err
			}
			savedRule = existing

			return uc.audit.Write(txCtx, audit.AuditEntry{
				ActorID:       actor.UserID,
				ActorUsername: actor.Username,
				ActorRole:     actor.Role,
				Action:        "UPDATE_RULE",
				ResourceType:  "price_schema_rule",
				ResourceID:    existing.ID,
				Before:        beforeSnapshot,
				After:         existing,
				CorrelationID: actor.CorrelationID,
			})
		}

		// Insert rule baru
		if err := uc.repo.CreateRule(txCtx, ruleCandidate); err != nil {
			return err
		}
		savedRule = ruleCandidate

		return uc.audit.Write(txCtx, audit.AuditEntry{
			ActorID:       actor.UserID,
			ActorUsername: actor.Username,
			ActorRole:     actor.Role,
			Action:        "CREATE_RULE",
			ResourceType:  "price_schema_rule",
			ResourceID:    ruleCandidate.ID,
			After:         ruleCandidate,
			CorrelationID: actor.CorrelationID,
		})
	})

	if err != nil {
		return nil, err
	}

	resp := ToPriceSchemaRuleResponse(savedRule)
	return &resp, nil
}

// =========================================================================
// 6. Delete Price Schema Rule Use Case
// =========================================================================

type DeletePriceSchemaRuleUseCase interface {
	Execute(ctx context.Context, schemaID, ruleID string, actor OperatorContext) error
}

type deletePriceSchemaRuleUseCase struct {
	repo  Repository
	tx    database.TransactionManager
	audit audit.AuditWriter
}

func NewDeletePriceSchemaRuleUseCase(repo Repository, tx database.TransactionManager, audit audit.AuditWriter) DeletePriceSchemaRuleUseCase {
	return &deletePriceSchemaRuleUseCase{repo: repo, tx: tx, audit: audit}
}

func (uc *deletePriceSchemaRuleUseCase) Execute(ctx context.Context, schemaID, ruleID string, actor OperatorContext) error {
	rule, err := uc.repo.FindRuleByID(ctx, ruleID)
	if err != nil {
		return err
	}
	if rule.PriceSchemaID != schemaID {
		return ErrRuleNotFound
	}

	// Invariant PRD §11 Q3: Rule SCHEMA_DEFAULT tidak boleh dihapus
	if rule.Level() == MatchLevelSchemaDefault {
		return ErrCannotDeleteDefaultRule
	}

	beforeSnapshot := *rule

	return uc.tx.WithTx(ctx, func(txCtx context.Context) error {
		if err := uc.repo.DeleteRule(txCtx, ruleID); err != nil {
			return err
		}

		return uc.audit.Write(txCtx, audit.AuditEntry{
			ActorID:       actor.UserID,
			ActorUsername: actor.Username,
			ActorRole:     actor.Role,
			Action:        "DELETE_RULE",
			ResourceType:  "price_schema_rule",
			ResourceID:    ruleID,
			Before:        beforeSnapshot,
			CorrelationID: actor.CorrelationID,
		})
	})
}

// =========================================================================
// 7. List Price Schema Rules Use Case (Query Path)
// =========================================================================

type ListPriceSchemaRulesUseCase interface {
	Execute(ctx context.Context, schemaID string, params ListRuleParams) ([]PriceSchemaRuleResponse, response.PaginationMeta, error)
}

type listPriceSchemaRulesUseCase struct {
	repo Repository
}

func NewListPriceSchemaRulesUseCase(repo Repository) ListPriceSchemaRulesUseCase {
	return &listPriceSchemaRulesUseCase{repo: repo}
}

func (uc *listPriceSchemaRulesUseCase) Execute(ctx context.Context, schemaID string, params ListRuleParams) ([]PriceSchemaRuleResponse, response.PaginationMeta, error) {
	rules, total, err := uc.repo.ListRulesBySchemaID(ctx, schemaID, params)
	if err != nil {
		return nil, response.PaginationMeta{}, err
	}

	res := make([]PriceSchemaRuleResponse, len(rules))
	for i := range rules {
		res[i] = ToPriceSchemaRuleResponse(&rules[i])
	}

	totalPages := 0
	if params.Limit > 0 && total > 0 {
		totalPages = int((total + int64(params.Limit) - 1) / int64(params.Limit))
	}

	meta := response.PaginationMeta{
		CurrentPage: params.Page,
		PerPage:     params.Limit,
		TotalItems:  total,
		TotalPages:  totalPages,
	}

	return res, meta, nil
}

// =========================================================================
// 8. Simulate Tariff Pricing Use Case (PRD §5, §10)
// =========================================================================

type SimulateTariffUseCase interface {
	Execute(ctx context.Context, req SimulateTariffRequest) (*SimulateTariffResponse, error)
}

type simulateTariffUseCase struct {
	tariffLookup TariffBaseLookup
	matcher      RuleMatcher
}

func NewSimulateTariffUseCase(tariffLookup TariffBaseLookup, matcher RuleMatcher) SimulateTariffUseCase {
	return &simulateTariffUseCase{
		tariffLookup: tariffLookup,
		matcher:      matcher,
	}
}

func (uc *simulateTariffUseCase) Execute(ctx context.Context, req SimulateTariffRequest) (*SimulateTariffResponse, error) {
	req.TariffClassID = strings.TrimSpace(req.TariffClassID)
	req.ItemID = strings.TrimSpace(req.ItemID)

	if req.TariffClassID == "" || req.ItemID == "" {
		return nil, ErrTariffNotConfigured
	}

	serviceDate := time.Now().UTC()
	if strings.TrimSpace(req.ServiceDate) != "" {
		if t, err := time.Parse("2006-01-02", strings.TrimSpace(req.ServiceDate)); err == nil {
			serviceDate = t
		} else if t, err := time.Parse(time.RFC3339, strings.TrimSpace(req.ServiceDate)); err == nil {
			serviceDate = t
		}
	}

	// 1. Ambil tarif dasar RS dari Price Plan
	base, err := uc.tariffLookup.LookupTariffBase(ctx, req.TariffClassID, req.ItemID, serviceDate, req.IsCito)
	if err != nil {
		return nil, err
	}
	if base == nil {
		return nil, ErrTariffNotConfigured
	}

	// 2. Pencocokan rule penyesuaian dari Price Schema (jika ada schema_id)
	var matchedRule *MatchedRule
	if req.PriceSchemaID != nil && strings.TrimSpace(*req.PriceSchemaID) != "" {
		q := MatchRuleQuery{
			PriceSchemaID: strings.TrimSpace(*req.PriceSchemaID),
			ItemID:        req.ItemID,
			DepartementID: stringFromPtr(req.DepartementID),
			ServiceUnitID: stringFromPtr(req.ServiceUnitID),
		}
		r, err := uc.matcher.FindMatchingRule(ctx, q)
		if err != nil {
			return nil, err
		}
		matchedRule = &r
	}

	// 3. Hitung harga final penyesuaian tarif
	adjType := AdjustmentTypeNone
	adjValue := decimal.Zero
	taxPercent := decimal.Zero
	payByPatientPercent := decimal.Zero
	payByCustomerPercent := decimal.NewFromInt(100)

	var ruleID *string
	var matchLevel *MatchLevel

	if matchedRule != nil {
		adjType = matchedRule.TariffAdjustmentType
		adjValue = matchedRule.TariffAdjustmentValue
		taxPercent = matchedRule.TaxPercent
		payByPatientPercent = matchedRule.PayByPatientPercent
		payByCustomerPercent = matchedRule.PayByCustomerPercent
		ruleID = &matchedRule.RuleID
		lvl := matchedRule.Level
		matchLevel = &lvl
	} else {
		// PRD §11 Q3: Untuk pasien umum (tanpa price_schema_id), tanggungan 100% Pasien
		payByPatientPercent = decimal.NewFromInt(100)
		payByCustomerPercent = decimal.Zero
	}

	finalPriceWithoutCito := ApplyAdjustment(base.TotalBasePrice, adjType, adjValue)

	// 4. Hitung Cito surcharge
	citoSurcharge := decimal.Zero
	finalPrice := finalPriceWithoutCito
	if req.IsCito {
		citoSurcharge = base.TotalCitoPrice.Sub(base.TotalBasePrice)
		if citoSurcharge.LessThan(decimal.Zero) {
			citoSurcharge = decimal.Zero
		}
		finalPrice = finalPrice.Add(citoSurcharge)
	}

	// 5. Hitung Pajak PPN (PRD §5.2 Langkah 5, §11 Q1)
	taxAmount := CalculateTax(finalPrice, taxPercent)
	gross := finalPrice.Add(taxAmount)

	// 6. Hitung Cost Sharing (PRD §5.2 Langkah 6)
	patientAmount, customerAmount := CalculateCostSharing(gross, payByPatientPercent)

	// 7. Distribusi pecahan komponen tarif (PRD §5.5)
	componentsOutput := DistributeToComponents(finalPriceWithoutCito, base.TotalBasePrice, req.IsCito, base.Components)

	return &SimulateTariffResponse{
		PricePlanID:           base.PricePlanID,
		PricePlanCode:         base.PricePlanCode,
		BasePrice:             base.TotalBasePrice,
		CitoSurcharge:         citoSurcharge,
		IsCito:                req.IsCito,
		PriceSchemaID:         req.PriceSchemaID,
		PriceSchemaRuleID:     ruleID,
		MatchLevel:            matchLevel,
		TariffAdjustmentType:  adjType,
		TariffAdjustmentValue: adjValue,
		FinalPrice:            finalPrice,
		TaxPercent:            taxPercent,
		TaxAmount:             taxAmount,
		GrossAmount:           gross,
		PayByPatientPercent:   payByPatientPercent,
		PayByCustomerPercent:  payByCustomerPercent,
		PatientAmount:         patientAmount,
		CustomerAmount:        customerAmount,
		Components:            componentsOutput,
	}, nil
}

func cleanStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func stringFromPtr(s *string) string {
	if s == nil {
		return ""
	}
	return strings.TrimSpace(*s)
}
