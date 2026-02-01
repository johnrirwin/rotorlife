package models

import (
	"encoding/json"
	"time"
)

// BatteryChemistry represents the type of battery chemistry
type BatteryChemistry string

const (
	ChemistryLIPO   BatteryChemistry = "LIPO"
	ChemistryLIPOHV BatteryChemistry = "LIPO_HV"
	ChemistryLIION  BatteryChemistry = "LIION"
)

// ValidChemistries returns all valid battery chemistries
func ValidChemistries() []BatteryChemistry {
	return []BatteryChemistry{ChemistryLIPO, ChemistryLIPOHV, ChemistryLIION}
}

// IsValidChemistry checks if a chemistry value is valid
func IsValidChemistry(c BatteryChemistry) bool {
	for _, valid := range ValidChemistries() {
		if c == valid {
			return true
		}
	}
	return false
}

// Battery represents a user's battery
type Battery struct {
	ID           string           `json:"id"`
	UserID       string           `json:"userId,omitempty"`
	BatteryCode  string           `json:"batteryCode"`     // Human-friendly ID like "BAT-A1B2"
	Name         string           `json:"name,omitempty"`  // Optional friendly name
	Chemistry    BatteryChemistry `json:"chemistry"`
	Cells        int              `json:"cells"`        // 1S through 8S
	CapacityMah  int              `json:"capacityMah"`
	CRating      *int             `json:"cRating,omitempty"`
	Connector    string           `json:"connector,omitempty"` // e.g., XT30, XT60
	PurchaseDate *time.Time       `json:"purchaseDate,omitempty"`
	Notes        string           `json:"notes,omitempty"`
	CreatedAt    time.Time        `json:"createdAt"`
	UpdatedAt    time.Time        `json:"updatedAt"`

	// Computed fields (populated on detail fetch)
	TotalCycles    int        `json:"totalCycles,omitempty"`
	LastLoggedDate *time.Time `json:"lastLoggedDate,omitempty"`
}

// BatteryLog represents a health/usage log entry for a battery
type BatteryLog struct {
	ID            string          `json:"id"`
	BatteryID     string          `json:"batteryId"`
	UserID        string          `json:"userId,omitempty"`
	LoggedAt      time.Time       `json:"loggedAt"`
	CycleDelta    int             `json:"cycleDelta,omitempty"`     // Usually 1, but can be more
	IRMohmPerCell json.RawMessage `json:"irMohmPerCell,omitempty"`  // JSON array of IR values per cell
	MinCellV      *float64        `json:"minCellV,omitempty"`       // Min cell voltage observed
	MaxCellV      *float64        `json:"maxCellV,omitempty"`       // Max cell voltage observed
	StorageOk     *bool           `json:"storageOk,omitempty"`      // Was it stored at storage voltage?
	Notes         string          `json:"notes,omitempty"`
	CreatedAt     time.Time       `json:"createdAt"`
}

// CreateBatteryParams defines parameters for creating a battery
type CreateBatteryParams struct {
	Name         string           `json:"name,omitempty"`
	Chemistry    BatteryChemistry `json:"chemistry"`
	Cells        int              `json:"cells"`
	CapacityMah  int              `json:"capacityMah"`
	CRating      *int             `json:"cRating,omitempty"`
	Connector    string           `json:"connector,omitempty"`
	PurchaseDate *time.Time       `json:"purchaseDate,omitempty"`
	Notes        string           `json:"notes,omitempty"`
}

// UpdateBatteryParams defines parameters for updating a battery
type UpdateBatteryParams struct {
	ID           string            `json:"id"`
	Name         *string           `json:"name,omitempty"`
	Chemistry    *BatteryChemistry `json:"chemistry,omitempty"`
	Cells        *int              `json:"cells,omitempty"`
	CapacityMah  *int              `json:"capacityMah,omitempty"`
	CRating      *int              `json:"cRating,omitempty"`
	Connector    *string           `json:"connector,omitempty"`
	PurchaseDate *time.Time        `json:"purchaseDate,omitempty"`
	Notes        *string           `json:"notes,omitempty"`
}

// BatteryListParams defines parameters for listing batteries
type BatteryListParams struct {
	Chemistry   BatteryChemistry `json:"chemistry,omitempty"`
	Cells       int              `json:"cells,omitempty"`
	MinCapacity int              `json:"minCapacity,omitempty"`
	MaxCapacity int              `json:"maxCapacity,omitempty"`
	Query       string           `json:"query,omitempty"`
	Sort        string           `json:"sort,omitempty"` // name, updated, logged, cycles
	Limit       int              `json:"limit,omitempty"`
	Offset      int              `json:"offset,omitempty"`
}

// BatteryListResponse represents the response for listing batteries
type BatteryListResponse struct {
	Batteries  []Battery `json:"batteries"`
	TotalCount int       `json:"totalCount"`
}

// BatteryDetailsResponse includes battery with logs
type BatteryDetailsResponse struct {
	Battery
	Logs []BatteryLog `json:"logs"`
}

// CreateBatteryLogParams defines parameters for creating a log entry
type CreateBatteryLogParams struct {
	BatteryID     string          `json:"batteryId"`
	LoggedAt      *time.Time      `json:"loggedAt,omitempty"` // Defaults to now
	CycleDelta    int             `json:"cycleDelta,omitempty"`
	IRMohmPerCell json.RawMessage `json:"irMohmPerCell,omitempty"`
	MinCellV      *float64        `json:"minCellV,omitempty"`
	MaxCellV      *float64        `json:"maxCellV,omitempty"`
	StorageOk     *bool           `json:"storageOk,omitempty"`
	Notes         string          `json:"notes,omitempty"`
}

// BatteryLogListResponse represents the response for listing logs
type BatteryLogListResponse struct {
	Logs       []BatteryLog `json:"logs"`
	TotalCount int          `json:"totalCount"`
}

// LabelSize represents the size of a printed label
type LabelSize string

const (
	LabelSizeSmall    LabelSize = "small"
	LabelSizeStandard LabelSize = "standard"
)
