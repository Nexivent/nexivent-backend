package adapter

import (
	"fmt"
	"strings"

	"github.com/Nexivent/nexivent-backend/internal/application/service/factiliza"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
)

type ValidacionDocumento struct {
	logger           logging.Logger
	factilizaService *factiliza.FactilizaService
}

func NewValidacionDocumentoAdapter(
	logger logging.Logger,
	factilizaToken string,
) *ValidacionDocumento {
	return &ValidacionDocumento{
		logger:           logger,
		factilizaService: factiliza.NewFactilizaService(factilizaToken),
	}
}

// ValidarDocumento valida un documento (DNI o RUC) usando la API de Factiliza
func (v *ValidacionDocumento) ValidarDocumento(req *schemas.ValidarDocumentoRequest) (*schemas.ValidarDocumentoResponse, error) {
	v.logger.Info("Validando documento:", req.TipoDocumento, req.NumeroDocumento)

	switch req.TipoDocumento {
	case "DNI":
        return v.validarDNI(req.NumeroDocumento)
    case "CE":
        return v.validarCE(req.NumeroDocumento)
    case "RUC_PERSONA":
        return v.validarRUCPersona(req.NumeroDocumento)
    case "RUC_EMPRESA":
        return v.validarRUCEmpresa(req.NumeroDocumento)
    default:
        return &schemas.ValidarDocumentoResponse{
            Success: false,
            Message: "Tipo de documento no válido. Use: DNI, CE, RUC_PERSONA o RUC_EMPRESA",
        }, nil
	}
}

func (v *ValidacionDocumento) validarDNI(dni string) (*schemas.ValidarDocumentoResponse, error) {
	if len(dni) != 8 {
        return &schemas.ValidarDocumentoResponse{
            Success: false,
            Message: "El DNI debe tener 8 dígitos",
        }, nil
    }

    v.logger.Info("Consultando DNI:", dni)
    dniData, err := v.factilizaService.ConsultarDNI(dni)
    if err != nil {
        v.logger.Error("Error consultando DNI:", err)
        return &schemas.ValidarDocumentoResponse{
            Success: false,
            Message: "Error al consultar el DNI",
        }, err
    }

    if !dniData.Success {
        return &schemas.ValidarDocumentoResponse{
            Success: false,
            Message: dniData.Message,
        }, nil
    }

    // Unir nombres y apellidos
    nombreCompleto := strings.TrimSpace(fmt.Sprintf("%s %s %s",
        dniData.Data.Nombres,
        dniData.Data.ApellidoPaterno,
        dniData.Data.ApellidoMaterno,
    ))

    return &schemas.ValidarDocumentoResponse{
        Success: true,
        Message: "DNI validado correctamente",
        Data: &schemas.ValidacionDocumentoData{
            TipoDocumento:   "DNI",
            NumeroDocumento: dni,
            NombreCompleto:  nombreCompleto,
            Direccion:       dniData.Data.DireccionCompleta,
            Departamento:    dniData.Data.Departamento,
            Provincia:       dniData.Data.Provincia,
            Distrito:        dniData.Data.Distrito,
            Ubigeo:          dniData.Data.UbigeoSunat,
            EsEmpresa:       false,
            Valido:          true,
        },
    }, nil
}

func (a *ValidacionDocumento) validarCE(ce string) (*schemas.ValidarDocumentoResponse, error) {
    if len(ce) != 9 && len(ce) != 12 {
        return &schemas.ValidarDocumentoResponse{
            Success: false,
            Message: "El Carnet de Extranjería debe tener 9 o 12 caracteres",
        }, nil
    }

    a.logger.Info("Consultando CE:", ce)
    ceData, err := a.factilizaService.ConsultarCE(ce)
    if err != nil {
        a.logger.Error("Error consultando CE:", err)
        return &schemas.ValidarDocumentoResponse{
            Success: false,
            Message: "Error al consultar el Carnet de Extranjería",
        }, err
    }

    if !ceData.Success {
        return &schemas.ValidarDocumentoResponse{
            Success: false,
            Message: ceData.Message,
        }, nil
    }

    nombreCompleto := strings.TrimSpace(fmt.Sprintf("%s %s %s",
        ceData.Data.Nombres,
        ceData.Data.ApellidoPaterno,
        ceData.Data.ApellidoMaterno,
    ))

    return &schemas.ValidarDocumentoResponse{
        Success: true,
        Message: "Carnet de Extranjería validado correctamente",
        Data: &schemas.ValidacionDocumentoData{
            TipoDocumento:   "CE",
            NumeroDocumento: ce,
            NombreCompleto:  nombreCompleto,
            EsEmpresa:       false,
            Valido:          true,
        },
    }, nil
}

