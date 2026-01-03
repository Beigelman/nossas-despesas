'use client'

import { BarChartBigIcon, DollarSign, LogInIcon } from 'lucide-react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'

import { cn } from '@/lib/utils'

import { Button } from './ui/button'

type SidebarProps = React.HTMLAttributes<HTMLDivElement>

export function Sidebar({ className }: SidebarProps) {
  const path = usePathname()

  return (
    <div className={cn('flex flex-col gap-2 p-2', className)}>
      <Link href={'/expenses'}>
        <Button
          variant={path === '/expenses' ? 'secondary' : 'ghost'}
          className="flex w-full items-center justify-start"
        >
          <DollarSign size={20} />
          <span className="ml-2 text-sm">Despesas</span>
        </Button>
      </Link>
      <Link href={'/incomes'}>
        <Button
          variant={path === '/incomes' ? 'secondary' : 'ghost'}
          className="flex w-full items-center justify-start"
        >
          <LogInIcon size={20} />
          <span className="ml-2 text-sm">Receitas</span>
        </Button>
      </Link>
      <Link href={'/insights'}>
        <Button
          variant={path === '/insights' ? 'secondary' : 'ghost'}
          className="flex w-full items-center justify-start"
        >
          <BarChartBigIcon size={20} />
          <span className="ml-2 text-sm">Insights</span>
        </Button>
      </Link>
    </div>
  )
}
