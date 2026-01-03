type Expense = {
  id: number
  name: string
  amount: number
  refundAmount?: number
  description: string
  categoryId: number
  payerId: number
  receiverId: number
  splitRatio: {
    payer: number
    receiver: number
  }
  splitType: 'equal' | 'proportional' | 'transfer'
  createdAt: Date
  updatedAt?: Date
  deletedAt?: Date
  version?: number
}

export type { Expense }
