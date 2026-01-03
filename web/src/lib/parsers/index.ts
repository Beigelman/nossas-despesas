// Create and export a default registry with all parsers
import { FlashCsvParser } from './implementations/csv-flash.parser'
import { InterCsvParser } from './implementations/csv-inter.parser'
import { InterCreditCardPdfParser } from './implementations/pdf-inter.parser'
import { ParserRegistry } from './registry'

export type { ExpenseParser, ParseResult, FileInfo } from './types'
export { ParserRegistry } from './registry'

// Parser implementations
export { FlashCsvParser } from './implementations/csv-flash.parser'
export { InterCsvParser } from './implementations/csv-inter.parser'
export { InterCreditCardPdfParser } from './implementations/pdf-inter.parser'

// Utilities (exported for potential reuse)
export { parseAmountToCents } from './utils/amount'
export { parseDateString } from './utils/date'
export * from './utils/csv'

/**
 * Default parser registry with all available parsers registered
 * Parsers are tried in this order:
 * 1. Inter Credit Card PDF (most specific)
 * 2. Flash CSV
 * 3. Inter CSV
 */
export const defaultParserRegistry = new ParserRegistry()
defaultParserRegistry.register(new InterCreditCardPdfParser())
defaultParserRegistry.register(new FlashCsvParser())
defaultParserRegistry.register(new InterCsvParser())
