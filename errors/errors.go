package errors

import (
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

type Error struct {
	Code    string
	Message string
}

var (
	// For 404 Not Found errors
	ObjectNotFoundError = struct {
		CommunityNotFound           Error
		ReservationNotFound         Error
		ProfessionalNotFound        Error
		LocalNotFound               Error
		UserNotFound                Error
		ServiceNotFound             Error
		PlanNotFound                Error
		MembershipNotFound          Error
		OnboardingNotFound          Error
		CommunityPlanNotFound       Error
		CommunityServiceNotFound    Error
		ServiceLocalNotFound        Error
		ServiceProfessionalNotFound Error
		SessionNotFound             Error
		EventoNotFound              Error
		CategoriaNotFound           Error
	}{
		CommunityNotFound: Error{
			Code:    "COMMUNITY_ERROR_001",
			Message: "Community not found",
		},
		ProfessionalNotFound: Error{
			Code:    "PROFESSIONAL_ERROR_001",
			Message: "Professional not found",
		},
		LocalNotFound: Error{
			Code:    "LOCAL_ERROR_001",
			Message: "Local not found",
		},
		ServiceNotFound: Error{
			Code:    "SERVICE_ERROR_001",
			Message: "Service not found",
		},
		PlanNotFound: Error{
			Code:    "PLAN_ERROR_001",
			Message: "Plan not found",
		},
		UserNotFound: Error{
			Code:    "USER_ERROR_001",
			Message: "User not found",
		},
		MembershipNotFound: Error{
			Code:    "MEMBERSHIP_ERROR_001",
			Message: "Membership not found",
		},
		OnboardingNotFound: Error{
			Code:    "ONBOARDING_ERROR_001",
			Message: "Onboarding not found",
		},
		ServiceLocalNotFound: Error{
			Code:    "SERVICE_LOCAL_ERROR_001",
			Message: "Service-Local association not found",
		},
		ServiceProfessionalNotFound: Error{
			Code:    "SERVICE_PROFESSIONAL_ERROR_001",
			Message: "Service-Professional association not found",
		},
		SessionNotFound: Error{
			Code:    "SESSION_ERROR_001",
			Message: "Session not found",
		},
		EventoNotFound: Error{
			Code:    "EVENTO_ERROR_001",
			Message: "Evento not found",
		},
		CategoriaNotFound: Error{
			Code:    "CATEGORIA_ERROR_001",
			Message: "Categoria not found",
		},
	}

	// For 422 Unprocessable Entity errors
	UnprocessableEntityError = struct {
		InvalidCommunityId           Error
		InvalidRequestBody           Error
		InvalidProfessionalId        Error
		InvalidServiceId             Error
		InvalidMembershipId          Error
		InvalidOnboardingId          Error
		InvalidUserEmail             Error
		InvalidUserId                Error
		InvalidCommunityPlanId       Error
		InvalidCommunityServiceId    Error
		InvalidParsingInteger        Error
		InvalidServiceLocalId        Error
		InvalidServiceProfessionalId Error
		InvalidSessionId             Error
		InvalidReservationId         Error
		InvalidDateFormat            Error
		EmailAlreadyRegistered       Error
		InvalidEventoId              Error
	}{
		InvalidRequestBody: Error{
			Code:    "REQUEST_ERROR_001",
			Message: "Invalid body request",
		},
		InvalidCommunityId: Error{
			Code:    "COMMUNITY_ERROR_004",
			Message: "Invalid community id",
		},
		InvalidProfessionalId: Error{
			Code:    "PROFESSIONAL_ERROR_004",
			Message: "Invalid professional id",
		},
		InvalidServiceId: Error{
			Code:    "SERVICE_ERROR_004",
			Message: "Invalid service id",
		},
		InvalidMembershipId: Error{
			Code:    "MEMBERSHIP_ERROR_001",
			Message: "Invalid membership id",
		},
		InvalidOnboardingId: Error{
			Code:    "ONBOARDING_ERROR_001",
			Message: "Invalid onboarding id",
		},
		InvalidUserEmail: Error{
			Code:    "USER_ERROR_001",
			Message: "Invalid user email",
		},
		InvalidUserId: Error{
			Code:    "USER_ERROR_004",
			Message: "Invalid user id",
		},
		InvalidCommunityPlanId: Error{
			Code:    "COMMUNITY_PLAN_ERROR_004",
			Message: "Invalid community_id or plan_id for association",
		},
		InvalidCommunityServiceId: Error{
			Code:    "COMMUNITY_SERVICE_ERROR_004",
			Message: "Invalid community_id or service_id for association",
		},
		InvalidParsingInteger: Error{
			Code:    "REQUEST_ERROR_004",
			Message: "Invalid parsing integer",
		},
		InvalidServiceLocalId: Error{
			Code:    "SERVICE_LOCAL_ERROR_004",
			Message: "Invalid service_id or local_id for association",
		},
		InvalidServiceProfessionalId: Error{
			Code:    "SERVICE_PROFESSIONAL_ERROR_004",
			Message: "Invalid service_id or professional_id for association",
		},
		InvalidSessionId: Error{
			Code:    "SESSION_ERROR_004",
			Message: "Invalid session id",
		},
		InvalidReservationId: Error{
			Code:    "RESERVATION_ERROR_004",
			Message: "Invalid reservation id",
		},
		InvalidDateFormat: Error{
			Code:    "DATE_FORMAT_ERROR_004",
			Message: "Invalid date format",
		},
		EmailAlreadyRegistered: Error{
			Code:    "USER_ERROR_006",
			Message: "Email is already registered",
		},
		InvalidEventoId: Error{
			Code:    "EVENTO_ERROR_004",
			Message: "Invalid evento_id",
		},
	}

	// For 400 Bad Request errors
	BadRequestError = struct {
		InvalidUpdatedByValue         Error
		CommunityNotCreated           Error
		CommunityNotUpdated           Error
		CommunityNotSoftDeleted       Error
		LocalNotCreated               Error
		LocalNotUpdated               Error
		LocalNotSoftDeleted           Error
		ProfessionalNotCreated        Error
		ProfessionalNotUpdated        Error
		ProfessionalNotSoftDeleted    Error
		ServiceNotCreated             Error
		ServiceNotUpdated             Error
		ServiceNotSoftDeleted         Error
		PlanNotCreated                Error
		PlanNotUpdated                Error
		PlanNotSoftDeleted            Error
		InvalidPlanType               Error
		MembershipNotCreated          Error
		MembershipNotUpdated          Error
		OnboardingNotCreated          Error
		OnboardingNotUpdated          Error
		UserNotCreated                Error
		UserNotUpdated                Error
		UserNotSoftDeleted            Error
		CommunityPlanNotCreated       Error
		CommunityPlanNotDeleted       Error
		CommunityServiceNotCreated    Error
		CommunityServiceNotDeleted    Error
		ServiceLocalNotCreated        Error
		ServiceLocalNotDeleted        Error
		ServiceProfessionalNotCreated Error
		ServiceProfessionalNotDeleted Error
		SessionNotCreated             Error
		SessionNotUpdated             Error
		SessionNotSoftDeleted         Error
		EventoNotCreated              Error
		EventoNotUpdated              Error
		EventoNotFound                Error
		CategoriaNotCreated           Error
		CategoriaNotUpdated           Error
		CategoriaNotFound             Error
		InvalidUploadMimeType         Error
		UploadURLNotCreated           Error
		InvalidIDParam                Error
	}{
		InvalidUpdatedByValue: Error{
			Code:    "REQUEST_ERROR_002",
			Message: "Invalid updated by value error",
		},
		LocalNotCreated: Error{
			Code:    "LOCAL_ERROR_002",
			Message: "Local not created",
		},
		LocalNotUpdated: Error{
			Code:    "LOCAL_ERROR_003",
			Message: "Local not updated",
		},
		LocalNotSoftDeleted: Error{
			Code:    "LOCAL_ERROR_005",
			Message: "Local not soft deleted",
		},
		MembershipNotCreated: Error{
			Code:    "MEMBERSHIP_ERROR_002",
			Message: "Membership not created",
		},
		MembershipNotUpdated: Error{
			Code:    "MEMBERSHIP_ERROR_003",
			Message: "Membership not updated",
		},
		OnboardingNotCreated: Error{
			Code:    "ONBOARDING_ERROR_002",
			Message: "Onboarding not created",
		},
		OnboardingNotUpdated: Error{
			Code:    "ONBOARDING_ERROR_003",
			Message: "Onboarding not updated",
		},
		UserNotCreated: Error{
			Code:    "USER_ERROR_002",
			Message: "User not created",
		},
		UserNotUpdated: Error{
			Code:    "USER_ERROR_003",
			Message: "User not updated",
		},
		ServiceNotCreated: Error{
			Code:    "SERVICE_ERROR_002",
			Message: "Service not created",
		},
		ServiceNotUpdated: Error{
			Code:    "SERVICE_ERROR_003",
			Message: "Service not updated",
		},
		ServiceNotSoftDeleted: Error{
			Code:    "SERVICE_ERROR_005",
			Message: "Service not soft deleted",
		},
		InvalidPlanType: Error{
			Code:    "PLAN_ERROR_005",
			Message: "Invalid plan type",
		},
		UserNotSoftDeleted: Error{
			Code:    "USER_ERROR_005",
			Message: "User not soft deleted",
		},
		CommunityPlanNotCreated: Error{
			Code:    "COMMUNITY_PLAN_ERROR_002",
			Message: "Community-Plan association not created",
		},
		CommunityPlanNotDeleted: Error{
			Code:    "COMMUNITY_PLAN_ERROR_005",
			Message: "Community-Plan association not deleted",
		},
		CommunityServiceNotCreated: Error{
			Code:    "COMMUNITY_SERVICE_ERROR_002",
			Message: "Community-Service association not created",
		},
		CommunityServiceNotDeleted: Error{
			Code:    "COMMUNITY_SERVICE_ERROR_005",
			Message: "Community-Service association not deleted",
		},
		ServiceLocalNotCreated: Error{
			Code:    "SERVICE_LOCAL_ERROR_002",
			Message: "Service-Local association not created",
		},
		ServiceLocalNotDeleted: Error{
			Code:    "SERVICE_LOCAL_ERROR_005",
			Message: "Service-Local association not deleted",
		},
		ServiceProfessionalNotCreated: Error{
			Code:    "SERVICE_PROFESSIONAL_ERROR_002",
			Message: "Service-Professional association not created",
		},
		ServiceProfessionalNotDeleted: Error{
			Code:    "SERVICE_PROFESSIONAL_ERROR_005",
			Message: "Service-Professional association not deleted",
		},
		SessionNotCreated: Error{
			Code:    "SESSION_ERROR_002",
			Message: "Session not created",
		},
		SessionNotUpdated: Error{
			Code:    "SESSION_ERROR_003",
			Message: "Session not updated",
		},
		SessionNotSoftDeleted: Error{
			Code:    "SESSION_ERROR_005",
			Message: "Session not soft deleted",
		},
		EventoNotCreated: Error{
			Code:    "EVENTO_ERROR_002",
			Message: "Evento not created",
		},
		EventoNotUpdated: Error{
			Code:    "EVENTO_ERROR_003",
			Message: "Evento not updated",
		},
		EventoNotFound: Error{
			Code:    "EVENTO_ERROR_001",
			Message: "Evento not found",
		},
		CategoriaNotCreated: Error{
			Code:    "CATEGORIA_ERROR_002",
			Message: "Categoria not created",
		},
		CategoriaNotUpdated: Error{
			Code:    "CATEGORIA_ERROR_003",
			Message: "Categoria not updated",
		},
		CategoriaNotFound: Error{
			Code:    "CATEGORIA_ERROR_001",
			Message: "Categoria not found",
		},
		InvalidUploadMimeType: Error{
			Code:    "UPLOAD_ERROR_001",
			Message: "Invalid file type. Only images and videos are allowed",
		},
		UploadURLNotCreated: Error{
			Code:    "UPLOAD_ERROR_002",
			Message: "Upload URL could not be generated",
		},
		InvalidIDParam: Error{
			Code:    "REQUEST_ERROR_003",
			Message: "Invalid ID parameter",
		},
	}

	// For 401 Unauthorized errors
	AuthenticationError = struct {
		UnauthorizedUser    Error
		InvalidRefreshToken Error
		InvalidAccessToken  Error
	}{
		UnauthorizedUser: Error{
			Code:    "AUTHENTICATION_ERROR_001",
			Message: "Unauthorized",
		},
		InvalidRefreshToken: Error{
			Code:    "AUTHENTICATION_ERROR_002",
			Message: "Invalid refresh token",
		},
		InvalidAccessToken: Error{
			Code:    "AUTHENTICATION_ERROR_003",
			Message: "Invalid access token",
		},
	}

	// For 409 Conflict errors
	ConflictError = struct {
		EmailAlreadyExists Error
		UserAlreadyExists  Error
		CuponAlreadyExists Error
	}{
		UserAlreadyExists: Error{
			Code:    "USER_ERROR_006",
			Message: "User already exists with this email",
		},
		CuponAlreadyExists: Error{
			Code:    "CUPON_ERROR_001",
			Message: "Cupon already exists in this event",
		},
		EmailAlreadyExists: Error{
			Code:    "USER_ERROR_007",
			Message: "Email already exists",
		},
	}

	// For 500 Internal Server errors
	InternalServerError = struct {
		Default               Error
		PasswordHashingFailed Error
	}{
		Default: Error{
			Code:    "INTERNAL_SERVER_ERROR_001",
			Message: "An unexpected error occurred.",
		},
		PasswordHashingFailed: Error{
			Code:    "INTERNAL_SERVER_ERROR_002",
			Message: "Password hashing failed",
		},
	}
)

// Helper function to check if an error is in a specific error group.
func isInErrorGroup(err Error, group interface{}) bool {
	val := reflect.ValueOf(group)
	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).Interface() == err {
			return true
		}
	}
	return false
}

// General error handler function for endpoints.
func HandleError(err Error, c echo.Context) error {
	var statusCode int
	switch {
	case isInErrorGroup(err, ObjectNotFoundError):
		statusCode = http.StatusNotFound

	case isInErrorGroup(err, UnprocessableEntityError):
		statusCode = http.StatusUnprocessableEntity

	case isInErrorGroup(err, BadRequestError):
		statusCode = http.StatusBadRequest

	case isInErrorGroup(err, ConflictError):
		statusCode = http.StatusConflict

	case isInErrorGroup(err, InternalServerError):
		statusCode = http.StatusInternalServerError

	default:
		statusCode = http.StatusInternalServerError // Default case for other errors
	}

	// Send JSON response with the error code and message
	return c.JSON(statusCode, err)
}
