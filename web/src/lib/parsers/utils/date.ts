/**
 * Parses a date string in various formats to ISO string
 * Handles formats: DD/MM/YYYY, DD-MM-YYYY, DD.MM.YYYY, MM/DD/YYYY, etc.
 */
export function parseDateString(value: string): string | undefined {
  if (!value) {
    return undefined
  }

  const sanitized = value.replace(/[^\d/.-]/g, '')
  const dateParts = sanitized.split(/[./-]/).filter(Boolean)

  if (dateParts.length < 3) {
    return undefined
  }

  let [first, second, third] = dateParts
  let day = first
  let month = second

  if (Number(first) > 12 && Number(second) <= 12) {
    day = first
    month = second
  } else if (Number(second) > 12 && Number(first) <= 12) {
    day = second
    month = first
  }

  if (third.length === 2) {
    const currentYear = new Date().getFullYear()
    const century = currentYear - (currentYear % 100)
    const year = Number(third)
    third = `${year > 50 ? century - 100 : century}${third.padStart(2, '0')}`
  }

  const isoDate = `${third.padStart(4, '0')}-${month.padStart(2, '0')}-${day.padStart(2, '0')}T00:00:00.000Z`
  const parsedDate = new Date(isoDate)

  if (Number.isNaN(parsedDate.getTime())) {
    return undefined
  }

  return parsedDate.toISOString()
}
