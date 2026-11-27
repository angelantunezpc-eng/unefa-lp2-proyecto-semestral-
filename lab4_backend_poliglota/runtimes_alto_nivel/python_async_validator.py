from __future__ import annotations

import asyncio
import re
from dataclasses import dataclass
from typing import Final

E164_PATTERN: Final[re.Pattern[str]] = re.compile(r"^\+[1-9][0-9]{6,14}$")


@dataclass(slots=True)
class TransactionPayload:
    tx_id: str
    msisdn: str
    amount_cents: int
    currency: str


@dataclass(slots=True)
class ValidationResult:
    tx_id: str
    is_valid: bool
    reasons: list[str]


async def validate_transaction(tx: TransactionPayload) -> ValidationResult:
    await asyncio.sleep(0.01)
    reasons: list[str] = []

    if not tx.tx_id.strip():
        reasons.append("MISSING_TX_ID")
    if not E164_PATTERN.match(tx.msisdn):
        reasons.append("INVALID_E164_MSISDN")
    if tx.amount_cents <= 0:
        reasons.append("NON_POSITIVE_AMOUNT")
    if len(tx.currency) != 3 or not tx.currency.isupper():
        reasons.append("INVALID_CURRENCY_CODE")

    return ValidationResult(
        tx_id=tx.tx_id,
        is_valid=len(reasons) == 0,
        reasons=reasons,
    )


async def validate_batch(batch: list[TransactionPayload]) -> list[ValidationResult]:
    tasks = [validate_transaction(item) for item in batch]
    return list(await asyncio.gather(*tasks))


async def main() -> None:
    sample_batch: list[TransactionPayload] = [
        TransactionPayload(tx_id="PY-001", msisdn="+584121234567", amount_cents=2500, currency="USD"),
        TransactionPayload(tx_id="PY-002", msisdn="04121234567", amount_cents=1200, currency="USD"),
        TransactionPayload(tx_id="PY-003", msisdn="+584149876543", amount_cents=-50, currency="VES"),
        TransactionPayload(tx_id="PY-004", msisdn="+584245551122", amount_cents=8000, currency="EUR"),
    ]

    results = await validate_batch(sample_batch)
    for res in results:
        status = "VALIDA" if res.is_valid else f"RECHAZADA ({', '.join(res.reasons)})"
        print(f"{res.tx_id} -> {status}")


if __name__ == "__main__":
    asyncio.run(main())