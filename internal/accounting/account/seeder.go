package account

import (
	"log"
	"time"
	"uuid"

	"hosim-go/pkg/enums"

	"gorm.io/gorm"
)

type SeedCOAItem struct {
	Code              string
	ParentCode        string
	Name              string
	Type              enums.AccountType
	Position          enums.AccountPosition
	IsPostable        bool
	IsTreasuryAccount bool
	BankName          *string
	BankAccountNumber *string
	Currency          string
	Description       *string
}

//go:fix inline
func strPtr(s string) *string {
	return new(s)
}

// =========================================================================
// 1. AKUN KEPALA / INDUK (HEAD & HEADER ACCOUNTS - NON-POSTABLE)
// =========================================================================

// DefaultRootHeadAccounts memuat 5 akun kepala root level 1 (Aset, Liabilitas, Ekuitas, Pendapatan, Beban)
var DefaultRootHeadAccounts = []SeedCOAItem{
	{
		Code:        "1000",
		ParentCode:  "",
		Name:        "ASET",
		Type:        enums.AccountTypeAsset,
		Position:    enums.AccountPositionDebit,
		IsPostable:  false,
		Currency:    "IDR",
		Description: new("Akun induk seluruh aset dan aktiva rumah sakit"),
	},
	{
		Code:        "2000",
		ParentCode:  "",
		Name:        "LIABILITAS",
		Type:        enums.AccountTypeLiability,
		Position:    enums.AccountPositionCredit,
		IsPostable:  false,
		Currency:    "IDR",
		Description: new("Akun induk seluruh kewajiban dan hutang rumah sakit"),
	},
	{
		Code:        "3000",
		ParentCode:  "",
		Name:        "EKUITAS",
		Type:        enums.AccountTypeEquity,
		Position:    enums.AccountPositionCredit,
		IsPostable:  false,
		Currency:    "IDR",
		Description: new("Akun induk modal dan ekuitas rumah sakit"),
	},
	{
		Code:        "4000",
		ParentCode:  "",
		Name:        "PENDAPATAN OPERASIONAL",
		Type:        enums.AccountTypeRevenue,
		Position:    enums.AccountPositionCredit,
		IsPostable:  false,
		Currency:    "IDR",
		Description: new("Akun induk seluruh pendapatan operasional pelayanan kesehatan"),
	},
	{
		Code:        "5000",
		ParentCode:  "",
		Name:        "BEBAN OPERASIONAL",
		Type:        enums.AccountTypeExpense,
		Position:    enums.AccountPositionDebit,
		IsPostable:  false,
		Currency:    "IDR",
		Description: new("Akun induk beban dan pengeluaran operasional rumah sakit"),
	},
}

