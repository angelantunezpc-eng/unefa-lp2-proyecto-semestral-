# Árbol de Fallos (FTA) y Diagrama de Ishikawa — Corrupción del Registro `RIP` por Desbordamiento de Pila

## 1. Diagrama de Ishikawa (Espina de Pescado)

```mermaid
flowchart LR
    subgraph Causas["Causas Potenciales"]
        direction TB
        M1["Código y Lenguaje<br/>(gets/strcpy sin límites, aritmética de punteros)"]
        M2["Compilación y Entorno<br/>(sin stack canary, NX/DEP deshabilitado, ASLR débil)"]
        M3["Datos de Entrada<br/>(payload mayor al buffer en Stack, falta de validación)"]
        M4["Proceso y Revisión<br/>(ausencia de análisis estático y pruebas de fuzzing)"]
    end
    Efecto(["Corrupción del registro RIP<br/>por desbordamiento de pila"])
    M1 --> Efecto
    M2 --> Efecto
    M3 --> Efecto
    M4 --> Efecto
