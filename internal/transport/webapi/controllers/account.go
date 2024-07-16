package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vitaliy-ukiru/bank-service/internal/application"
	"github.com/vitaliy-ukiru/bank-service/internal/domain/account"
	"github.com/vitaliy-ukiru/bank-service/internal/transport/webapi/request"
	"github.com/vitaliy-ukiru/bank-service/internal/transport/webapi/response"
)

type Usecase interface {
	CreateAccount(ctx context.Context) (int64, error)
	DepositBalance(ctx context.Context, cmd application.DepositBalanceCommand) error
	WithdrawBalance(ctx context.Context, cmd application.WithdrawBalanceCommand) error
	GetBalance(ctx context.Context, cmd application.GetBalanceCommand) (float64, error)
}

type AccountController struct {
	uc Usecase
}

func NewAccountController(uc Usecase) *AccountController {
	return &AccountController{uc: uc}
}

func (a AccountController) Bind(e *echo.Echo) {
	g := e.Group("/accounts")
	g.POST("", a.CreateAccount)
	g.POST("/:id/deposit", a.Deposit)
	g.POST("/:id/withdraw", a.Withdraw)
	g.GET("/:id/balance", a.GetAccountBalance)
}

const unknownError = "unknown error occurred"

func (a AccountController) CreateAccount(c echo.Context) error {
	ctx := c.Request().Context()
	accountId, err := a.uc.CreateAccount(ctx)
	if err != nil {
		return processError(c, err)
	}

	return c.JSON(http.StatusCreated, response.Ok(response.M{
		"id": accountId,
	}))
}

func (a AccountController) Deposit(c echo.Context) error {
	var req request.DepositRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusUnprocessableEntity,
			response.Error(fmt.Errorf("invalid request format: %w", err)),
		)
	}

	ctx := getContext(c)
	err := a.uc.DepositBalance(ctx, application.DepositBalanceCommand{
		AccountId: req.AccountId,
		Amount:    req.Amount,
	})
	if err != nil {
		return processError(c, err)
	}

	return c.JSON(http.StatusOK, response.OkStatus)
}

func (a AccountController) Withdraw(c echo.Context) error {
	var req request.WithdrawRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusUnprocessableEntity,
			response.Error(fmt.Errorf("invalid request format: %w", err)),
		)
	}

	ctx := getContext(c)
	err := a.uc.WithdrawBalance(ctx, application.WithdrawBalanceCommand{
		AccountId: req.AccountId,
		Amount:    req.Amount,
	})
	if err != nil {
		return processError(c, err)
	}

	return c.JSON(http.StatusOK, response.OkStatus)
}

func (a AccountController) GetAccountBalance(c echo.Context) error {
	var req request.GetBalanceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusUnprocessableEntity,
			response.Error(fmt.Errorf("invalid request format: %w", req)),
		)
	}

	ctx := getContext(c)
	balance, err := a.uc.GetBalance(ctx, application.GetBalanceCommand{
		AccountId: req.AccountId,
	})
	if err != nil {
		return processError(c, err)
	}

	return c.JSON(http.StatusOK, response.Ok(response.M{
		"balance": balance,
	}))
}

func getContext(c echo.Context) context.Context {
	return c.Request().Context()
}

func processError(c echo.Context, err error) error {
	resp := response.Error(errors.Unwrap(err))
	if errors.Is(err, application.ErrAccountNotFound) {
		return c.JSON(http.StatusNotFound, resp)
	}

	if errors.Is(err, account.ErrNegativeAmount) || errors.Is(err, account.ErrZeroAmount) {
		return c.JSON(http.StatusBadRequest, resp)
	}

	if errors.Is(err, account.ErrNotEnoughBalance) {
		return c.JSON(http.StatusConflict, resp)
	}

	return c.JSON(500, response.Fail(unknownError))
}