// DefaultHeadAccounts memuat seluruh akun induk/kepala (Level 1 Root Header, Level 2 Header, dan Level 3 Header)
var DefaultHeadAccounts = []SeedCOAItem{
	// ==========================================
	// LEVEL 1: AKUN INDUK UTAMA (ROOT HEADER)
	// ==========================================
	{
		Code:        "1000",
		ParentCode:  "",
		Name:        "ASET",
		Type:        enums.AccountTypeAsset,
		Position:    enums.AccountPositionDebit,
		IsPostable:  false,
		Currency:    "IDR",
		Description: new("Akun induk seluruh aset dan aktiva rumah sakit"),
	},
	{
		Code:        "2000",
		ParentCode:  "",
		Name:        "LIABILITAS",
		Type:        enums.AccountTypeLiability,
		Position:    enums.AccountPositionCredit,
		IsPostable:  false,
		Currency:    "IDR",
		Description: new("Akun induk seluruh kewajiban dan hutang rumah sakit"),
	},
	{
		Code:        "3000",
		ParentCode:  "",
		Name:        "EKUITAS",
		Type:        enums.AccountTypeEquity,
		Position:    enums.AccountPositionCredit,
		IsPostable:  false,
		Currency:    "IDR",
		Description: new("Akun induk modal dan ekuitas rumah sakit"),
	},
	{
		Code:        "4000",
		ParentCode:  "",
		Name:        "PENDAPATAN OPERASIONAL",
		Type:        enums.AccountTypeRevenue,
		Position:    enums.AccountPositionCredit,
		IsPostable:  false,
		Currency:    "IDR",
		Description: new("Akun induk seluruh pendapatan operasional pelayanan kesehatan"),
	},
	{
		Code:        "5000",
		ParentCode:  "",
		Name:        "BEBAN OPERASIONAL",
		Type:        enums.AccountTypeExpense,
		Position:    enums.AccountPositionDebit,
		IsPostable:  false,
		Currency:    "IDR",
		Description: new("Akun induk beban dan pengeluaran operasional rumah sakit"),
	},

	// ==========================================
	// LEVEL 2: KELOMPOK AKUN HEADER
	// ==========================================
	{
		Code:        "1100",
		ParentCode:  "1000",
		Name:        "ASET LANCAR",
		Type:        enums.AccountTypeAsset,
		Position:    enums.AccountPositionDebit,
		IsPostable:  false,
		Currency:    "IDR",
		Description: new("Kelompok aset lancar, kas, piutang, dan persediaan"),
	},
	{
		Code:        "1200",
		ParentCode:  "1000",
		Name:        "ASET TETAP",
		Type:        enums.AccountTypeAsset,
		Position:    enums.AccountPositionDebit,
		IsPostable:  false,
		Currency:    "IDR",
		Description: new("Kelompok aset tetap dan peralatan medis"),
	},
	{
		Code:        "2100",
		ParentCode:  "2000",
		Name:        "LIABILITAS JANGKA PENDEK",
		Type:        enums.AccountTypeLiability,
		Position:    enums.AccountPositionCredit,
		IsPostable:  false,
		Currency:    "IDR",
		Description: new("Hutang lancar operasional dan kewajiban jangka pendek"),
	},
	{
		Code:        "4100",
		ParentCode:  "4000",
		Name:        "Pendapatan Pelayanan Medis RS",
		Type:        enums.AccountTypeRevenue,
		Position:    enums.AccountPositionCredit,
		IsPostable:  false,
		Currency:    "IDR",
		Description: new("Kelompok pendapatan pelayanan tindakan dan sarana"),
	},
	{
		Code:        "4200",
		ParentCode:  "4000",
		Name:        "Pendapatan Farmasi & BMHP",
		Type:        enums.AccountTypeRevenue,
		Position:    enums.AccountPositionCredit,
		IsPostable:  false,
		Currency:    "IDR",
		Description: new("Kelompok pendapatan dari penjualan obat dan BMHP"),
	},

	// ==========================================
	// LEVEL 3: SUB-KELOMPOK HEADER
	// ==========================================
	{
		Code:        "1110",
		ParentCode:  "1100",
		Name:        "Kas dan Setara Kas",
		Type:        enums.AccountTypeAsset,
		Position:    enums.AccountPositionDebit,
		IsPostable:  false,
		Currency:    "IDR",
		Description: new("Header kas fisik dan rekening bank operasional"),
	},
	{
		Code:        "1120",
		ParentCode:  "1100",
		Name:        "Piutang Pelayanan Pasien",
		Type:        enums.AccountTypeAsset,
		Position:    enums.AccountPositionDebit,
		IsPostable:  false,
		Currency:    "IDR",
		Description: new("Header piutang tagihan pasien dan asuransi/BPJS"),
	},
	{
		Code:        "1130",
		ParentCode:  "1100",
		Name:        "Persediaan Medis & Non-Medis",
		Type:        enums.AccountTypeAsset,
		Position:    enums.AccountPositionDebit,
		IsPostable:  false,
		Currency:    "IDR",
		Description: new("Header persediaan obat, BMHP, dan logistik"),
	},
}

// =========================================================================
// 2. AKUN DETAIL / POSTABLE (LEAF ACCOUNTS - LEVEL 2, 3, 4)
// =========================================================================

