package schemas

// =============== ENTRADAS / ITEMS DE LA ORDEN ===============

// Item de entrada dentro del hold
// "entradas": [ { "idTarifa": "", "cantidad": "" } ]
type EntradaOrdenRequest struct {
	IdTarifa int64 `json:"idTarifa"`
	Cantidad int64 `json:"cantidad"`
}

// =============== POST /api/orders/hold ===============

// Request:
// {
//   "idEvento": "",
//   "idFechaEvento": "",
//   "idUsuario": "",
//   "entradas": [
//     { "idTarifa": "", "cantidad": "" }
//   ]
// }
type CrearOrdenTemporalRequest struct {
	IdEvento      int64                 `json:"idEvento"`
	IdFechaEvento int64                 `json:"idFechaEvento"`
	IdUsuario     int64                 `json:"idUsuario"`
	Entradas      []EntradaOrdenRequest `json:"entradas"`
}

// Response 201:
// {
//   "orderId": "",
//   "estado": "TEMPORAL",
//   "total": "",
//   "startedAt": "",
//   "expiresAt": "",
//   "ttlSeconds": ""
// }
type CrearOrdenTemporalResponse struct {
	OrderID    int64   `json:"orderId"`
	Estado     string  `json:"estado"` // "TEMPORAL"
	Total      float64 `json:"total"`
	StartedAt  string  `json:"startedAt"`  // RFC3339
	ExpiresAt  string  `json:"expiresAt"`  // RFC3339
	TTLSeconds int64   `json:"ttlSeconds"` // segundos
}

// =============== GET /api/orders/{orderId}/hold ===============

// Response 200:
// {
//   "orderId": "",
//   "estado": "BORRADOR",
//   "remainingSeconds": "",
//   "startedAt": "",
//   "expiresAt": "",
//   "total": ""
// }
type ObtenerHoldResponse struct {
	OrderID       int64   `json:"orderId"`
	Estado        string  `json:"estado"` // "BORRADOR" mientras esté TEMPORAL
	RemainingSecs int64   `json:"remainingSeconds"`
	StartedAt     string  `json:"startedAt"` // RFC3339
	ExpiresAt     string  `json:"expiresAt"` // RFC3339
	Total         float64 `json:"total"`
}

// =============== POST /api/orders/{orderId}/confirm ===============

// Request:
// { "paymentId": "" }
type ConfirmarOrdenRequest struct {
	PaymentID string `json:"paymentId"`
}

// Response 200:
// { "orderId": "", "estado": "CONFIRMADA", "mensaje": "Compra confirmada" }
type ConfirmarOrdenResponse struct {
	OrderID int64  `json:"orderId"`
	Estado  string `json:"estado"`  // "CONFIRMADA"
	Mensaje string `json:"mensaje"` // "Compra confirmada"
}
