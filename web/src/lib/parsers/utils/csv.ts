/**
 * Detects the delimiter used in a CSV line
 */
export function detectDelimiter(line: string): string {
  const delimiters = [',', ';', '\t', '|']
  let bestDelimiter = ','
  let bestCount = -1

  for (const delimiter of delimiters) {
    const count = line.split(delimiter).length - 1
    if (count > bestCount) {
      bestCount = count
      bestDelimiter = delimiter
    }
  }

  return bestDelimiter
}

/**
 * Splits a CSV row respecting quoted values
 */
export function splitCsvRow(row: string, delimiter: string): string[] {
  const values: string[] = []
  let current = ''
  let insideQuotes = false

  for (let i = 0; i < row.length; i++) {
    const char = row[i]
    const isQuote = char === '"'

    if (isQuote) {
      if (insideQuotes && row[i + 1] === '"') {
        current += '"'
        i++
      } else {
        insideQuotes = !insideQuotes
      }
      continue
    }

    if (char === delimiter && !insideQuotes) {
      values.push(current.trim())
      current = ''
    } else {
      current += char
    }
  }

  values.push(current.trim())
  return values.map((value) => value.replace(/^['"]|['"]$/g, '').trim())
}

/**
 * Picks the description from CSV cells
 */
export function pickDescription(cells: string[], descriptionIndex: number, historicoIndex?: number): string {
  const parts: string[] = []

  // Add histórico if available (e.g., "Pix enviado", "Pagamento efetuado")
  if (historicoIndex !== undefined && historicoIndex >= 0 && historicoIndex < cells.length) {
    const historico = cells[historicoIndex].trim()
    if (historico) {
      parts.push(historico)
    }
  }

  // Add description if available (e.g., "Pjbank", "Fatura Cartao BTG")
  if (descriptionIndex >= 0 && descriptionIndex < cells.length) {
    const desc = cells[descriptionIndex].trim()
    if (desc) {
      parts.push(desc)
    }
  }

  // If we have both, combine them
  if (parts.length > 0) {
    return parts.join(' - ')
  }

  // Fallback: try to find a non-numeric cell
  return cells.find((cell, index) => index !== 0 && isNaN(Number(cell.replace(/\D/g, '')))) ?? cells.join(' - ')
}

/**
 * Picks the amount from CSV cells
 */
export function pickAmount(cells: string[], amountIndex: number): string {
  if (amountIndex >= 0 && amountIndex < cells.length) {
    return cells[amountIndex]
  }

  const candidate = cells.find((cell) => /[-+]?\d/.test(cell))
  return candidate ?? ''
}

/**
 * Picks the date from CSV cells
 */
export function pickDate(cells: string[], dateIndex: number): string {
  if (dateIndex >= 0 && dateIndex < cells.length) {
    return cells[dateIndex]
  }

  const candidate = cells.find((cell) => /\d{1,2}[/-]\d{1,2}[/-]\d{2,4}/.test(cell))
  return candidate ?? ''
}
