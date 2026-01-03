import { randomUUID } from 'node:crypto'

import { ExpenseParser, FileInfo, ParseResult } from '../types'
import { parseAmountToCents } from '../utils/amount'
import { detectDelimiter, splitCsvRow } from '../utils/csv'
import { parseDateString } from '../utils/date'

const FLASH_HEADER_TOKENS = ['data', 'hora', 'movimentacao', 'valor', 'meio de pagamento']

function normalize(value: string): string {
  return value
    .trim()
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
}

function looksLikeFlashHeader(line: string | undefined): boolean {
  if (!line) {
    return false
  }

  const normalized = normalize(line)
  return FLASH_HEADER_TOKENS.every((token) => normalized.includes(token))
}

/**
 * Parser para arquivos CSV exportados do aplicativo Flash
 * Estrutura esperada: Data, Hora, Movimentação, Valor, Meio de Pagamento, Saldo
 */
export class FlashCsvParser extends ExpenseParser {
  readonly name = 'Flash CSV Parser'
  readonly supportedFormats = ['text/csv', 'csv']

  canParse(fileInfo: FileInfo): boolean {
    const isCsvMime = fileInfo.mimeType === 'text/csv'
    const isCsvExtension = fileInfo.filename?.toLowerCase().endsWith('.csv')

    if (typeof fileInfo.content !== 'string' || !(isCsvMime || isCsvExtension)) {
      return false
    }

    const firstLine = fileInfo.content
      .split(/\r?\n/)
      .map((line) => line.trim())
      .find((line) => line.length > 0)

    return looksLikeFlashHeader(firstLine)
  }

  async parse(fileInfo: FileInfo): Promise<ParseResult> {
    if (typeof fileInfo.content !== 'string') {
      return { expenses: [], errors: ['Content must be a string for CSV parsing'] }
    }

    const content = fileInfo.content
    const rows = content
      .split(/\r?\n/)
      .map((line) => line.trim())
      .filter((line) => line.length > 0)

    if (rows.length <= 1) {
      return { expenses: [] }
    }

    const headerRow = rows[0]
    if (!looksLikeFlashHeader(headerRow)) {
      return { expenses: [], errors: ['Formato de CSV não reconhecido como extrato Flash'] }
    }

    const delimiter = detectDelimiter(headerRow)
    const headers = splitCsvRow(headerRow, delimiter).map(normalize)
    const dateIndex = headers.findIndex((header) => header === 'data')
    const timeIndex = headers.findIndex((header) => header === 'hora')
    const descriptionIndex = headers.findIndex((header) => header === 'movimentacao')
    const amountIndex = headers.findIndex((header) => header === 'valor')

    const expenses = rows.slice(1).reduce<ParseResult['expenses']>((acc, row) => {
      const cells = splitCsvRow(row, delimiter)
      if (!cells.length || cells.every((cell) => cell === '')) {
        return acc
      }

      const description = descriptionIndex >= 0 ? cells[descriptionIndex]?.trim() : undefined
      const dateValue = dateIndex >= 0 ? cells[dateIndex] : undefined
      const timeValue = timeIndex >= 0 ? cells[timeIndex] : undefined
      const rawAmount = amountIndex >= 0 ? cells[amountIndex] : undefined

      const amountInCents = parseAmountToCents(rawAmount ?? '')
      if (!description || typeof amountInCents !== 'number' || amountInCents >= 0) {
        return acc
      }

      let date = parseDateString(dateValue ?? '')
      if (date && timeValue) {
        const [hours, minutes] = timeValue.split(':', 2).map((value) => Number(value.replace(/\D/g, '')))
        if (!Number.isNaN(hours)) {
          const parsedDate = new Date(date)
          parsedDate.setUTCHours(hours, Number.isNaN(minutes) ? 0 : minutes)
          date = parsedDate.toISOString()
        }
      }

      acc.push({
        id: randomUUID(),
        description,
        amountInCents: String(Math.abs(amountInCents)),
        date,
        raw: row,
        splitType: 'proportional',
      })

      return acc
    }, [])

    return { expenses }
  }
}
