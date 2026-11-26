package main

import (
	"context"
	"fmt"

	"unefa-lp2/lab2_arquitectura_hexagonal/internal/adapters"
	"unefa-lp2/lab2_arquitectura_hexagonal/internal/application"
	"unefa-lp2/lab2_arquitectura_hexagonal/internal/ports"
)

func main() {
	ctx := context.Background()
	repo := adapters.NuevoRepositorioMemoria()
	pasarela := adapters.NuevaPasarelaSimulada()
	servicio := application.NuevoServicioRecarga(repo, pasarela)

	cmd := ports.ComandoRecarga{
		IdempotencyKey: "key-unefa-2026-001",
		Phone:          "+584121234567",
		Amount:         15.50,
		Currency:       "USD",
	}

	for i := 1; i <= 3; i++ {
		res, err := servicio.EjecutarRecarga(ctx, cmd)
		if err != nil {
			fmt.Printf("Intento %d error: %v\n", i, err)
			continue
		}
		fmt.Printf("Intento %d -> ID: %s | Idempotente: %v | Estado: %s\n",
			i, res.Recarga.ID, res.Idempotente, res.Recarga.Estado)
	}

	_, errInvalido := servicio.EjecutarRecarga(ctx, ports.ComandoRecarga{
		IdempotencyKey: "key-unefa-2026-002",
		Phone:          "04121234567",
		Amount:         10.0,
		Currency:       "USD",
	})
	fmt.Printf("Prueba telefono invalido -> %v\n", errInvalido)
	fmt.Printf("Cobros efectivos en pasarela: %d\n", pasarela.TotalCobros())
}
