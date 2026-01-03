import { randomUUID } from 'node:crypto'

import { ExpenseParser, FileInfo, ParseResult } from '../types'
import { parseAmountToCents } from '../utils/amount'
import { detectDelimiter, pickAmount, pickDate, pickDescription, splitCsvRow } from '../utils/csv'
import { parseDateString } from '../utils/date'

/**
 * Parser for Inter Bank CSV exports
 * Handles the specific CSV format used by Banco Inter
 */
export class InterCsvParser extends ExpenseParser {
  readonly name = 'Inter CSV Parser'
  readonly supportedFormats = ['text/csv', 'csv']

  canParse(fileInfo: FileInfo): boolean {
    const isCsvMime = fileInfo.mimeType === 'text/csv'
    const isCsvExtension = fileInfo.filename?.toLowerCase().endsWith('.csv')
    const isStringContent = typeof fileInfo.content === 'string'

    return Boolean(isStringContent && (isCsvMime || isCsvExtension))
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

    if (!rows.length) {
      return { expenses: [] }
    }

    // Find header row (contains "valor" and has semicolons/commas)
    const headerRowIndex = rows.findIndex((row) => /\bvalor\b/i.test(row) && row.includes(';'))
    const headerRow = rows[headerRowIndex >= 0 ? headerRowIndex : 0]

    const delimiter = detectDelimiter(headerRow)
    const headers = splitCsvRow(headerRow, delimiter).map((header) => header.toLowerCase())
    const descriptionIndex = headers.findIndex((header) => /descr/i.test(header))
    const historicoIndex = headers.findIndex((header) => /hist[oó]rico/i.test(header))
    const amountIndex = headers.findIndex((header) => /(valor|amount)/i.test(header))
    const dateIndex = headers.findIndex((header) => /(data|date)/i.test(header))

    const dataRows = rows.slice(headerRowIndex >= 0 ? headerRowIndex + 1 : 1)

    const expenses = dataRows.reduce<ParseResult['expenses']>((acc, row) => {
      const cells = splitCsvRow(row, delimiter)
      if (cells.every((cell) => cell === '')) {
        return acc
      }

      const description = pickDescription(cells, descriptionIndex, historicoIndex)
      const amountInCents = parseAmountToCents(pickAmount(cells, amountIndex))
      const date = parseDateString(pickDate(cells, dateIndex))

      // Only include expenses (negative values), skip income (positive values)
      if (description && typeof amountInCents === 'number' && amountInCents < 0) {
        acc.push({
          id: randomUUID(),
          description,
          amountInCents: String(Math.abs(amountInCents)),
          date,
          raw: row,
          splitType: 'proportional',
        })
      }

      return acc
    }, [])

    return { expenses }
  }
}
