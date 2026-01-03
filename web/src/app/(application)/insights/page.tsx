'use client'

import { lastDayOfMonth, setDate } from 'date-fns'
import { useState } from 'react'
import { DateRange } from 'react-day-picker'

import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { YearPicker } from '@/components/year-picker'

import { CalendarDateRangePicker } from '../../../components/date-rage-picker'
import { ExpensesPeriodInsight } from './components/expenses-period-insight'

export default function InsightsPage() {
  const [dateRange, setDateRange] = useState<DateRange>({
    from: setDate(new Date().setHours(0, 0, 0, 0), 1),
    to: lastDayOfMonth(new Date().setHours(0, 0, 0, 0)),
  })

  const [selectedYear, setSelectedYear] = useState(new Date().getFullYear())

  return (
    <div className="mx-auto my-10 w-full md:w-[786px] lg:w-[1024px]">
      <Tabs defaultValue="period" className="space-y-4 px-2 md:p-0">
        <div className="flex flex-col items-center md:flex-row md:justify-between">
          <TabsList>
            <TabsTrigger value="period">Período</TabsTrigger>
            <TabsTrigger value="yearly">Anual</TabsTrigger>
          </TabsList>
          <TabsContent value="period">
            <CalendarDateRangePicker dateRange={dateRange} setDateRange={setDateRange} />
          </TabsContent>
          <TabsContent value="yearly">
            <YearPicker onSelect={setSelectedYear} selectedYear={selectedYear} />
          </TabsContent>
        </div>
        <TabsContent value="period" className="space-y-4">
          <ExpensesPeriodInsight dateRange={dateRange} aggregate="day" />
        </TabsContent>
        <TabsContent value="yearly" className="space-y-4">
          <ExpensesPeriodInsight
            dateRange={{
              from: new Date(selectedYear, 0, 1, 0, 0, 0, 0),
              to: new Date(selectedYear, 11, 31, 0, 0, 0, 0),
            }}
            aggregate="month"
          />
        </TabsContent>
      </Tabs>
    </div>
  )
}
