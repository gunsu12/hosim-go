package enums

// enums untuk tipe gudang / storage
type StorageType string

const (
	// gudang sentral untuk kedatangan barang
	StorageTypeCentral StorageType = "central"
	// depo ini ada di pecahan strage mana yang ada di unit
	StorageTypeDepo StorageType = "depo"
	// ini tipe pembantu karena ketika distribusi maka barang akan nyangkut disini sebelum di konfirmasi
	StorageTypeVirtual StorageType = "virtual"
)

func (s StorageType) String() string {
	return string(s)
}
