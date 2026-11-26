package ports

import (
	"context"

	"unefa-lp2/lab2_arquitectura_hexagonal/internal/domain"
)

type ComandoRecarga struct {
	IdempotencyKey string
	Phone          string
	Amount         float64
	Currency       string
}

type ResultadoRecarga struct {
	Recarga    domain.Recarga
	Idempotente bool
}

type PuertoRecargaEntrada interface {
	EjecutarRecarga(ctx context.Context, cmd ComandoRecarga) (ResultadoRecarga, error)
}

type RepositorioIdempotencia interface {
	BuscarPorClave(ctx context.Context, key string) (domain.Recarga, bool, error)
	Guardar(ctx context.Context, recarga domain.Recarga) error
}

type PasarelaCobroSecundaria interface {
	Cobrar(ctx context.Context, recarga domain.Recarga) error
}
