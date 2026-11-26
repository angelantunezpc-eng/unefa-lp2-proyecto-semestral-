package domain

import (
	"errors"
	"regexp"
	"time"
)

var (
	ErrInvalidPhone      = errors.New("numero telefonico invalido: debe cumplir formato E.164 (+[1-9][0-9]{6,14})")
	ErrInvalidAmount     = errors.New("monto invalido: debe ser mayor que cero")
	ErrInvalidCurrency   = errors.New("moneda invalida: debe tener 3 letras")
	ErrMissingIdempotency = errors.New("clave de idempotencia requerida")
)

var e164Regex = regexp.MustCompile(`^\+[1-9][0-9]{6,14}$`)

type EstadoTransaccion string

const (
	EstadoAplicada EstadoTransaccion = "APLICADA"
)

type Recarga struct {
	ID             string
	IdempotencyKey string
	Phone          string
	Amount         float64
	Currency       string
	Estado         EstadoTransaccion
	CreadaEn       time.Time
}

func NuevaRecarga(id, idempotencyKey, phone string, amount float64, currency string, ahora time.Time) (Recarga, error) {
	if idempotencyKey == "" {
		return Recarga{}, ErrMissingIdempotency
	}
	if !e164Regex.MatchString(phone) {
		return Recarga{}, ErrInvalidPhone
	}
	if amount <= 0 {
		return Recarga{}, ErrInvalidAmount
	}
	if len(currency) != 3 {
		return Recarga{}, ErrInvalidCurrency
	}

	return Recarga{
		ID:             id,
		IdempotencyKey: idempotencyKey,
		Phone:          phone,
		Amount:         amount,
		Currency:       currency,
		Estado:         EstadoAplicada,
		CreadaEn:       ahora.UTC(),
	}, nil
}
