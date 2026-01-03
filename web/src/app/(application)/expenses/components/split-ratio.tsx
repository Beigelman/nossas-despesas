import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { cn } from '@/lib/utils'

type SplitRatioProps = {
  payerRatio: number
  receiverRatio: number
  className?: string
}

function SplitRatio({ className, receiverRatio, payerRatio }: SplitRatioProps) {
  return (
    <Tooltip>
      <TooltipTrigger className={cn('', className)}>
        <span className="text-xs text-slate-400">
          {payerRatio}/{receiverRatio}
        </span>
      </TooltipTrigger>
      <TooltipContent>
        Razão da divisão das despesa, sendo o primeiro valor a parte do pagador e o segundo da contra parte
      </TooltipContent>
    </Tooltip>
  )
}

export { SplitRatio }
