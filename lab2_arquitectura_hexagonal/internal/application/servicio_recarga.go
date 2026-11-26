package application

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"unefa-lp2/lab2_arquitectura_hexagonal/internal/domain"
	"unefa-lp2/lab2_arquitectura_hexagonal/internal/ports"
)

type ServicioRecarga struct {
	repo     ports.RepositorioIdempotencia
	pasarela ports.PasarelaCobroSecundaria
	reloj    func() time.Time
	seq      atomic.Uint64
}

func NuevoServicioRecarga(repo ports.RepositorioIdempotencia, pasarela ports.PasarelaCobroSecundaria) *ServicioRecarga {
	return &ServicioRecarga{
		repo:     repo,
		pasarela: pasarela,
		reloj:    time.Now,
	}
}

func (s *ServicioRecarga) EjecutarRecarga(ctx context.Context, cmd ports.ComandoRecarga) (ports.ResultadoRecarga, error) {
	existente, ok, err := s.repo.BuscarPorClave(ctx, cmd.IdempotencyKey)
	if err != nil {
		return ports.ResultadoRecarga{}, err
	}
	if ok {
		return ports.ResultadoRecarga{
			Recarga:     existente,
			Idempotente: true,
		}, nil
	}

	idNum := s.seq.Add(1)
	id := fmt.Sprintf("TXN-%06d", idNum)

	recarga, err := domain.NuevaRecarga(id, cmd.IdempotencyKey, cmd.Phone, cmd.Amount, cmd.Currency, s.reloj())
	if err != nil {
		return ports.ResultadoRecarga{}, err
	}

	if err := s.pasarela.Cobrar(ctx, recarga); err != nil {
		return ports.ResultadoRecarga{}, err
	}

	if err := s.repo.Guardar(ctx, recarga); err != nil {
		return ports.ResultadoRecarga{}, err
	}

	return ports.ResultadoRecarga{
		Recarga:     recarga,
		Idempotente: false,
	}, nil
}
