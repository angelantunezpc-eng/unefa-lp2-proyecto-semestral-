package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

const WorkerCount = 4

type Transaccion struct {
	ID            string
	MSISDN        string
	MontoCentavos int64
}

type ResultadoAuditoria struct {
	TransaccionID string
	Aprobada      bool
	Motivo        string
	WorkerID      int
}

type ProcesadorTransaccion interface {
	Auditar(ctx context.Context, tx Transaccion, workerID int) ResultadoAuditoria
}

type AuditorReglas struct {
	UmbralCentavos int64
}

func (a AuditorReglas) Auditar(_ context.Context, tx Transaccion, workerID int) ResultadoAuditoria {
	if tx.MontoCentavos <= 0 {
		return ResultadoAuditoria{
			TransaccionID: tx.ID,
			Aprobada:      false,
			Motivo:        "MONTO_NO_POSITIVO",
			WorkerID:      workerID,
		}
	}
	if tx.MontoCentavos > a.UmbralCentavos {
		return ResultadoAuditoria{
			TransaccionID: tx.ID,
			Aprobada:      false,
			Motivo:        "EXCEDE_UMBRAL_RIESGO",
			WorkerID:      workerID,
		}
	}
	return ResultadoAuditoria{
		TransaccionID: tx.ID,
		Aprobada:      true,
		Motivo:        "OK",
		WorkerID:      workerID,
	}
}

func EjecutarWorkerPool(
	ctx context.Context,
	workers int,
	procesador ProcesadorTransaccion,
	jobs <-chan Transaccion,
	results chan<- ResultadoAuditoria,
	procesadasTotal *atomic.Int64,
	alertasTotal *atomic.Int64,
) {
	if workers <= 0 {
		workers = WorkerCount
	}

	var wg sync.WaitGroup
	wg.Add(workers)

	for w := 1; w <= workers; w++ {
		workerID := w
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case tx, ok := <-jobs:
					if !ok {
						return
					}
					res := procesador.Auditar(ctx, tx, workerID)
					procesadasTotal.Add(1)
					if !res.Aprobada {
						alertasTotal.Add(1)
					}
					select {
					case <-ctx.Done():
						return
					case results <- res:
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lote := []Transaccion{
		{ID: "TX-101", MSISDN: "+584121234567", MontoCentavos: 2500},
		{ID: "TX-102", MSISDN: "+584141234567", MontoCentavos: 250000},
		{ID: "TX-103", MSISDN: "+584241234567", MontoCentavos: 1800},
		{ID: "TX-104", MSISDN: "+584161234567", MontoCentavos: 0},
		{ID: "TX-105", MSISDN: "+584261234567", MontoCentavos: 9500},
		{ID: "TX-106", MSISDN: "+584129876543", MontoCentavos: 150000},
	}

	jobs := make(chan Transaccion, len(lote))
	results := make(chan ResultadoAuditoria, len(lote))

	var procesadas atomic.Int64
	var alertas atomic.Int64

	auditor := AuditorReglas{UmbralCentavos: 100000}
	EjecutarWorkerPool(ctx, WorkerCount, auditor, jobs, results, &procesadas, &alertas)

	for _, tx := range lote {
		jobs <- tx
	}
	close(jobs)

	for res := range results {
		fmt.Printf(
			"[Worker %d] tx=%s aprobada=%v motivo=%s\n",
			res.WorkerID, res.TransaccionID, res.Aprobada, res.Motivo,
		)
	}

	fmt.Printf(
		"Total procesadas=%d | Total alertas=%d\n",
		procesadas.Load(), alertas.Load(),
	)
}
