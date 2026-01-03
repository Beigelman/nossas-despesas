import { ImportedExpenseDraft } from '@/types/expense-import'

export interface ParseResult {
  expenses: ImportedExpenseDraft[]
  errors?: string[]
}

export interface FileInfo {
  content: string | Uint8Array
  filename?: string
  mimeType?: string
}

export abstract class ExpenseParser {
  abstract readonly name: string
  abstract readonly supportedFormats: string[]

  /**
   * Determines if this parser can handle the given file
   */
  abstract canParse(fileInfo: FileInfo): boolean

  /**
   * Parses the file and returns expenses
   */
  abstract parse(fileInfo: FileInfo): Promise<ParseResult>
}
