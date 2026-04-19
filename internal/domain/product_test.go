package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProduct_IsValid(t *testing.T) {
	now := time.Now()
	brand := "Yuasa"
	model := "YTX9-BS"
	voltage := 12.0
	amperage := 9.0
	batteryTypeSeca := BatteryTypeSeca
	batteryTypeLiquida := BatteryTypeLiquida
	polarity := PolarityDerecha
	acidLiters := 0.5
	vehicleTypeAuto := VehicleTypeAuto
	vehicleTypeMoto := VehicleTypeMoto
	description := "Batería de alta calidad"

	tests := []struct {
		name    string
		product Product
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid battery type seca",
			product: Product{
				ID:           uuid.New(),
				Name:         "Batería Yuasa YTX9-BS",
				Description:  &description,
				ProductType:  ProductTypeBateria,
				Brand:        &brand,
				Model:        &model,
				Voltage:      &voltage,
				Amperage:     &amperage,
				BatteryType:  &batteryTypeSeca,
				Polarity:     &polarity,
				VehicleType:  &vehicleTypeAuto,
				MinSalePrice: 150.00,
				Active:       true,
				CreatedAt:    now,
				CreatedBy:    uuid.New(),
			},
			wantErr: false,
		},
		{
			name: "valid battery type liquida",
			product: Product{
				ID:           uuid.New(),
				Name:         "Batería Líquida 12N9",
				ProductType:  ProductTypeBateria,
				Brand:        &brand,
				Model:        &model,
				Voltage:      &voltage,
				Amperage:     &amperage,
				BatteryType:  &batteryTypeLiquida,
				Polarity:     &polarity,
				AcidLiters:   &acidLiters,
				VehicleType:  &vehicleTypeMoto,
				MinSalePrice: 120.00,
				Active:       true,
				CreatedAt:    now,
				CreatedBy:    uuid.New(),
			},
			wantErr: false,
		},
		{
			name: "valid accessory",
			product: Product{
				ID:           uuid.New(),
				Name:         "Cargador de batería",
				Description:  &description,
				ProductType:  ProductTypeAccesorio,
				MinSalePrice: 45.00,
				Active:       true,
				CreatedAt:    now,
				CreatedBy:    uuid.New(),
			},
			wantErr: false,
		},
		{
			name: "invalid - empty name",
			product: Product{
				ID:           uuid.New(),
				Name:         "",
				ProductType:  ProductTypeBateria,
				Brand:        &brand,
				Model:        &model,
				BatteryType:  &batteryTypeSeca,
				VehicleType:  &vehicleTypeAuto,
				MinSalePrice: 150.00,
				Active:       true,
				CreatedAt:    now,
				CreatedBy:    uuid.New(),
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "invalid - negative price",
			product: Product{
				ID:           uuid.New(),
				Name:         "Batería",
				ProductType:  ProductTypeBateria,
				Brand:        &brand,
				Model:        &model,
				BatteryType:  &batteryTypeSeca,
				VehicleType:  &vehicleTypeAuto,
				MinSalePrice: -10.00,
				Active:       true,
				CreatedAt:    now,
				CreatedBy:    uuid.New(),
			},
			wantErr: true,
			errMsg:  "min_sale_price must be positive",
		},
		{
			name: "invalid - zero price",
			product: Product{
				ID:           uuid.New(),
				Name:         "Batería",
				ProductType:  ProductTypeBateria,
				Brand:        &brand,
				Model:        &model,
				BatteryType:  &batteryTypeSeca,
				VehicleType:  &vehicleTypeAuto,
				MinSalePrice: 0,
				Active:       true,
				CreatedAt:    now,
				CreatedBy:    uuid.New(),
			},
			wantErr: true,
			errMsg:  "min_sale_price must be positive",
		},
		{
			name: "invalid - battery missing brand",
			product: Product{
				ID:           uuid.New(),
				Name:         "Batería",
				ProductType:  ProductTypeBateria,
				Brand:        nil,
				Model:        &model,
				BatteryType:  &batteryTypeSeca,
				VehicleType:  &vehicleTypeAuto,
				MinSalePrice: 150.00,
				Active:       true,
				CreatedAt:    now,
				CreatedBy:    uuid.New(),
			},
			wantErr: true,
			errMsg:  "brand is required for batteries",
		},
		{
			name: "invalid - battery empty brand",
			product: Product{
				ID:           uuid.New(),
				Name:         "Batería",
				ProductType:  ProductTypeBateria,
				Brand:        strPtr(""),
				Model:        &model,
				BatteryType:  &batteryTypeSeca,
				VehicleType:  &vehicleTypeAuto,
				MinSalePrice: 150.00,
				Active:       true,
				CreatedAt:    now,
				CreatedBy:    uuid.New(),
			},
			wantErr: true,
			errMsg:  "brand is required for batteries",
		},
		{
			name: "invalid - battery missing model",
			product: Product{
				ID:           uuid.New(),
				Name:         "Batería",
				ProductType:  ProductTypeBateria,
				Brand:        &brand,
				Model:        nil,
				BatteryType:  &batteryTypeSeca,
				VehicleType:  &vehicleTypeAuto,
				MinSalePrice: 150.00,
				Active:       true,
				CreatedAt:    now,
				CreatedBy:    uuid.New(),
			},
			wantErr: true,
			errMsg:  "model is required for batteries",
		},
		{
			name: "invalid - battery missing battery_type",
			product: Product{
				ID:           uuid.New(),
				Name:         "Batería",
				ProductType:  ProductTypeBateria,
				Brand:        &brand,
				Model:        &model,
				BatteryType:  nil,
				VehicleType:  &vehicleTypeAuto,
				MinSalePrice: 150.00,
				Active:       true,
				CreatedAt:    now,
				CreatedBy:    uuid.New(),
			},
			wantErr: true,
			errMsg:  "battery_type is required for batteries",
		},
		{
			name: "invalid - battery missing vehicle_type",
			product: Product{
				ID:           uuid.New(),
				Name:         "Batería",
				ProductType:  ProductTypeBateria,
				Brand:        &brand,
				Model:        &model,
				BatteryType:  &batteryTypeSeca,
				VehicleType:  nil,
				MinSalePrice: 150.00,
				Active:       true,
				CreatedAt:    now,
				CreatedBy:    uuid.New(),
			},
			wantErr: true,
			errMsg:  "vehicle_type is required for batteries",
		},
		{
			name: "invalid - liquida battery missing acid_liters",
			product: Product{
				ID:           uuid.New(),
				Name:         "Batería Líquida",
				ProductType:  ProductTypeBateria,
				Brand:        &brand,
				Model:        &model,
				BatteryType:  &batteryTypeLiquida,
				AcidLiters:   nil,
				VehicleType:  &vehicleTypeAuto,
				MinSalePrice: 150.00,
				Active:       true,
				CreatedAt:    now,
				CreatedBy:    uuid.New(),
			},
			wantErr: true,
			errMsg:  "acid_liters required for liquid batteries",
		},
		{
			name: "invalid - liquida battery zero acid_liters",
			product: Product{
				ID:           uuid.New(),
				Name:         "Batería Líquida",
				ProductType:  ProductTypeBateria,
				Brand:        &brand,
				Model:        &model,
				BatteryType:  &batteryTypeLiquida,
				AcidLiters:   floatPtr(0),
				VehicleType:  &vehicleTypeAuto,
				MinSalePrice: 150.00,
				Active:       true,
				CreatedAt:    now,
				CreatedBy:    uuid.New(),
			},
			wantErr: true,
			errMsg:  "acid_liters required for liquid batteries",
		},
		{
			name: "invalid - seca battery with acid_liters",
			product: Product{
				ID:           uuid.New(),
				Name:         "Batería Seca",
				ProductType:  ProductTypeBateria,
				Brand:        &brand,
				Model:        &model,
				BatteryType:  &batteryTypeSeca,
				AcidLiters:   &acidLiters,
				VehicleType:  &vehicleTypeAuto,
				MinSalePrice: 150.00,
				Active:       true,
				CreatedAt:    now,
				CreatedBy:    uuid.New(),
			},
			wantErr: true,
			errMsg:  "acid_liters not allowed for dry batteries",
		},
		{
			name: "invalid - accessory with battery fields",
			product: Product{
				ID:           uuid.New(),
				Name:         "Cargador",
				ProductType:  ProductTypeAccesorio,
				Brand:        &brand,
				MinSalePrice: 45.00,
				Active:       true,
				CreatedAt:    now,
				CreatedBy:    uuid.New(),
			},
			wantErr: true,
			errMsg:  "battery-specific fields not allowed for accessories",
		},
		{
			name: "invalid - unknown product_type",
			product: Product{
				ID:           uuid.New(),
				Name:         "Producto",
				ProductType:  ProductType("unknown"),
				MinSalePrice: 100.00,
				Active:       true,
				CreatedAt:    now,
				CreatedBy:    uuid.New(),
			},
			wantErr: true,
			errMsg:  "invalid product_type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.IsValid()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestProduct_IsBattery(t *testing.T) {
	tests := []struct {
		name     string
		prodType ProductType
		expected bool
	}{
		{"battery is battery", ProductTypeBateria, true},
		{"accessory is not battery", ProductTypeAccesorio, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Product{ProductType: tt.prodType}
			assert.Equal(t, tt.expected, p.IsBattery())
		})
	}
}

func TestProduct_IsAccessory(t *testing.T) {
	tests := []struct {
		name     string
		prodType ProductType
		expected bool
	}{
		{"accessory is accessory", ProductTypeAccesorio, true},
		{"battery is not accessory", ProductTypeBateria, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Product{ProductType: tt.prodType}
			assert.Equal(t, tt.expected, p.IsAccessory())
		})
	}
}

func TestProductType_Constants(t *testing.T) {
	assert.Equal(t, ProductType("bateria"), ProductTypeBateria)
	assert.Equal(t, ProductType("accesorio"), ProductTypeAccesorio)
}

func TestBatteryType_Constants(t *testing.T) {
	assert.Equal(t, BatteryType("seca"), BatteryTypeSeca)
	assert.Equal(t, BatteryType("liquida"), BatteryTypeLiquida)
}

func TestPolarity_Constants(t *testing.T) {
	assert.Equal(t, Polarity("izquierda"), PolarityIzquierda)
	assert.Equal(t, Polarity("derecha"), PolarityDerecha)
}

func TestVehicleType_Constants(t *testing.T) {
	assert.Equal(t, VehicleType("auto"), VehicleTypeAuto)
	assert.Equal(t, VehicleType("moto"), VehicleTypeMoto)
	assert.Equal(t, VehicleType("otro"), VehicleTypeOtro)
}

// Helper functions
func strPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}
