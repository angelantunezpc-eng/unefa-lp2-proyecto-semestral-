<?php

declare(strict_types=1);

enum TipoRed: string
{
    case LTE_4G = '4G_LTE';
    case NR_5G_NSA = '5G_NSA';
    case NR_5G_SA = '5G_SA';
    case ROAMING = 'ROAMING';
}

readonly class TarifaCalculada
{
    public function __construct(
        public string $msisdn,
        public TipoRed $tipoRed,
        public float $tarifaBasePorMb,
        public int $mbConsumidos,
        public float $totalUsd,
    ) {}
}

readonly class MotorTarifasJit
{
    public function calcularTarifa(string $msisdn, TipoRed $tipoRed, int $mbConsumidos): TarifaCalculada
    {
        if ($mbConsumidos < 0) {
            throw new InvalidArgumentException('Los MB consumidos no pueden ser negativos.');
        }

        $tarifaBasePorMb = match ($tipoRed) {
            TipoRed::LTE_4G => 0.0025,
            TipoRed::NR_5G_NSA => 0.0040,
            TipoRed::NR_5G_SA => 0.0065,
            TipoRed::ROAMING => 0.0250,
        };

        $totalUsd = round($tarifaBasePorMb * $mbConsumidos, 4);

        return new TarifaCalculada(
            msisdn: $msisdn,
            tipoRed: $tipoRed,
            tarifaBasePorMb: $tarifaBasePorMb,
            mbConsumidos: $mbConsumidos,
            totalUsd: $totalUsd,
        );
    }
}

$motor = new MotorTarifasJit();
$pruebas = [
    $motor->calcularTarifa('+584121234567', TipoRed::LTE_4G, 1024),
    $motor->calcularTarifa('+584141234567', TipoRed::NR_5G_NSA, 2048),
    $motor->calcularTarifa('+584241234567', TipoRed::NR_5G_SA, 512),
    $motor->calcularTarifa('+584161234567', TipoRed::ROAMING, 100),
];

foreach ($pruebas as $item) {
    echo sprintf(
        "MSISDN=%s | Red=%s | MB=%d | Total=%.4f USD\n",
        $item->msisdn,
        $item->tipoRed->value,
        $item->mbConsumidos,
        $item->totalUsd,
    );
}