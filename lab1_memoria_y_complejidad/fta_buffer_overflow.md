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
flowchart TD
    TOP["Evento Superior:<br/>Corrupción de RIP por Stack Buffer Overflow"]
    AND1{{"AND"}}
    OR1{{"OR"}}
    OR2{{"OR"}}
    E1["Escritura fuera de límites<br/>en buffer del Stack"]
    E2["Sobrescritura de dirección<br/>de retorno guardada (saved RIP)"]
    E3["Copia sin comprobación<br/>de longitud de entrada"]
    E4["Payload externo excede<br/>capacidad asignada"]
    E5["Canary ausente o eludido<br/>en marco de pila"]
    E6["Flujo de retorno alcanza<br/>instrucción ret con RIP alterado"]

    TOP --> AND1
    AND1 --> OR1
    AND1 --> OR2
    OR1 --> E1
    OR1 --> E2
    E1 --> E3
    E1 --> E4
    OR2 --> E5
    OR2 --> E6
