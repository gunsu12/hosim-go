package priceschema

import (
	"strings"

	"github.com/shopspring/decimal"
)

var (
	zeroDec    = decimal.Zero
	hundredDec = decimal.NewFromInt(100)
)

// ComponentInput merepresentasikan data komponen dari buku tarif dasar (Price Plan).
type ComponentInput struct {
	ComponentID       string          `json:"component_id"`
	ComponentName     string          `json:"component_name"`
	ComponentType     string          `json:"component_type"`
	COACode           string          `json:"coa_code"`
	BaseAmount        decimal.Decimal `json:"base_amount"`
	CitoAmount        decimal.Decimal `json:"cito_amount"`
	IsHospitalRevenue bool            `json:"is_hospital_revenue"`
}

// ComponentOutput merepresentasikan hasil distribusi penyesuaian tarif ke komponen.
type ComponentOutput struct {
	ComponentID   string          `json:"component_id"`
	ComponentName string          `json:"component_name"`
	ComponentType string          `json:"component_type"`
	COACode       string          `json:"coa_code"`
	Amount        decimal.Decimal `json:"amount"`
}

// ApplyAdjustment menghitung harga tarif setelah penyesuaian (PRD §5.2 Langkah 3).
// Nominal akhir di-clamp minimum 0 dan dibulatkan dengan ROUND_HALF_UP ke 0 desimal.
func ApplyAdjustment(basePrice decimal.Decimal, adjType AdjustmentType, adjValue decimal.Decimal) decimal.Decimal {
	var finalPrice decimal.Decimal

	switch adjType {
	case AdjustmentTypePercent:
		factor := decimal.NewFromInt(1).Add(adjValue.Div(hundredDec))
		finalPrice = basePrice.Mul(factor)
	case AdjustmentTypeFixedAmount:
		finalPrice = basePrice.Add(adjValue)
	case AdjustmentTypeFixedPrice:
		finalPrice = adjValue
	case AdjustmentTypeNone:
		fallthrough
	default:
		finalPrice = basePrice
	}

	if finalPrice.LessThan(zeroDec) {
		finalPrice = zeroDec
	}

	return finalPrice.Round(0)
}

// CalculateTax menghitung nilai pajak PPN dari harga final (PRD §5.2 Langkah 5).
func CalculateTax(finalPrice decimal.Decimal, taxPercent decimal.Decimal) decimal.Decimal {
	if taxPercent.IsZero() || finalPrice.IsZero() {
		return zeroDec
	}
	return finalPrice.Mul(taxPercent).Div(hundredDec).Round(0)
}

// CalculateCostSharing membagi total tagihan bruto ke tanggungan pasien dan customer/penjamin (PRD §5.2 Langkah 6).
// Sisa pembulatan dialokasikan ke porsi customer agar patient_amount + customer_amount == gross.
func CalculateCostSharing(gross decimal.Decimal, payByPatientPercent decimal.Decimal) (patientAmount, customerAmount decimal.Decimal) {
	if gross.IsZero() {
		return zeroDec, zeroDec
	}
	patientAmount = gross.Mul(payByPatientPercent).Div(hundredDec).Round(0)
	customerAmount = gross.Sub(patientAmount)
	return patientAmount, customerAmount
}

// DistributeToComponents mendistribusikan penyesuaian tarif dan cito ke masing-masing komponen biaya (PRD §5.5).
func DistributeToComponents(
	finalPriceTanpaCito decimal.Decimal,
	totalBasePrice decimal.Decimal,
	isCito bool,
	baseComponents []ComponentInput,
) []ComponentOutput {
	if len(baseComponents) == 0 {
		return []ComponentOutput{}
	}

	delta := finalPriceTanpaCito.Sub(totalBasePrice)
	output := make([]ComponentOutput, len(baseComponents))

	// Cari indeks komponen Jasa RS / Sarana untuk menampung sisa pembulatan
	targetHospitalIdx := 0
	foundHospital := false
	for idx, c := range baseComponents {
		t := strings.ToUpper(c.ComponentType)
		n := strings.ToUpper(c.ComponentName)
		if c.IsHospitalRevenue || strings.Contains(t, "RS") || strings.Contains(t, "SARANA") ||
			strings.Contains(n, "RS") || strings.Contains(n, "SARANA") {
			targetHospitalIdx = idx
			foundHospital = true
			break
		}
	}
	if !foundHospital {
		targetHospitalIdx = 0
	}

	if totalBasePrice.IsZero() {
		// Jika base price 0, seluruh delta masuk ke komponen Jasa RS
		for i, c := range baseComponents {
			amt := c.BaseAmount
			if i == targetHospitalIdx {
				amt = amt.Add(delta)
			}
			output[i] = ComponentOutput{
				ComponentID:   c.ComponentID,
				ComponentName: c.ComponentName,
				ComponentType: c.ComponentType,
				COACode:       c.COACode,
				Amount:        amt.Round(0),
			}
		}
	} else {
		// Hitung proporsi masing-masing komponen: amount_i = ROUND(base_i + delta * base_i / total_base)
		sumDistributed := zeroDec
		for i, c := range baseComponents {
			adjPart := delta.Mul(c.BaseAmount).Div(totalBasePrice)
			amt := c.BaseAmount.Add(adjPart).Round(0)
			if amt.LessThan(zeroDec) {
				amt = zeroDec
			}
			output[i] = ComponentOutput{
				ComponentID:   c.ComponentID,
				ComponentName: c.ComponentName,
				ComponentType: c.ComponentType,
				COACode:       c.COACode,
				Amount:        amt,
			}
			sumDistributed = sumDistributed.Add(amt)
		}

		// Sisa pembulatan dialokasikan ke komponen Jasa RS / Sarana
		roundingDiff := finalPriceTanpaCito.Sub(sumDistributed)
		if !roundingDiff.IsZero() {
			output[targetHospitalIdx].Amount = output[targetHospitalIdx].Amount.Add(roundingDiff)
		}
	}

	// Jika CITO, tambahkan cito surcharge per komponen dari Price Plan
	if isCito {
		for i, c := range baseComponents {
			citoSurcharge := c.CitoAmount.Sub(c.BaseAmount)
			if citoSurcharge.GreaterThan(zeroDec) {
				output[i].Amount = output[i].Amount.Add(citoSurcharge.Round(0))
			}
		}
	}

	return output
}
