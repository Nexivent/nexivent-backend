package factiliza

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
    baseURL = "https://api.factiliza.com/v1"
)

type FactilizaService struct {
    token      string
    httpClient *http.Client
}

type DNIResponse struct {
    Status  int    `json:"status"`
    Success bool   `json:"success"`
    Message string `json:"message"`
    Data    struct {
        Numero              string `json:"numero"`
        Nombres             string `json:"nombres"`
        ApellidoPaterno     string `json:"apellido_paterno"`
        ApellidoMaterno     string `json:"apellido_materno"`
        NombreCompleto      string `json:"nombre_completo"`
        Departamento        string `json:"departamento"`
        Provincia           string `json:"provincia"`
        Distrito            string `json:"distrito"`
        DireccionCompleta   string `json:"direccion_completa"`
        UbigeoReniec        string `json:"ubigeo_reniec"`
        UbigeoSunat         string `json:"ubigeo_sunat"`
        Ubigeo              []string `json:"ubigeo"`
        FechaNacimiento     string `json:"fecha_nacimiento"`
    } `json:"data"`
}

type CEResponse struct {
    Status  int    `json:"status"`
    Success bool   `json:"success"`
    Message string `json:"message"`
    Data    struct {
        Numero          string `json:"numero"`
        Nombres         string `json:"nombres"`
        ApellidoPaterno string `json:"apellido_paterno"`
        ApellidoMaterno string `json:"apellido_materno"`
        NombreCompleto  string `json:"nombre_completo"`
    } `json:"data"`
}

type RUCResponse struct {
    Status  int    `json:"status"`
    Success bool   `json:"success"`
    Message string   `json:"message"`
    Data    struct {
        Numero                 string   `json:"numero"`
        NombreORazonSocial     string   `json:"nombre_o_razon_social"`
        TipoContribuyente      string   `json:"tipo_contribuyente"`
        Estado                 string   `json:"estado"`
        Condicion              string   `json:"condicion"`
        Direccion              string   `json:"direccion"`
        DireccionCompleta      string   `json:"direccion_completa"`
        Departamento           string   `json:"departamento"`
        Provincia              string   `json:"provincia"`
        Distrito               string   `json:"distrito"`
        UbigeoSunat            string   `json:"ubigeo_sunat"`
        Ubigeo                 []string `json:"ubigeo"`
    } `json:"data"`
}

type RUCRepresentanteResponse struct {
    Status  int    `json:"status"`
    Success bool   `json:"success"`
    Message string `json:"message"`
    Data    []struct {
        TipoDocumento   string `json:"tipo_documento"`
        NumeroDocumento string `json:"numero_documento"`
        Nombre          string `json:"nombre"`
        Cargo           string `json:"cargo"`
        FechaDesde      string `json:"fecha_desde"`
        Estado          string `json:"estado"`
        Condicion       string `json:"condicion"`
        Departamento    string `json:"departamento"`
        Provincia       string `json:"provincia"`
        Distrito        string `json:"distrito"`
        Direccion       string `json:"direccion"`
    } `json:"data"`
}

func NewFactilizaService(token string) *FactilizaService {
    return &FactilizaService{
        token: token,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func makeRequest[T any](s *FactilizaService, url string) (*T, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    // Según la documentación, el formato del header es: Authorization: Bearer <token>
    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.token))

    resp, err := s.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error making request: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response: %w", err)
    }

    // Log para debugging (opcional, quitar en producción)
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
    }

    var response T
    if err := json.Unmarshal(body, &response); err != nil {
        return nil, fmt.Errorf("error parsing response: %w. Body: %s", err, string(body))
    }

    return &response, nil
}

// ConsultarDNI consulta información de un DNI
// Endpoint: GET /v1/dni/info/{dni}
func (s *FactilizaService) ConsultarDNI(dni string) (*DNIResponse, error) {
    url := fmt.Sprintf("%s/dni/info/%s", baseURL, dni)
    return makeRequest[DNIResponse](s, url)
}

// ConsultarCE consulta información de un Carnet de Extranjería
// Endpoint: GET /v1/cee/info/{cee}
func (s *FactilizaService) ConsultarCE(ce string) (*CEResponse, error) {
    // ⚠️ CORRECCIÓN: La documentación usa /cee/info/{cee}, no /ce/info/{ce}
    url := fmt.Sprintf("%s/cee/info/%s", baseURL, ce)
    return makeRequest[CEResponse](s, url)
}

// ConsultarRUC consulta información de un RUC (empresa)
// Endpoint: GET /v1/ruc/info/{ruc}
func (s *FactilizaService) ConsultarRUC(ruc string) (*RUCResponse, error) {
    url := fmt.Sprintf("%s/ruc/info/%s", baseURL, ruc)
    return makeRequest[RUCResponse](s, url)
}

// ConsultarRUCRepresentante consulta representantes de un RUC (persona natural)
// Endpoint: GET /v1/ruc/representante/{ruc}
func (s *FactilizaService) ConsultarRUCRepresentante(ruc string) (*RUCRepresentanteResponse, error) {
    url := fmt.Sprintf("%s/ruc/representante/%s", baseURL, ruc)
    return makeRequest[RUCRepresentanteResponse](s, url)
}

