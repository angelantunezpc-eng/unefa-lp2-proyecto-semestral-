# UNEFA — LP2: Proyecto Semestral (TelCo / FinTech Convergente)

## Arquitectura General (Modelo C4 — Nivel 1 y 2)

```mermaid
C4Context
    title Diagrama de Contexto (C4 Nivel 1) — Pasarela y Telemetría 5G
    Person(abonado, "Abonado / Operador", "Solicita recargas y monitorea celdas 5G")
    System(sistema, "Plataforma Convergente LP2", "Procesa recargas idempotentes, telemetría 5G y motor antifraude políglota")
    System_Ext(pasarela, "Pasarela Bancaria / Core TelCo", "Autoriza débitos y registra CDRs")
    Rel(abonado, sistema, "HTTPS / REST / UI Reactiva", "JSON / Idempotency-Key")
    Rel(sistema, pasarela, "Autorización y Liquidación", "Puerto Secundario")