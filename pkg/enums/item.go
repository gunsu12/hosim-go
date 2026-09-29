package enums

// ItemType merepresentasikan tipe universal barang dan layanan di katalog rumah sakit
type ItemType string

const (
	ItemTypeMedication ItemType = "MEDICATION"
	ItemTypeGeneral    ItemType = "GENERAL"
	ItemTypeAsset      ItemType = "ASSET"
	ItemTypeTariff     ItemType = "TARIFF"
)

func (t ItemType) IsValid() bool {
	switch t {
	case ItemTypeMedication, ItemTypeGeneral, ItemTypeAsset, ItemTypeTariff:
		return true
	default:
		return false
	}
}

func (t ItemType) String() string {
	return string(t)
}

// GeneralType merepresentasikan subtipe barang umum / BMHP / logistik non-medis
type GeneralType string

const (
	GeneralTypeBMHPMedis      GeneralType = "BMHP_MEDIS"
	GeneralTypeInstrumenMedis GeneralType = "INSTRUMEN_MEDIS"
	GeneralTypeATK            GeneralType = "ATK"
	GeneralTypeLinen          GeneralType = "LINEN"
	GeneralTypeKebersihan     GeneralType = "KEBERSIHAN"
	GeneralTypeDapur          GeneralType = "DAPUR"
	GeneralTypeLainnya        GeneralType = "LAINNYA"
)

func (gt GeneralType) IsValid() bool {
	switch gt {
	case GeneralTypeBMHPMedis,
		GeneralTypeInstrumenMedis,
		GeneralTypeATK,
		GeneralTypeLinen,
		GeneralTypeKebersihan,
		GeneralTypeDapur,
		GeneralTypeLainnya:
		return true
	default:
		return false
	}
}

func (gt GeneralType) String() string {
	return string(gt)
}

// MedicationType merepresentasikan golongan obat menurut regulasi farmasi
type MedicationType string

const (
	MedicationTypeBebas        MedicationType = "Bebas"
	MedicationTypeKeras        MedicationType = "Keras"
	MedicationTypeNarkotika    MedicationType = "Narkotika"
	MedicationTypePsikotropika MedicationType = "Psikotropika"
	MedicationTypePrekursor    MedicationType = "Prekursor"
	MedicationTypeLainnya      MedicationType = "Lainnya"
)

func (mt MedicationType) IsValid() bool {
	switch mt {
	case MedicationTypeBebas,
		MedicationTypeKeras,
		MedicationTypeNarkotika,
		MedicationTypePsikotropika,
		MedicationTypePrekursor,
		MedicationTypeLainnya:
		return true
	default:
		return false
	}
}

func (mt MedicationType) String() string {
	return string(mt)
}

// SterilizationMethod merepresentasikan metode sterilisasi siklus CSSD
type SterilizationMethod string

const (
	SterilizationMethodSteamAutoclave SterilizationMethod = "STEAM_AUTOCLAVE"
	SterilizationMethodEOGas          SterilizationMethod = "EO_GAS"
	SterilizationMethodPlasma         SterilizationMethod = "PLASMA"
	SterilizationMethodDryHeat        SterilizationMethod = "DRY_HEAT"
)

func (sm SterilizationMethod) IsValid() bool {
	switch sm {
	case SterilizationMethodSteamAutoclave,
		SterilizationMethodEOGas,
		SterilizationMethodPlasma,
		SterilizationMethodDryHeat:
		return true
	default:
		return false
	}
}

func (sm SterilizationMethod) String() string {
	return string(sm)
}

// DepreciationMethod merepresentasikan metode penyusutan aktiva tetap / alat modal
type DepreciationMethod string

const (
	DepreciationMethodStraightLine    DepreciationMethod = "STRAIGHT_LINE"
	DepreciationMethodDoubleDeclining DepreciationMethod = "DOUBLE_DECLINING"
	DepreciationMethodSumOfYears      DepreciationMethod = "SUM_OF_YEARS"
)

func (dm DepreciationMethod) IsValid() bool {
	switch dm {
	case DepreciationMethodStraightLine,
		DepreciationMethodDoubleDeclining,
		DepreciationMethodSumOfYears:
		return true
	default:
		return false
	}
}

func (dm DepreciationMethod) String() string {
	return string(dm)
}

// ChargeType merepresentasikan klasifikasi beban tagihan layanan pasien
type ChargeType string

const (
	ChargeTypeAdministrasi ChargeType = "ADMINISTRASI"
	ChargeTypeAkomodasi    ChargeType = "AKOMODASI"
	ChargeTypeTindakan     ChargeType = "TINDAKAN"
	ChargeTypePenunjang    ChargeType = "PENUNJANG"
	ChargeTypeLainnya      ChargeType = "LAINNYA"
)

func (ct ChargeType) IsValid() bool {
	switch ct {
	case ChargeTypeAdministrasi,
		ChargeTypeAkomodasi,
		ChargeTypeTindakan,
		ChargeTypePenunjang,
		ChargeTypeLainnya:
		return true
	default:
		return false
	}
}

func (ct ChargeType) String() string {
	return string(ct)
}
