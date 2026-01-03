import { CalendarDaysIcon } from 'lucide-react'
import { useState } from 'react'

import { Button } from '@/components/ui/button'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { cn } from '@/lib/utils'

type YearPickerProps = {
  onSelect: (year: number) => void
  selectedYear: number
}

function YearPicker({ selectedYear, onSelect }: YearPickerProps) {
  const [open, setOpen] = useState(false)
  const years = Array.from({ length: 6 }, (_, i) => new Date().getFullYear() - i)

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button className="w-[240px] justify-start text-left" variant="outline">
          <CalendarDaysIcon className="mr-1 h-4 w-4 -translate-x-1" />
          {selectedYear ?? 'Escolha um ano'}
        </Button>
      </PopoverTrigger>
      <PopoverContent align="start" className="flex w-auto flex-col space-y-2 p-2">
        <div className="grid grid-cols-3 gap-2">
          {years.map((year) => (
            <Button
              key={year}
              className={cn(selectedYear === year && 'bg-accent text-accent-foreground', 'rounded-lg')}
              variant="ghost"
              onClick={() => {
                onSelect(year)
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

export { YearPicker }
