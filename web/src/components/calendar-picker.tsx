import { format } from 'date-fns'
import { CalendarIcon } from 'lucide-react'
import { useState } from 'react'
import { SelectSingleEventHandler } from 'react-day-picker'

import { cn } from '@/lib/utils'

import { Button } from './ui/button'
import { Calendar } from './ui/calendar'
import { Popover, PopoverContent, PopoverTrigger } from './ui/popover'

type CalendarPickerProps = {
  value: Date
  setValue: (date: Date) => void
}

function CalendarPicker({ value, setValue }: CalendarPickerProps) {
  const [open, setOpen] = useState(false)
  const handleSelect: SelectSingleEventHandler = (_, selectedDay) => {
    setValue(selectedDay)
  }

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant={'outline'}
          className={cn('w-[240px] pl-3 text-left font-normal', !value && 'text-muted-foreground')}
        >
          {value ? format(value, 'PPP') : <span>Escolha uma data</span>}
          <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-auto p-0" align="start">
        <Calendar
          mode="single"
          selected={value}
          onSelect={handleSelect}
          onDayClick={() => setOpen(false)}
          disabled={(date) => date > new Date() || date < new Date('1900-01-01')}
          initialFocus
        />
      </PopoverContent>
    </Popover>
  )
}

export { CalendarPicker }
