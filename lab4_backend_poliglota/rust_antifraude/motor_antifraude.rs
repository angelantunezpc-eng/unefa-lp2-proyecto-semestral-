use std::sync::{Arc, Mutex};
use std::thread;

#[derive(Debug)]
struct AlertaFraude<'a> {
    id_tx: &'a str,
    msisdn: &'a str,
    motivo: &'a str,
    monto_centavos: u64,
}

#[derive(Debug, Default)]
struct EstadisticasRiesgo {
    transacciones_auditadas: u64,
    alertas_generadas: u64,
    monto_riesgo_centavos: u64,
}

fn evaluar_lote<'a>(
    lote: &[(&'a str, &'a str, u64)],
    umbral_centavos: u64,
    stats_arc: &Arc<Mutex<EstadisticasRiesgo>>,
) -> Vec<AlertaFraude<'a>> {
    let mut alertas_locales = Vec::new();
    let mut auditadas_locales = 0_u64;
    let mut alertas_locales_cnt = 0_u64;
    let mut monto_riesgo_local = 0_u64;

    for &(id_tx, msisdn, monto_centavos) in lote {
        auditadas_locales += 1;

        let motivo = if monto_centavos >= umbral_centavos {
            Some("MONTO_ATIPICO_EXCEDE_UMBRAL")
        } else if !msisdn.starts_with('+') || msisdn.len() < 8 {
            Some("MSISDN_NO_E164_SOSPECHOSO")
        } else {
            None
        };

        if let Some(razon) = motivo {
            alertas_locales_cnt += 1;
            monto_riesgo_local = monto_riesgo_local.saturating_add(monto_centavos);
            alertas_locales.push(AlertaFraude {
                id_tx,
                msisdn,
                motivo: razon,
                monto_centavos,
            });
        }
    }

    if let Ok(mut guard) = stats_arc.lock() {
        guard.transacciones_auditadas += auditadas_locales;
        guard.alertas_generadas += alertas_locales_cnt;
        guard.monto_riesgo_centavos = guard
            .monto_riesgo_centavos
            .saturating_add(monto_riesgo_local);
    }

    alertas_locales
}

fn main() {
    let transacciones: Vec<(&str, &str, u64)> = vec![
        ("TX-9001", "+584121110001", 1_500),
        ("TX-9002", "+584149990002", 250_000),
        ("TX-9003", "04165550003", 4_200),
        ("TX-9004", "+584241110004", 3_100),
        ("TX-9005", "+584128880005", 500_000),
        ("TX-9006", "+584162220006", 2_000),
        ("TX-9007", "+584267770007", 310_000),
        ("TX-9008", "+584123330008", 1_800),
    ];

    let umbral_centavos: u64 = 100_000;
    let num_hilos: usize = 4;
    let stats = Arc::new(Mutex::new(EstadisticasRiesgo::default()));

    let chunk_size = transacciones.len().div_ceil(num_hilos);
    let mut handles = Vec::with_capacity(num_hilos);
    let mut todas_alertas: Vec<AlertaFraude<'_>> = Vec::new();

    thread::scope(|scope| {
        let mut scoped_handles = Vec::new();
        for lote in transacciones.chunks(chunk_size) {
            let stats_ref = Arc::clone(&stats);
            let h = scope.spawn(move || evaluar_lote(lote, umbral_centavos, &stats_ref));
            scoped_handles.push(h);
        }

        for h in scoped_handles {
            if let Ok(mut parciales) = h.join() {
                todas_alertas.append(&mut parciales);
            }
        }
    });

    let _ = handles.len();

    for alerta in &todas_alertas {
        println!(
            "ALERTA -> id_tx={} msisdn={} motivo={} monto_centavos={}",
            alerta.id_tx, alerta.msisdn, alerta.motivo, alerta.monto_centavos
        );
    }

    let final_stats = stats.lock().expect("mutex envenenado");
    println!(
        "RESUMEN -> auditadas={} alertas={} monto_riesgo_centavos={}",
        final_stats.transacciones_auditadas,
        final_stats.alertas_generadas,
        final_stats.monto_riesgo_centavos
    );
}