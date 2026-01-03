type Income = {
  id: number
  userId: number
  type: 'salary' | 'benefit' | 'vacation' | 'thirteenth_salary' | 'other'
  amount: number
  createdAt: Date
}

function incomeLabel(type: string) {
  switch (type) {
    case 'salary':
      return 'Salário'
    case 'benefit':
      return 'Benefício'
    case 'vacation':
      return 'Férias'
    case 'thirteenth_salary':
      return 'Décimo Terceiro'
    case 'other':
      return 'Outros'
    default:
      return ''
  }
}

export { incomeLabel }
export type { Income }
