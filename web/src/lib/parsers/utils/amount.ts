/**
 * Parses a monetary amount string to cents (integer)
 * Handles various formats: 1.234,56 / 1,234.56 / -1234.56 / R$ 1.234,56
 */
export function parseAmountToCents(value: string): number | undefined {
  if (!value) {
    return undefined
  }

  const sanitized = value.replace(/[^\d,.-]/g, '')
  if (!sanitized) {
    return undefined
  }

  const lastComma = sanitized.lastIndexOf(',')
  const lastDot = sanitized.lastIndexOf('.')
  const separatorIndex = Math.max(lastComma, lastDot)

  let normalized = sanitized

  if (separatorIndex >= 0) {
    const integerPart = sanitized
      .slice(0, separatorIndex)
      .replace(/[^0-9-]/g, '')
      .replace(/(?!^)-/g, '')
    const fractionPart = sanitized.slice(separatorIndex + 1).replace(/\D/g, '')
    normalized = `${integerPart}.${fractionPart}`
  } else {
    normalized = sanitized.replace(/[^0-9-]/g, '')
  }

  const amount = Number(normalized)

  if (Number.isNaN(amount)) {
    return undefined
  }

  return Math.round(amount * 100)
}
