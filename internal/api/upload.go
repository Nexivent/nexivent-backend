package api

import (
	"net/http"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/labstack/echo/v4"
)

// @Summary 			Generate presigned S3 upload URL
// @Description 		Returns a presigned PUT URL to upload an image or short video directly to S3.
// @Tags 				Media
// @Accept 				json
// @Produce 			json
// @Param               request body schemas.PresignUploadRequest true "Upload request"
// @Success 			200 {object} schemas.PresignUploadResponse "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/media/upload-url [post]
func (a *Api) GenerateUploadURL(c echo.Context) error {
	if a.BllController.Media == nil {
		return errors.HandleError(errors.BadRequestError.UploadURLNotCreated, c)
	}

	var req schemas.PresignUploadRequest
	if err := c.Bind(&req); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	resp, err := a.BllController.Media.GenerateUploadURL(c.Request().Context(), req)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, resp)
}
