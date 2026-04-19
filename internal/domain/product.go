package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ProductType represents the type of product.
type ProductType string

const (
	ProductTypeBateria   ProductType = "bateria"
	ProductTypeAccesorio ProductType = "accesorio"
)

// BatteryType represents the battery chemistry type.
type BatteryType string

const (
	BatteryTypeSeca    BatteryType = "seca"
	BatteryTypeLiquida BatteryType = "liquida"
)

// Polarity represents the battery terminal polarity.
type Polarity string

const (
	PolarityIzquierda Polarity = "izquierda"
	PolarityDerecha   Polarity = "derecha"
)

// VehicleType represents the type of vehicle the battery is for.
type VehicleType string

const (
	VehicleTypeAuto VehicleType = "auto"
	VehicleTypeMoto VehicleType = "moto"
	VehicleTypeOtro VehicleType = "otro"
)

// Product represents a product in the system.
type Product struct {
	ID            uuid.UUID    `json:"id"`
	Name          string       `json:"name"`
	Description   *string      `json:"description,omitempty"`
	ProductType   ProductType  `json:"product_type"`
	Brand         *string      `json:"brand,omitempty"`
	Model         *string      `json:"model,omitempty"`
	Voltage       *float64     `json:"voltage,omitempty"`
	Amperage      *float64     `json:"amperage,omitempty"`
	BatteryType   *BatteryType `json:"battery_type,omitempty"`
	Polarity      *Polarity    `json:"polarity,omitempty"`
	AcidLiters    *float64     `json:"acid_liters,omitempty"`
	VehicleType   *VehicleType `json:"vehicle_type,omitempty"`
	MinSalePrice  float64      `json:"min_sale_price"`
	EffectiveDate *time.Time   `json:"effective_date,omitempty"`
	PreviousPrice *float64     `json:"previous_price,omitempty"`
	Active        bool         `json:"active"`
	CreatedAt     time.Time    `json:"created_at"`
	CreatedBy     uuid.UUID    `json:"created_by"`
}

// IsValid validates the product entity according to business rules.
func (p Product) IsValid() error {
	if p.Name == "" {
		return errors.New("name is required")
	}

	if p.MinSalePrice <= 0 {
		return errors.New("min_sale_price must be positive")
	}

	switch p.ProductType {
	case ProductTypeBateria:
		if p.Brand == nil || *p.Brand == "" {
			return errors.New("brand is required for batteries")
		}
		if p.Model == nil || *p.Model == "" {
			return errors.New("model is required for batteries")
		}
		if p.BatteryType == nil {
			return errors.New("battery_type is required for batteries")
		}
		if p.VehicleType == nil {
			return errors.New("vehicle_type is required for batteries")
		}
		if *p.BatteryType == BatteryTypeLiquida {
			if p.AcidLiters == nil || *p.AcidLiters <= 0 {
				return errors.New("acid_liters required for liquid batteries")
			}
		}
		if *p.BatteryType == BatteryTypeSeca && p.AcidLiters != nil {
			return errors.New("acid_liters not allowed for dry batteries")
		}
	case ProductTypeAccesorio:
		if p.Brand != nil || p.Model != nil || p.Voltage != nil || p.AcidLiters != nil || p.VehicleType != nil {
			return errors.New("battery-specific fields not allowed for accessories")
		}
	default:
		return fmt.Errorf("invalid product_type: %s", p.ProductType)
	}

	return nil
}

// IsBattery returns true if the product is a battery.
func (p Product) IsBattery() bool {
	return p.ProductType == ProductTypeBateria
}

// IsAccessory returns true if the product is an accessory.
func (p Product) IsAccessory() bool {
	return p.ProductType == ProductTypeAccesorio
}
