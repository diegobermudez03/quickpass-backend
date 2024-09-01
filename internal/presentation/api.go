package presentation

import (
	"net/http"

	"github.com/foxinuni/quickpass-backend/internal/presentation/middlewares"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type QuickpassAPIOptions interface {
	GetListenAddress() string
}

type QuickpassAPI struct {
	options QuickpassAPIOptions
}

func NewQuickpassAPI(options QuickpassAPIOptions) *QuickpassAPI {
	return &QuickpassAPI{
		options: options,
	}
}

func (api *QuickpassAPI) Listen() error {
	// Create a new Echo instance
	app := echo.New()

	// Hide the banner and port
	app.HideBanner = true
	app.HidePort = true

	// Use Go-Playground's validator for DTOs
	app.Validator = &CustomValidator{validator: validator.New()}

	// Register the middlewares
	app.Use(middlewares.RequestLogMiddleware)
	app.Use(middlewares.ErrorHandlerMiddleware)

	app.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	log.Info().Msgf("HTTP server is now listening on %s", api.options.GetListenAddress())
	if err := app.Start(api.options.GetListenAddress()); err != nil {
		return err
	}

	return nil
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(dto interface{}) error {
	if err := cv.validator.Struct(dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}