func (v *ValidacionDocumento) validarRUCPersona(ruc string) (*schemas.ValidarDocumentoResponse, error) {
    if len(ruc) != 11 {
        return &schemas.ValidarDocumentoResponse{
            Success: false,
            Message: "El RUC debe tener 11 dígitos",
        }, nil
    }

    // RUC de persona natural empieza con 10
    if !strings.HasPrefix(ruc, "10") {
        return &schemas.ValidarDocumentoResponse{
            Success: false,
            Message: "El RUC de persona natural debe empezar con 10",
        }, nil
    }

    v.logger.Info("Consultando RUC Persona (Representante):", ruc)
    rucData, err := v.factilizaService.ConsultarRUCRepresentante(ruc)
    if err != nil {
        v.logger.Error("Error consultando RUC Representante:", err)
        return &schemas.ValidarDocumentoResponse{
            Success: false,
            Message: "Error al consultar el RUC",
        }, err
    }

    if !rucData.Success || len(rucData.Data) == 0 {
        return &schemas.ValidarDocumentoResponse{
            Success: false,
            Message: rucData.Message,
        }, nil
    }

    // Usar el primer representante
    representante := rucData.Data[0]

    return &schemas.ValidarDocumentoResponse{
        Success: true,
        Message: "RUC de persona natural validado correctamente",
        Data: &schemas.ValidacionDocumentoData{
            TipoDocumento:          "RUC_PERSONA",
            NumeroDocumento:        ruc,
            NombreCompleto:         representante.Nombre,
            Direccion:              representante.Direccion,
            Departamento:           representante.Departamento,
            Provincia:              representante.Provincia,
            Distrito:               representante.Distrito,
            EstadoContribuyente:    representante.Estado,
            CondicionContribuyente: representante.Condicion,
            EsEmpresa:              false,
            Valido:                 true,
        },
    }, nil
}

func (a *ValidacionDocumento) validarRUCEmpresa(ruc string) (*schemas.ValidarDocumentoResponse, error) {
    if len(ruc) != 11 {
        return &schemas.ValidarDocumentoResponse{
            Success: false,
            Message: "El RUC debe tener 11 dígitos",
        }, nil
    }

    // RUC de empresa empieza con 20
    if !strings.HasPrefix(ruc, "20") {
        return &schemas.ValidarDocumentoResponse{
            Success: false,
            Message: "El RUC de empresa debe empezar con 20",
        }, nil
    }

    a.logger.Info("Consultando RUC Empresa:", ruc)
    rucData, err := a.factilizaService.ConsultarRUC(ruc)
    if err != nil {
        a.logger.Error("Error consultando RUC:", err)
        return &schemas.ValidarDocumentoResponse{
            Success: false,
            Message: "Error al consultar el RUC",
        }, err
    }

    if !rucData.Success {
        return &schemas.ValidarDocumentoResponse{
            Success: false,
            Message: rucData.Message,
        }, nil
    }

    return &schemas.ValidarDocumentoResponse{
        Success: true,
        Message: "RUC de empresa validado correctamente",
        Data: &schemas.ValidacionDocumentoData{
            TipoDocumento:          "RUC_EMPRESA",
            NumeroDocumento:        ruc,
            RazonSocial:            rucData.Data.NombreORazonSocial,
            Direccion:              rucData.Data.DireccionCompleta,
            Departamento:           rucData.Data.Departamento,
            Provincia:              rucData.Data.Provincia,
            Distrito:               rucData.Data.Distrito,
            Ubigeo:                 rucData.Data.UbigeoSunat,
            EstadoContribuyente:    rucData.Data.Estado,
            CondicionContribuyente: rucData.Data.Condicion,
            EsEmpresa:              true,
            Valido:                 true,
        },
    }, nil
}