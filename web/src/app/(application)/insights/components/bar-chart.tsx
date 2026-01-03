'use client'

import React from 'react'
import { Bar, BarChart as ReBarChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts'

import { Spinning } from '@/components/ui/spinning'
import { formatCurrency } from '@/lib/utils'

type BarChartProps = {
  customTooltip?: React.ReactElement
  data: unknown[] | undefined
  isLoading?: boolean
  direction?: 'horizontal' | 'vertical'
  XKey: string
  YKey: string
}

export function BarChart({ data, isLoading, XKey, YKey, direction, customTooltip }: BarChartProps) {
  if (isLoading) {
    return (
      <ResponsiveContainer width="100%" height={350} className="flex items-center justify-center">
        <Spinning className="h-24 w-24 border-8" />
      </ResponsiveContainer>
    )
  }

  return (
    <ResponsiveContainer width="100%" height={350}>
      <ReBarChart data={data} layout={direction ?? 'horizontal'}>
        <XAxis
          dataKey={direction === 'vertical' ? undefined : XKey}
          type={direction === 'vertical' ? 'number' : 'category'}
          stroke="#888888"
          fontSize={10}
          tickLine={false}
          axisLine={false}
          tickFormatter={direction === 'vertical' ? (value) => formatCurrency(value) : undefined}
        />
        <YAxis
          dataKey={direction === 'vertical' ? XKey : undefined}
          type={direction === 'vertical' ? 'category' : 'number'}
          stroke="#888888"
          fontSize={10}
          tickLine={false}
          axisLine={false}
          tickFormatter={direction === 'vertical' ? undefined : (value) => formatCurrency(value)}
        />
        {customTooltip || <Tooltip />}
        <Bar
          dataKey={YKey}
          fill="currentColor"
          radius={direction === 'vertical' ? [0, 5, 5, 0] : [5, 5, 0, 0]}
          className="fill-primary"
        />
      </ReBarChart>
    </ResponsiveContainer>
  )
}
