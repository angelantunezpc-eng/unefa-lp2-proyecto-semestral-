package adapters

import (
	"context"
	"sync"
	"sync/atomic"

	"unefa-lp2/lab2_arquitectura_hexagonal/internal/domain"
)

type RepositorioMemoria struct {
	mu      sync.RWMutex
	registros map[string]domain.Recarga
}

func NuevoRepositorioMemoria() *RepositorioMemoria {
	return &RepositorioMemoria{
		registros: make(map[string]domain.Recarga),
	}
}

func (r *RepositorioMemoria) BuscarPorClave(_ context.Context, key string) (domain.Recarga, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rec, ok := r.registros[key]
	return rec, ok, nil
}

func (r *RepositorioMemoria) Guardar(_ context.Context, recarga domain.Recarga) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.registros[recarga.IdempotencyKey] = recarga
	return nil
}

type PasarelaSimulada struct {
	cobrosEfectivos atomic.Int64
}

func NuevaPasarelaSimulada() *PasarelaSimulada {
	return &PasarelaSimulada{}
}

func (p *PasarelaSimulada) Cobrar(_ context.Context, _ domain.Recarga) error {
	p.cobrosEfectivos.Add(1)
	return nil
}

func (p *PasarelaSimulada) TotalCobros() int64 {
	return p.cobrosEfectivos.Load()
}
