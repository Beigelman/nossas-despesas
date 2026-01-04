'use client'

import { format, getMonth, setMonth, setYear } from 'date-fns'
import { enUS, ptBR } from 'date-fns/locale'
import { CalendarDaysIcon } from 'lucide-react'
import { useLocale, useTranslations } from 'next-intl'
import { useState } from 'react'

import { Button } from '@/components/ui/button'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { cn } from '@/lib/utils'

import { Separator } from './ui/separator'

type YearMonthPickerProps = {
  onSelectDate: (date: Date) => void
  selectedDate: Date
}

function YearMonthPicker({ selectedDate, onSelectDate }: YearMonthPickerProps) {
  const [open, setOpen] = useState(false)
  const locale = useLocale()
  const t = useTranslations()
  const dateLocale = locale === 'en' ? enUS : ptBR
  const years = Array.from({ length: 6 }, (_, i) => new Date().getFullYear() - i)
  const months = [
    { month: locale === 'en' ? 'Jan' : 'Jan', number: 0 },
    { month: locale === 'en' ? 'Feb' : 'Fev', number: 1 },
    { month: locale === 'en' ? 'Mar' : 'Mar', number: 2 },
    { month: locale === 'en' ? 'Apr' : 'Abr', number: 3 },
    { month: locale === 'en' ? 'May' : 'Mai', number: 4 },
    { month: locale === 'en' ? 'Jun' : 'Jun', number: 5 },
    { month: locale === 'en' ? 'Jul' : 'Jul', number: 6 },
    { month: locale === 'en' ? 'Aug' : 'Ago', number: 7 },
    { month: locale === 'en' ? 'Sep' : 'Set', number: 8 },
    { month: locale === 'en' ? 'Oct' : 'Out', number: 9 },
    { month: locale === 'en' ? 'Nov' : 'Nov', number: 10 },
    { month: locale === 'en' ? 'Dec' : 'Dez', number: 11 },
  ]

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button className="w-[240px] justify-start text-left" variant="outline">
          <CalendarDaysIcon className="mr-1 h-4 w-4 -translate-x-1" />
          {selectedDate ? format(selectedDate, 'MMM-yyyy', { locale: dateLocale }) : t('date.chooseDate')}
        </Button>
      </PopoverTrigger>
      <PopoverContent align="start" className="flex w-auto space-x-4 p-2">
        <div className="grid grid-cols-2 gap-2">
          {months.map(({ month, number }) => (
            <Button
              key={number}
              className={cn(getMonth(selectedDate) === number && 'bg-accent text-accent-foreground', 'rounded-lg')}
              variant="ghost"
              onClick={() => {
                onSelectDate(setMonth(selectedDate, number))
                setOpen(false)
              }}
            >
              {month}
            </Button>
          ))}
        </div>
        <Separator orientation="vertical" className="h-auto" />
        <div className="flex flex-col gap-2">
          {years.map((year) => (
            <Button
              key={year}
              className={cn(selectedDate.getFullYear() === year && 'bg-accent text-accent-foreground', 'rounded-lg')}
              variant="ghost"
              onClick={() => {
                onSelectDate(setYear(selectedDate, year))
                setOpen(false)
              }}
            >
              {year}
            </Button>
          ))}
        </div>
      </PopoverContent>
    </Popover>
  )
}

export { YearMonthPicker }
