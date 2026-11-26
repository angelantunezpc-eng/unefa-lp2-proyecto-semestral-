#include <inttypes.h>
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

static uint64_t g_contador_cdr_global = 42ULL;

struct CDR_Desalineado {
    char tipo_red;
    uint64_t imsi_abonado;
    char bandera_roaming;
    uint32_t celda_id;
    char codigo_estado;
    uint16_t puerto_origen;
};

struct CDR_Optimizado {
    uint64_t imsi_abonado;
    uint32_t celda_id;
    uint16_t puerto_origen;
    char tipo_red;
    char bandera_roaming;
    char codigo_estado;
};

int main(void) {
    uint64_t local_stack_marker = 99ULL;

    uint64_t *heap_marker = (uint64_t *)malloc(sizeof(uint64_t));
    if (heap_marker == NULL) {
        fprintf(stderr, "Error al asignar memoria en el Heap\n");
        return EXIT_FAILURE;
    }
    *heap_marker = 12345ULL;

    printf("=== (a) DIRECCIONES VIRTUALES DE MEMORIA ===\n");
    printf("Global (.data): %p (valor: %" PRIu64 ")\n", (const void *)&g_contador_cdr_global, g_contador_cdr_global);
    printf("Local (Stack) : %p (valor: %" PRIu64 ")\n", (const void *)&local_stack_marker, local_stack_marker);
    printf("Dinamico(Heap): %p (valor: %" PRIu64 ")\n\n", (const void *)heap_marker, *heap_marker);

    free(heap_marker);
    heap_marker = NULL;

    printf("=== (b) TAMANO Y OFFSETS DE ESTRUCTURAS ===\n");
    printf("CDR_Desalineado sizeof = %zu bytes\n", sizeof(struct CDR_Desalineado));
    printf("  offset tipo_red       = %zu\n", offsetof(struct CDR_Desalineado, tipo_red));
    printf("  offset imsi_abonado   = %zu\n", offsetof(struct CDR_Desalineado, imsi_abonado));
    printf("  offset bandera_roaming= %zu\n", offsetof(struct CDR_Desalineado, bandera_roaming));
    printf("  offset celda_id       = %zu\n", offsetof(struct CDR_Desalineado, celda_id));
    printf("  offset codigo_estado  = %zu\n", offsetof(struct CDR_Desalineado, codigo_estado));
    printf("  offset puerto_origen  = %zu\n\n", offsetof(struct CDR_Desalineado, puerto_origen));

    printf("CDR_Optimizado sizeof = %zu bytes\n", sizeof(struct CDR_Optimizado));
    printf("  offset imsi_abonado   = %zu\n", offsetof(struct CDR_Optimizado, imsi_abonado));
    printf("  offset celda_id       = %zu\n", offsetof(struct CDR_Optimizado, celda_id));
    printf("  offset puerto_origen  = %zu\n", offsetof(struct CDR_Optimizado, puerto_origen));
    printf("  offset tipo_red       = %zu\n", offsetof(struct CDR_Optimizado, tipo_red));
    printf("  offset bandera_roaming= %zu\n", offsetof(struct CDR_Optimizado, bandera_roaming));
    printf("  offset codigo_estado  = %zu\n\n", offsetof(struct CDR_Optimizado, codigo_estado));

    const double lote = 10000000.0;
    const size_t l1_line_bytes = 64;
    double bytes_desalineado = lote * (double)sizeof(struct CDR_Desalineado);
    double bytes_optimizado = lote * (double)sizeof(struct CDR_Optimizado);
    double millones_bytes_ahorrados = (bytes_desalineado - bytes_optimizado) / 1000000.0;

    printf("=== (c) AHORRO EN 10,000,000 CDR Y DENSIDAD EN LINEA L1 (64B) ===\n");
    printf("RAM lote desalineado : %.2f millones de bytes\n", bytes_desalineado / 1000000.0);
    printf("RAM lote optimizado  : %.2f millones de bytes\n", bytes_optimizado / 1000000.0);
    printf("RAM ahorrada         : %.2f millones de bytes\n", millones_bytes_ahorrados);
    printf("Registros por linea L1 (64B) - Desalineado: %zu\n", l1_line_bytes / sizeof(struct CDR_Desalineado));
    printf("Registros por linea L1 (64B) - Optimizado : %zu\n", l1_line_bytes / sizeof(struct CDR_Optimizado));

    return EXIT_SUCCESS;
}
