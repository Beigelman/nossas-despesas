type ImportedExpenseDraft = {
  id: string
  description: string
  amountInCents: string
  date?: string
  raw?: string
  splitType: 'equal' | 'proportional' | 'transfer'
}

type ExpenseDraftRow = ImportedExpenseDraft & {
  include: boolean
  categoryId: number
  payerId?: number
  status: 'idle' | 'saving' | 'success' | 'error'
  errorMessage?: string
}

export type { ImportedExpenseDraft, ExpenseDraftRow }
