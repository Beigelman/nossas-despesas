'use client'

import { useLocale } from 'next-intl'

type DateDisplayProps = {
  date: Date
  className?: string
}

function DateDisplay({ date }: DateDisplayProps) {
  const locale = useLocale()

  return (
    <div className="mr-2 flex flex-col text-center">
      <span className="text-xl font-bold leading-4 text-primary">{date.getDate().toString().padStart(2, '0')}</span>
      <span className="text-sm leading-4 text-gray-500">
        {date.toLocaleString(locale === 'en' ? 'en-US' : 'pt-BR', { month: 'short' })}
      </span>
    </div>
  )
}

export { DateDisplay }
