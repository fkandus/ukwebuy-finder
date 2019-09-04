package main

// StoresResponse represents the topmost level
type StoresResponse struct {
	Response CexStoresResponse
}

// CexStoresResponse is the Response to nearest stores request
type CexStoresResponse struct {
	Data StoresDataResponse
}

// StoresDataResponse is the list of nearest stores
type StoresDataResponse struct {
	NearestStores []NearestStoresResponse
}

// NearestStoresResponse represents the important data of a store
type NearestStoresResponse struct {
	StoreName      string
	QuantityOnHand interface{}
}

// DetailResponse represents the topmost level
type DetailResponse struct {
	Response CexDetailResponse
}

// CexDetailResponse is the Response to nearest stores request
type CexDetailResponse struct {
	Data DetailDataResponse
}

// DetailDataResponse is the list of nearest stores
type DetailDataResponse struct {
	BoxDetails []ItemDetailResponse
}

// ItemDetailResponse represents the important data of a store
type ItemDetailResponse struct {
	BoxName       string
	SellPrice     float64
	ExchangePrice float64
}
