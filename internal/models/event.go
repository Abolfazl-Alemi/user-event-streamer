package models

type GA4Event struct {
	EventName       string              `json:"event_name"`
	EventParams     []EventParam        `json:"event_params"`
	UserProperties  map[string]Property `json:"user_properties,omitempty"`
	UserID          string              `json:"user_id,omitempty"`
	SessionID       int64               `json:"session_id,omitempty"`
	TimestampMicros int64               `json:"timestamp_micros"`
	Geo             *GeoInfo            `json:"geo,omitempty"`
	Device          *DeviceInfo         `json:"device,omitempty"`
	TrafficSource   *TrafficSource      `json:"traffic_source,omitempty"`
	Ecommerce       *EcommerceInfo      `json:"ecommerce,omitempty"`
}

type EventParam struct {
	Key   string     `json:"key"`
	Value ParamValue `json:"value"`
}

type ParamValue struct {
	StringValue string  `json:"string_value,omitempty"`
	IntValue    int64   `json:"int_value,omitempty"`
	FloatValue  float32 `json:"float_value,omitempty"`
	DoubleValue float64 `json:"double_value,omitempty"`
}

type Property struct {
	Value string `json:"value"`
}

type GeoInfo struct {
	Country string `json:"country,omitempty"`
	Region  string `json:"region,omitempty"`
	City    string `json:"city,omitempty"`
}

type DeviceInfo struct {
	DeviceCategory  string `json:"device_category,omitempty"`
	MobileBrandName string `json:"mobile_brand_name,omitempty"`
	MobileModelName string `json:"mobile_model_name,omitempty"`
	Platform        string `json:"platform,omitempty"`
}

type TrafficSource struct {
	Medium   string `json:"medium,omitempty"`
	Source   string `json:"source,omitempty"`
	Campaign string `json:"campaign,omitempty"`
}

type EcommerceInfo struct {
	TransactionID string          `json:"transaction_id,omitempty"`
	Value         float64         `json:"value,omitempty"`
	Currency      string          `json:"currency,omitempty"`
	Items         []EcommerceItem `json:"items,omitempty"`
}

type EcommerceItem struct {
	ItemID   string  `json:"item_id,omitempty"`
	ItemName string  `json:"item_name,omitempty"`
	Price    float64 `json:"price,omitempty"`
	Quantity int     `json:"quantity,omitempty"`
}
