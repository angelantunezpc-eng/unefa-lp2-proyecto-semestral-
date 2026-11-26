package main

import (
	"fmt"
	"time"
)

type Sesion5G struct {
	AbonadoID  int64
	CeldaID    uint32
	LatenciaMs float64
	Activa     bool
}

func CrearEnStackPorValor(abonadoID int64, celdaID uint32, latencia float64) Sesion5G {
	s := Sesion5G{
		AbonadoID:  abonadoID,
		CeldaID:    celdaID,
		LatenciaMs: latencia,
		Activa:     true,
	}
	return s
}

func CrearConEscapeAlHeap(abonadoID int64, celdaID uint32, latencia float64) *Sesion5G {
	s := Sesion5G{
		AbonadoID:  abonadoID,
		CeldaID:    celdaID,
		LatenciaMs: latencia,
		Activa:     true,
	}
	return &s
}

func buscarLineal(slice []Sesion5G, objetivo int64) bool {
	for i := range slice {
		if slice[i].AbonadoID == objetivo {
			return true
		}
	}
	return false
}

func main() {
	sStack := CrearEnStackPorValor(1001, 501, 11.4)
	sHeap := CrearConEscapeAlHeap(1002, 502, 14.8)
	fmt.Printf("Sesion Stack: %+v | Sesion Heap: %+v\n", sStack, *sHeap)

	const N = 100_000
	lote := make([]Sesion5G, N)
	indiceHash := make(map[int64]Sesion5G, N)

	for i := 0; i < N; i++ {
		id := int64(i + 1)
		reg := Sesion5G{
			AbonadoID:  id,
			CeldaID:    uint32(i % 256),
			LatenciaMs: float64(5 + (i % 45)),
			Activa:     true,
		}
		lote[i] = reg
		indiceHash[id] = reg
	}

	const consultas = 1_000
	objetivos := make([]int64, consultas)
	for i := 0; i < consultas; i++ {
		objetivos[i] = int64(((i * 97) % N) + 1)
	}

	inicioLineal := time.Now()
	encontradosLineal := 0
	for _, obj := range objetivos {
		if buscarLineal(lote, obj) {
			encontradosLineal++
		}
	}
	duracionLinealUs := float64(time.Since(inicioLineal).Microseconds())

	inicioHash := time.Now()
	encontradosHash := 0
	for _, obj := range objetivos {
		if _, ok := indiceHash[obj]; ok {
			encontradosHash++
		}
	}
	duracionHashUs := float64(time.Since(inicioHash).Microseconds())

	fmt.Printf("Benchmark N=%d, consultas=%d:\n", N, consultas)
	fmt.Printf("  Busqueda Lineal O(n): %.2f us (encontrados=%d)\n", duracionLinealUs, encontradosLineal)
	fmt.Printf("  Busqueda Mapa   O(1): %.2f us (encontrados=%d)\n", duracionHashUs, encontradosHash)
}
