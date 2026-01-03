import { z } from 'zod'

const expenseFormSchema = z.object({
  description: z.string().min(3, 'Descrição muito curta'),
  amount: z.string().min(1, 'Valor inválido. Deve ser maior que 0'),
  payerId: z.number().min(1, 'Pagador inválido'),
  splitType: z.enum(['equal', 'proportional', 'transfer']),
  categoryId: z.number(),
  date: z.date(),
  refundAmount: z.string(),
})

const incomeFormSchema = z.object({
  userId: z.number(),
  amount: z.string().min(1, 'Valor inválido'),
  date: z.date(),
  type: z.enum(['salary', 'benefit', 'vacation', 'thirteenth_salary', 'other']),
})

const signUpFormSchema = z.object({
  name: z.string().min(3, 'Descrição muito curta'),
  email: z.string().email('Email inválido'),
  password: z.string().min(6, 'Senha muito curta'),
  passwordConfirmation: z.string().min(6, 'Senha muito curta'),
})

const signInFormSchema = z.object({
  email: z.string().email('Email inválido'),
  password: z.string().min(6, 'Senha muito curta'),
})

export { signInFormSchema, signUpFormSchema, expenseFormSchema, incomeFormSchema }