// DefaultDetailAccounts memuat seluruh akun detail yang siap dipilih transaksi (is_postable = true)
var DefaultDetailAccounts = []SeedCOAItem{
	// Level 2 Ekuitas Leaf
	{
		Code:        "3100",
		ParentCode:  "3000",
		Name:        "Modal Disetor / Modal Pemilik",
		Type:        enums.AccountTypeEquity,
		Position:    enums.AccountPositionCredit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Modal awal disetor oleh pemilik/yayasan"),
	},
	{
		Code:        "3200",
		ParentCode:  "3000",
		Name:        "Saldo Laba Ditahan",
		Type:        enums.AccountTypeEquity,
		Position:    enums.AccountPositionCredit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Akumulasi laba/rugi periode berjalan dan periode lalu"),
	},

	// Level 2 Beban Leaf
	{
		Code:        "5100",
		ParentCode:  "5000",
		Name:        "Harga Pokok Penjualan / HPP Obat & BMHP",
		Type:        enums.AccountTypeExpense,
		Position:    enums.AccountPositionDebit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Harga pokok obat dan bahan medis habis pakai yang terjual"),
	},
	{
		Code:        "5200",
		ParentCode:  "5000",
		Name:        "Beban Jasa Dokter & Tenaga Medis",
		Type:        enums.AccountTypeExpense,
		Position:    enums.AccountPositionDebit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Beban jasa medis dan remunerasi dokter"),
	},
	{
		Code:        "5300",
		ParentCode:  "5000",
		Name:        "Beban Gaji & Kesejahteraan Karyawan",
		Type:        enums.AccountTypeExpense,
		Position:    enums.AccountPositionDebit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Beban gaji, tunjangan, dan asuransi staf rumah sakit"),
	},
	{
		Code:        "5400",
		ParentCode:  "5000",
		Name:        "Beban Pemeliharaan & Operasional Sarana RS",
		Type:        enums.AccountTypeExpense,
		Position:    enums.AccountPositionDebit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Biaya listrik, air, kalibrasi alat medis, dan pemeliharaan gedung"),
	},

	// Level 3 Aset Tetap Leaf
	{
		Code:        "1210",
		ParentCode:  "1200",
		Name:        "Peralatan Medis & Laboratorium",
		Type:        enums.AccountTypeAsset,
		Position:    enums.AccountPositionDebit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Nilai perolehan peralatan medis dan mesin penunjang"),
	},
	{
		Code:        "1219",
		ParentCode:  "1200",
		Name:        "Akumulasi Penyusutan Alat Medis",
		Type:        enums.AccountTypeAsset,
		Position:    enums.AccountPositionCredit, // Akun Kontra
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Akumulasi amortisasi/penyusutan peralatan medis"),
	},

	// Level 3 Liabilitas Leaf
	{
		Code:        "2110",
		ParentCode:  "2100",
		Name:        "Hutang Usaha / Pembelian Obat PBF",
		Type:        enums.AccountTypeLiability,
		Position:    enums.AccountPositionCredit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Hutang tagihan pembelian obat ke Pedagang Besar Farmasi"),
	},
	{
		Code:        "2120",
		ParentCode:  "2100",
		Name:        "Hutang Jasa Medis Dokter",
		Type:        enums.AccountTypeLiability,
		Position:    enums.AccountPositionCredit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Kewajiban jasa medis dokter yang belum dibayarkan"),
	},
	{
		Code:        "2130",
		ParentCode:  "2100",
		Name:        "Hutang Sewa Alat Medis Rekanan",
		Type:        enums.AccountTypeLiability,
		Position:    enums.AccountPositionCredit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Hutang biaya sewa operasional alat diagnostik"),
	},
	{
		Code:        "2140",
		ParentCode:  "2100",
		Name:        "Uang Muka Biaya Rawat Inap",
		Type:        enums.AccountTypeLiability,
		Position:    enums.AccountPositionCredit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Deposit / uang muka yang disetorkan pasien rawat inap"),
	},

	// Level 3 Pendapatan Leaf
	{
		Code:        "4101",
		ParentCode:  "4100",
		Name:        "Pendapatan Jasa Sarana & Kamar RS",
		Type:        enums.AccountTypeRevenue,
		Position:    enums.AccountPositionCredit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Pendapatan sewa tempat tidur dan fasilitas ruang rawat"),
	},
	{
		Code:        "4102",
		ParentCode:  "4100",
		Name:        "Pendapatan Jasa Medis Tindakan Dokter",
		Type:        enums.AccountTypeRevenue,
		Position:    enums.AccountPositionCredit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Pendapatan jasa visitasi, konsultasi, dan operasi dokter"),
	},
	{
		Code:        "4103",
		ParentCode:  "4100",
		Name:        "Pendapatan Jasa Paramedis & Asuhan Keperawatan",
		Type:        enums.AccountTypeRevenue,
		Position:    enums.AccountPositionCredit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Pendapatan tindakan keperawatan dan bidan"),
	},
	{
		Code:        "4104",
		ParentCode:  "4100",
		Name:        "Pendapatan Tindakan Laboratorium & Radiologi",
		Type:        enums.AccountTypeRevenue,
		Position:    enums.AccountPositionCredit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Pendapatan pemeriksaan laboratorium patologi dan rontgen"),
	},
	{
		Code:        "4201",
		ParentCode:  "4200",
		Name:        "Pendapatan Penjualan Obat Rawat Jalan",
		Type:        enums.AccountTypeRevenue,
		Position:    enums.AccountPositionCredit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Pendapatan resep farmasi instalasi rawat jalan"),
	},
	{
		Code:        "4202",
		ParentCode:  "4200",
		Name:        "Pendapatan Penjualan Obat Rawat Inap",
		Type:        enums.AccountTypeRevenue,
		Position:    enums.AccountPositionCredit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Pendapatan paket obat dan BMHP pasien ranap"),
	},

	// Level 4 Kas & Setara Kas (Treasury Accounts)
	{
		Code:              "1111",
		ParentCode:        "1110",
		Name:              "Kas Kasir Utama",
		Type:              enums.AccountTypeAsset,
		Position:          enums.AccountPositionDebit,
		IsPostable:        true,
		IsTreasuryAccount: true,
		Currency:          "IDR",
		Description:       new("Kas fisik loket kasir penerimaan pembayaran"),
	},
	{
		Code:              "1112",
		ParentCode:        "1110",
		Name:              "Kas Kecil Depo Farmasi & IGD",
		Type:              enums.AccountTypeAsset,
		Position:          enums.AccountPositionDebit,
		IsPostable:        true,
		IsTreasuryAccount: true,
		Currency:          "IDR",
		Description:       new("Kas kecil operasional untuk depo farmasi dan kasir IGD 24 jam"),
	},
	{
		Code:              "1113",
		ParentCode:        "1110",
		Name:              "Bank Mandiri Operasional",
		Type:              enums.AccountTypeAsset,
		Position:          enums.AccountPositionDebit,
		IsPostable:        true,
		IsTreasuryAccount: true,
		BankName:          new("Bank Mandiri"),
		BankAccountNumber: new("137-00-1234567-8"),
		Currency:          "IDR",
		Description:       new("Rekening giro operasional Rumah Sakit"),
	},
	{
		Code:              "1114",
		ParentCode:        "1110",
		Name:              "Bank BCA Operasional",
		Type:              enums.AccountTypeAsset,
		Position:          enums.AccountPositionDebit,
		IsPostable:        true,
		IsTreasuryAccount: true,
		BankName:          new("Bank Central Asia"),
		BankAccountNumber: new("028-9876543"),
		Currency:          "IDR",
		Description:       new("Rekening operasional penerimaan EDC dan transfer"),
	},

	// Level 4 Piutang Leaf
	{
		Code:        "1121",
		ParentCode:  "1120",
		Name:        "Piutang Pasien Umum / Pribadi",
		Type:        enums.AccountTypeAsset,
		Position:    enums.AccountPositionDebit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Tagihan pelayanan pasien umum belum terbayar"),
	},
	{
		Code:        "1122",
		ParentCode:  "1120",
		Name:        "Piutang BPJS Kesehatan",
		Type:        enums.AccountTypeAsset,
		Position:    enums.AccountPositionDebit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Klaim tagihan BPJS Kesehatan dalam proses verifikasi dan pembayaran"),
	},
	{
		Code:        "1123",
		ParentCode:  "1120",
		Name:        "Piutang Asuransi Swasta & Korporasi",
		Type:        enums.AccountTypeAsset,
		Position:    enums.AccountPositionDebit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Piutang jaminan perusahaan rekanan dan asuransi komersial"),
	},

	// Level 4 Persediaan Leaf
	{
		Code:        "1131",
		ParentCode:  "1130",
		Name:        "Persediaan Obat-Obatan Farmasi",
		Type:        enums.AccountTypeAsset,
		Position:    enums.AccountPositionDebit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Nilai persediaan obat di gudang farmasi dan depo"),
	},
	{
		Code:        "1132",
		ParentCode:  "1130",
		Name:        "Persediaan Bahan Medis Habis Pakai / BMHP",
		Type:        enums.AccountTypeAsset,
		Position:    enums.AccountPositionDebit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Nilai persediaan spuit, infus, benang bedah, dan kassa"),
	},
	{
		Code:        "1133",
		ParentCode:  "1130",
		Name:        "Persediaan Logistik Umum & ATK",
		Type:        enums.AccountTypeAsset,
		Position:    enums.AccountPositionDebit,
		IsPostable:  true,
		Currency:    "IDR",
		Description: new("Persediaan perlengkapan kantor, form rekam medis, dan linen"),
	},
}

