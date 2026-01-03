import { ExpenseParser, FileInfo, ParseResult } from './types'

/**
 * Registry that manages multiple parsers and selects the appropriate one
 * for a given file. Parsers are tried in the order they were registered.
 */
export class ParserRegistry {
  private parsers: ExpenseParser[] = []

  /**
   * Registers a new parser
   * Parsers registered first have higher priority
   */
  register(parser: ExpenseParser): void {
    this.parsers.push(parser)
  }

  /**
   * Finds the first parser that can handle the given file
   */
  findParser(fileInfo: FileInfo): ExpenseParser | undefined {
    return this.parsers.find((parser) => parser.canParse(fileInfo))
  }

  /**
   * Parses a file using the first compatible parser
   * Returns empty expenses array if no parser can handle the file
   */
  async parse(fileInfo: FileInfo): Promise<ParseResult> {
    const parser = this.findParser(fileInfo)

    if (!parser) {
      return {
        expenses: [],
        errors: ['No compatible parser found for this file type'],
      }
    }

    return parser.parse(fileInfo)
  }

  /**
   * Gets all registered parsers
   */
  getParsers(): ExpenseParser[] {
    return [...this.parsers]
  }
}
