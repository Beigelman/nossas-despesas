import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

import { isSameDate } from './date'

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

const currencyFormatter = new Intl.NumberFormat('pt-BR', {
  style: 'currency',
  currency: 'BRL',
}).format

function formatCurrency(valueInCents?: number) {
  if (valueInCents === undefined) {
    return currencyFormatter(0)
  }
  return currencyFormatter(valueInCents / 100)
}

const moneyMask = (value = '') => {
  if (value === '') {
    return '0,00'
  }

  value = value.replace('.', '').replace(',', '').replace(/\D/g, '')
  return new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL', minimumFractionDigits: 2 }).format(
    parseFloat(value) / 100,
  )
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function findDifferences(obj1: Record<string, any>, obj2: Record<string, any>, basePath = '') {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  let differences = {} as Record<string, any>

  Object.keys(obj1).forEach((key) => {
    const path = basePath ? `${basePath}.${key}` : key

    // If key is not present in obj2, record it as a difference
    if (!(key in obj2)) {
      differences[path] = { old: obj1[key], actual: undefined }
      return
    }
    if (obj1[key] instanceof Date && obj2[key] instanceof Date && !isSameDate(obj1[key], obj2[key])) {
      differences[path] = { old: obj1[key], actual: obj2[key] }
    } else if (
      typeof obj1[key] === 'object' &&
      typeof obj2[key] === 'object' &&
      obj1[key] !== null &&
      obj2[key] !== null
    ) {
      const deeperDifferences = findDifferences(obj1[key], obj2[key], path)
      differences = { ...differences, ...deeperDifferences }
    } else if (obj1[key] !== obj2[key]) {
      // If values are different, record the difference
      differences[path] = { old: obj1[key], actual: obj2[key] }
    }
  })

  // Check for keys in obj2 that are not in obj1
  Object.keys(obj2).forEach((key) => {
    const path = basePath ? `${basePath}.${key}` : key
    if (!(key in obj1)) {
      differences[path] = { old: undefined, actual: obj2[key] }
    }
  })

  return differences
}

export { cn, formatCurrency, moneyMask, findDifferences }