// DefaultHospitalCOA mendefinisikan bagan akun lengkap standar rumah sakit (Head + Detail)
var DefaultHospitalCOA = append(append([]SeedCOAItem{}, DefaultHeadAccounts...), DefaultDetailAccounts...)

// =========================================================================
// 3. FUNGSI SEEDER
// =========================================================================

// seedAccountsList menginisialisasi daftar akun secara bertahap dan idempotent
func seedAccountsList(db *gorm.DB, items []SeedCOAItem) error {
	now := time.Now()

	accountMap := make(map[string]ChartOfAccount)

	var existingAccounts []ChartOfAccount
	if err := db.Find(&existingAccounts).Error; err != nil {
		return err
	}
	for _, acc := range existingAccounts {
		accountMap[acc.Code] = acc
	}

	for _, item := range items {
		if existing, found := accountMap[item.Code]; found {
			if existing.IsPostable != item.IsPostable && !item.IsPostable {
				_ = db.Model(&ChartOfAccount{}).Where("id = ?", existing.ID).Update("is_postable", item.IsPostable)
			}
			continue
		}

		var parentID *string
		var parentCode *string
		level := 1

		if item.ParentCode != "" {
			parent, parentFound := accountMap[item.ParentCode]
			if parentFound {
				pID := parent.ID
				pCode := parent.Code
				parentID = &pID
				parentCode = &pCode
				level = parent.AccountLevel + 1

				if parent.IsPostable {
					if err := db.Model(&ChartOfAccount{}).Where("id = ?", parent.ID).Update("is_postable", false).Error; err != nil {
						return err
					}
					parent.IsPostable = false
					accountMap[item.ParentCode] = parent
				}
			}
		}

		newAccount := ChartOfAccount{
			ID:                uuid.NewV7().String(),
			ParentID:          parentID,
			ParentCode:        parentCode,
			Code:              item.Code,
			Name:              item.Name,
			Description:       item.Description,
			Type:              item.Type,
			Position:          item.Position,
			AccountLevel:      level,
			IsPostable:        item.IsPostable,
			IsTreasuryAccount: item.IsTreasuryAccount,
			BankName:          item.BankName,
			BankAccountNumber: item.BankAccountNumber,
			Currency:          item.Currency,
			IsActive:          true,
			CreatedBy:         "system-seeder",
			UpdatedBy:         "system-seeder",
			CreatedAt:         now,
			UpdatedAt:         now,
		}

		if err := db.Select("*").Create(&newAccount).Error; err != nil {
			return err
		}

		accountMap[newAccount.Code] = newAccount
	}

	return nil
}

