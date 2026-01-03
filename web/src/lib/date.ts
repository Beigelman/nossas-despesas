function newUTCDate(date?: Date | number): Date {
  if (date === undefined) {
    date = new Date()
  } else if (typeof date === 'number') {
    date = new Date(date)
  }

  return new Date(
    Date.UTC(
      date.getFullYear(),
      date.getMonth(),
      date.getDate(),
      date.getHours(),
      date.getMinutes(),
      date.getSeconds(),
      date.getMilliseconds(),
    ),
  )
}

function isSameDate(date1?: Date, date2?: Date): boolean {
  if (!date1 || !date2) {
    return false
  }

  return (
    date1.getFullYear() === date2.getFullYear() &&
    date1.getMonth() === date2.getMonth() &&
    date1.getDate() === date2.getDate()
  )
}

export { newUTCDate, isSameDate }
