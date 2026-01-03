import { randomUUID } from 'node:crypto'

import { ExpenseParser, FileInfo, ParseResult } from '../types'
import { parseAmountToCents } from '../utils/amount'

/**
 * Parser for Inter Bank credit card PDF invoices
 * Recognizes the specific format used by Banco Inter
 */
export class InterCreditCardPdfParser extends ExpenseParser {
  readonly name = 'Inter Credit Card PDF Parser'
  readonly supportedFormats = ['application/pdf', 'pdf']

  canParse(fileInfo: FileInfo): boolean {
    if (typeof fileInfo.content !== 'string') {
      return false
    }

    // Check if it's a PDF and contains Inter-specific markers
    const isPdf = fileInfo.mimeType === 'application/pdf' || fileInfo.filename?.toLowerCase().endsWith('.pdf')
    const hasInterMarker = fileInfo.content.includes('Despesas da fatura')

    return Boolean(isPdf && hasInterMarker)
  }

  async parse(fileInfo: FileInfo): Promise<ParseResult> {
    if (typeof fileInfo.content !== 'string') {
      return { expenses: [], errors: ['Content must be a string for PDF parsing'] }
    }

    const text = fileInfo.content
    const lines = text.split(/\r?\n/).map((line) => line.trim())
    const expenses: ParseResult['expenses'] = []

    // Find the start of transactions section
    const startIndex = lines.findIndex((line) => line.includes('Despesas da fatura'))
    if (startIndex === -1) {
      return { expenses: [] }
    }

    // Portuguese month mapping
    const monthMap: Record<string, number> = {
      jan: 0,
      fev: 1,
      mar: 2,
      abr: 3,
      mai: 4,
      jun: 5,
      jul: 6,
      ago: 7,
      set: 8,
      out: 9,
      nov: 10,
      dez: 11,
    }

    // Process lines after "Despesas da fatura"
    for (let i = startIndex + 1; i < lines.length; i++) {
      const line = lines[i]

      // Match Inter date format at the start of line: "07 de nov. 2025"
      const dateMatch = line.match(/^(\d{1,2})\s+de\s+(\w{3})\.\s+(\d{4})\s+(.+)/)
      if (!dateMatch) continue

      const [, day, monthAbbr, year, rest] = dateMatch
      const month = monthMap[monthAbbr.toLowerCase()]
      if (month === undefined) continue

      // The rest of the line contains: Description - Beneficiary Value
      // Match value at the end: "- R$ 632,31" or "- + R$ 319,90"
      const valueMatch = rest.match(/([+-]?\s*R\$\s*[\d.,]+)\s*$/)
      if (!valueMatch) continue

      const amountInCents = parseAmountToCents(valueMatch[1])
      if (typeof amountInCents !== 'number' || amountInCents >= 0) continue

      // Extract description (everything before the last " - ")
      const lastDashIndex = rest.lastIndexOf(' - ')
      if (lastDashIndex === -1) continue

      const description = rest.substring(0, lastDashIndex).trim()
      if (!description) continue

      // Parse date
      const date = new Date(parseInt(year), month, parseInt(day))

      expenses.push({
        id: randomUUID(),
        description: description.replace(/\(Parcela.*?\)/i, '').trim(),
        amountInCents: String(Math.abs(amountInCents)),
        date: date.toISOString(),
        raw: line,
        splitType: 'proportional',
      })
    }

    return { expenses }
  }
}