// SeedRootHeadAccounts menginisialisasi 5 akun kepala root level 1 (Aset, Liabilitas, Ekuitas, Pendapatan, Beban)
func SeedRootHeadAccounts(db *gorm.DB) error {
	log.Println("[SEEDER] Memeriksa & melakukan seeding Akun Induk Root (Level 1)...")
	if err := seedAccountsList(db, DefaultRootHeadAccounts); err != nil {
		return err
	}
	log.Println("[SEEDER] Seeding Akun Induk Root berhasil.")
	return nil
}

// SeedHeadAccounts menginisialisasi seluruh akun kepala / header (Level 1 s/d Level 3 Header)
func SeedHeadAccounts(db *gorm.DB) error {
	log.Println("[SEEDER] Memeriksa & melakukan seeding seluruh Akun Kepala / Header (Chart of Accounts)...")
	if err := seedAccountsList(db, DefaultHeadAccounts); err != nil {
		return err
	}
	log.Println("[SEEDER] Seeding Akun Kepala / Header berhasil.")
	return nil
}

// SeedChartOfAccounts menginisialisasi seluruh bagan akun rumah sakit lengkap (Head + Detail)
func SeedChartOfAccounts(db *gorm.DB) error {
	log.Println("[SEEDER] Memeriksa & melakukan seeding Bagan Akun Lengkap (Chart of Accounts)...")
	if err := seedAccountsList(db, DefaultHospitalCOA); err != nil {
		return err
	}
	log.Println("[SEEDER] Seeding Bagan Akun Lengkap berhasil.")
	return nil
}
