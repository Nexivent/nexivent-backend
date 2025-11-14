package convert

// MapEstadoToInt16 converts string estado to int16
func MapEstadoToInt16(estado string) int16 {
	switch estado {
	case "BORRADOR":
		return 0
	case "PUBLICADO":
		return 1
	case "FINALIZADO":
		return 2
	default:
		return 0
	}
}

// MapEstadoToString converts int16 estado to string
func MapEstadoToString(estado int16) string {
	switch estado {
	case 0:
		return "BORRADOR"
	case 1:
		return "PUBLICADO"
	case 2:
		return "FINALIZADO"
	default:
		return "BORRADOR"
	}
}
