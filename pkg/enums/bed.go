package enums

type BedStatus string

const (
	BedStatusAvailable   BedStatus = "AVAILABLE"
	BedStatusReserved    BedStatus = "RESERVED"
	BedStatusOccupied    BedStatus = "OCCUPIED"
	BedStatusCleaning    BedStatus = "CLEANING"
	BedStatusMaintenance BedStatus = "MAINTENANCE"
)

func (e BedStatus) IsValid() bool {
	switch e {
	case BedStatusAvailable,
		BedStatusReserved,
		BedStatusOccupied,
		BedStatusCleaning,
		BedStatusMaintenance:
		return true
	default:
		return false
	}
}

func (e BedStatus) String() string {
	return string(e)
}
