package main

import (
	"github.com/Rhyanz46/go-map-validator/map_validator"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
)

type LoginData struct {
	Email string `map_validator:"email"`
	Pw    string `map_validator:"password"`
}

func handleLogin(c echo.Context) error {
	data := &LoginData{}
	jsonHttp, err := map_validator.NewValidateBuilder().SetRules(map[string]map_validator.Rules{
		"email":    {Email: true, Max: map_validator.SetTotal(100)},
		"password": {Type: reflect.String, Min: map_validator.SetTotal(6), Max: map_validator.SetTotal(30)},
	}).LoadJsonHttp(c.Request())
	if err != nil {
		switch err {
		case map_validator.ErrNoData:
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "you need to input json body",
			})
		case map_validator.ErrInvalidFormat, map_validator.ErrUnsupportType:
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "json format is invalid",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "internal server error",
		})
	}
	jsonData, err := jsonHttp.RunValidate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	err = jsonData.Bind(data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
	return c.NoContent(http.StatusOK)
}

func main() {
	e := echo.New()
	e.POST("/login", handleLogin)
	e.Start(":3000")
}
