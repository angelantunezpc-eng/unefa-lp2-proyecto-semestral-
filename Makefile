.PHONY: all lab1 lab2 lab3 lab4 clean

all: lab1 lab2 lab3 lab4

lab1:
	@echo "=== [LAB 1] Perfilador de Memoria C y Escape Analysis Go ==="
	gcc -Wall -Wextra -Werror -O0 lab1_memoria_y_complejidad/cdr_memory_profiler.c -o cdr_profiler
	./cdr_profiler
	go build -gcflags="-m" -o escape_analyzer lab1_memoria_y_complejidad/escape_analyzer.go
	go run lab1_memoria_y_complejidad/escape_analyzer.go

lab2:
	@echo "=== [LAB 2] Arquitectura Hexagonal Idempotente (Go) ==="
	cd lab2_arquitectura_hexagonal && go vet ./... && go run ./cmd/api/main.go

lab3:
	@echo "=== [LAB 3] Auditoria Frontend Reactivo (React 19 vs Svelte 5) ==="
	@test -f lab3_frontend_reactivo/nextjs_react19_vdom/TelemetryPanelClient.tsx
	@test -f lab3_frontend_reactivo/svelte5_runes_aot/+page.svelte
	@test -f lab3_frontend_reactivo/comparativa_web_vitals.md
	@echo "Archivos del Lab 3 verificados correctamente."

lab4:
	@echo "=== [LAB 4] Motor Backend Poliglota ==="
	rustc -O lab4_backend_poliglota/rust_antifraude/motor_antifraude.rs -o motor_antifraude
	./motor_antifraude
	go run -race lab4_backend_poliglota/go_worker_pool/worker_pool.go
	node lab4_backend_poliglota/runtimes_alto_nivel/node_crypto_worker.mjs
	python3 lab4_backend_poliglota/runtimes_alto_nivel/python_async_validator.py
	php lab4_backend_poliglota/runtimes_alto_nivel/php_jit_tariff_engine.php

clean:
	rm -f cdr_profiler cdr_profiler.exe escape_analyzer escape_analyzer.exe motor_antifraude motor_antifraude.exe