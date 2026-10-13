package priceschema

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestApplyAdjustment(t *testing.T) {
	tests := []struct {
		name      string
		basePrice decimal.Decimal
		adjType   AdjustmentType
		adjValue  decimal.Decimal
		expected  decimal.Decimal
	}{
		{
			name:      "None adjustment",
			basePrice: decimal.NewFromInt(200000),
			adjType:   AdjustmentTypeNone,
			adjValue:  decimal.Zero,
			expected:  decimal.NewFromInt(200000),
		},
		{
			name:      "Discount 10 percent",
			basePrice: decimal.NewFromInt(200000),
			adjType:   AdjustmentTypePercent,
			adjValue:  decimal.NewFromInt(-10),
			expected:  decimal.NewFromInt(180000),
		},
		{
			name:      "Markup 15 percent",
			basePrice: decimal.NewFromInt(100000),
			adjType:   AdjustmentTypePercent,
			adjValue:  decimal.NewFromInt(15),
			expected:  decimal.NewFromInt(115000),
		},
		{
			name:      "Fixed amount markup",
			basePrice: decimal.NewFromInt(100000),
			adjType:   AdjustmentTypeFixedAmount,
			adjValue:  decimal.NewFromInt(25000),
			expected:  decimal.NewFromInt(125000),
		},
		{
			name:      "Fixed amount discount clamp to zero",
			basePrice: decimal.NewFromInt(50000),
			adjType:   AdjustmentTypeFixedAmount,
			adjValue:  decimal.NewFromInt(-80000),
			expected:  decimal.Zero,
		},
		{
			name:      "Fixed price override",
			basePrice: decimal.NewFromInt(250000),
			adjType:   AdjustmentTypeFixedPrice,
			adjValue:  decimal.NewFromInt(175000),
			expected:  decimal.NewFromInt(175000),
		},
		{
			name:      "Rounding half up",
			basePrice: decimal.NewFromFloat(100.5),
			adjType:   AdjustmentTypePercent,
			adjValue:  decimal.NewFromInt(5), // 100.5 * 1.05 = 105.525 -> 106
			expected:  decimal.NewFromInt(106),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApplyAdjustment(tt.basePrice, tt.adjType, tt.adjValue)
			if !got.Equal(tt.expected) {
				t.Errorf("ApplyAdjustment() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestCalculateTaxAndCostSharing(t *testing.T) {
	// Harga final 180.000, pajak 11% -> 19.800, gross = 199.800
	finalPrice := decimal.NewFromInt(180000)
	taxPercent := decimal.NewFromInt(11)

	tax := CalculateTax(finalPrice, taxPercent)
	expectedTax := decimal.NewFromInt(19800)
	if !tax.Equal(expectedTax) {
		t.Fatalf("CalculateTax() = %v, expected %v", tax, expectedTax)
	}

	gross := finalPrice.Add(tax) // 199.800

	// Cost sharing: Patient 20%, Customer 80%
	// Patient = 199.800 * 0.2 = 39.960
	// Customer = 199.800 - 39.960 = 159.840
	patientAmt, customerAmt := CalculateCostSharing(gross, decimal.NewFromInt(20))
	expectedPatient := decimal.NewFromInt(39960)
	expectedCustomer := decimal.NewFromInt(159840)

	if !patientAmt.Equal(expectedPatient) {
		t.Errorf("CalculateCostSharing() patient = %v, expected %v", patientAmt, expectedPatient)
	}
	if !customerAmt.Equal(expectedCustomer) {
		t.Errorf("CalculateCostSharing() customer = %v, expected %v", customerAmt, expectedCustomer)
	}
	if !patientAmt.Add(customerAmt).Equal(gross) {
		t.Errorf("Sum of patient and customer (%v) != gross (%v)", patientAmt.Add(customerAmt), gross)
	}
}

func TestDistributeToComponents_PRDExample(t *testing.T) {
	// PRD §5.6: Base 200.000 (Jasa RS 120.000, Jasa Dokter 80.000), Discount 10% -> Final 180.000
	baseComponents := []ComponentInput{
		{
			ComponentID:       "comp-rs",
			ComponentName:     "Jasa RS",
			ComponentType:     "SARANA",
			COACode:           "410.01",
			BaseAmount:        decimal.NewFromInt(120000),
			CitoAmount:        decimal.NewFromInt(150000),
			IsHospitalRevenue: true,
		},
		{
			ComponentID:       "comp-dr",
			ComponentName:     "Jasa Dokter",
			ComponentType:     "MEDIS",
			COACode:           "214.01",
			BaseAmount:        decimal.NewFromInt(80000),
			CitoAmount:        decimal.NewFromInt(100000),
			IsHospitalRevenue: false,
		},
	}

	totalBasePrice := decimal.NewFromInt(200000)
	finalPriceWithoutCito := decimal.NewFromInt(180000)

	// Reguler (bukan cito)
	dist := DistributeToComponents(finalPriceWithoutCito, totalBasePrice, false, baseComponents)
	if len(dist) != 2 {
		t.Fatalf("expected 2 components, got %d", len(dist))
	}

	if !dist[0].Amount.Equal(decimal.NewFromInt(108000)) {
		t.Errorf("Jasa RS amount = %v, expected 108000", dist[0].Amount)
	}
	if !dist[1].Amount.Equal(decimal.NewFromInt(72000)) {
		t.Errorf("Jasa Dokter amount = %v, expected 72000", dist[1].Amount)
	}

	// CITO test: Cito surcharge RS = 150.000 - 120.000 = 30.000, Dokter = 100.000 - 80.000 = 20.000
	distCito := DistributeToComponents(finalPriceWithoutCito, totalBasePrice, true, baseComponents)
	if !distCito[0].Amount.Equal(decimal.NewFromInt(138000)) { // 108.000 + 30.000
		t.Errorf("Cito Jasa RS amount = %v, expected 138000", distCito[0].Amount)
	}
	if !distCito[1].Amount.Equal(decimal.NewFromInt(92000)) { // 72.000 + 20.000
		t.Errorf("Cito Jasa Dokter amount = %v, expected 92000", distCito[1].Amount)
	}
}
