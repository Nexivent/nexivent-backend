package schemas

// Item de entrada dentro del hold
// "entradas": [ { "idTarifa": "", "cantidad": "" } ]
type EntradaOrdenRequest struct {
	IdTarifa int64 `json:"idTarifa"`
	IdPerfil int64 `json:"idPerfil"`
	IdSector int64 `json:"idSector"`
	Cantidad int64 `json:"cantidad"`
}

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

// Request:
// { "paymentId": "" }
type ConfirmarOrdenRequest struct {
	PaymentID   string `json:"paymentId"`
	IdEvento    int64  `json:"idEvento"`
	FechaEvento string `json:"fechaEvento"` // "YYYY-MM-DD"
}

// Response 200:
// { "orderId": "", "estado": "CONFIRMADA", "mensaje": "Compra confirmada" }
type ConfirmarOrdenResponse struct {
	OrderID int64  `json:"orderId"`
	Estado  string `json:"estado"`  // "CONFIRMADA"
	Mensaje string `json:"mensaje"` // "Compra confirmada"
}
