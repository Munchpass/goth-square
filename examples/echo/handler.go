package square

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

func OAuthStart(c echo.Context) error {
	ctx := context.WithValue(c.Request().Context(), gothic.ProviderParamKey, "square")
	newReq := c.Request().WithContext(ctx)
	user, err := gothic.CompleteUserAuth(c.Response(), newReq)
	if err != nil {
		fmt.Println("err: ", err)
		gothic.BeginAuthHandler(c.Response(), newReq)
		return nil
	}

	fmt.Println("user (oauth start): ", user)
	return nil
}

func OAuthCallback(c echo.Context) error {
	ctx := context.WithValue(c.Request().Context(), gothic.ProviderParamKey, "square")
	newReq := c.Request().WithContext(ctx)
	user, err := gothic.CompleteUserAuth(c.Response(), newReq)
	if err != nil {
		return err
	}

	fmt.Println("user (from callback): ", user)
	return nil
}
