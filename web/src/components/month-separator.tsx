function monthToString(month?: number): string {
  switch (month) {
    case 0:
      return 'Janeiro'
    case 1:
      return 'Fevereiro'
    case 2:
      return 'Março'
    case 3:
      return 'Abril'
    case 4:
      return 'Maio'
    case 5:
      return 'Junho'
    case 6:
      return 'Julho'
    case 7:
      return 'Agosto'
    case 8:
      return 'Setembro'
    case 9:
      return 'Outubro'
    case 10:
      return 'Novembro'
    case 11:
      return 'Dezembro'
    default:
      return ''
  }
}

function MonthSeparator({ currentDate, nextDate }: { currentDate: Date; nextDate?: Date }): React.ReactNode {
  if (
    nextDate &&
    (nextDate?.getMonth() < currentDate.getMonth() || nextDate?.getFullYear() < currentDate.getFullYear())
  ) {
    return (
      <div className="bg-zinc-300 px-2 text-sm uppercase dark:bg-zinc-800">
        {`${monthToString(nextDate?.getMonth())} de ${nextDate?.getFullYear()}`}
      </div>
    )
  }

  return null
}

export { MonthSeparator }
