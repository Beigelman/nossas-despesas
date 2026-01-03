import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { Income, incomeLabel } from '@/domain/income'
import { useIncome } from '@/hooks/use-income'
type DeleteIncomeButtonProps = {
  children: React.ReactNode
  income: Income
}

function DeleteIncomeButton({ income, children }: DeleteIncomeButtonProps) {
  const { deleteIncome } = useIncome()

  return (
    <Dialog>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Deseja continuar?</DialogTitle>
          <DialogDescription>{`Deletar receita ${incomeLabel(income.type)}`}</DialogDescription>
        </DialogHeader>
        <DialogFooter className="gap-2">
          <DialogClose>
            <Button variant="secondary" className="w-full md:w-fit">
              Cancelar
            </Button>
          </DialogClose>
          <Button variant="destructive" onClick={() => deleteIncome(income.id)}>
            Deletar
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

export { DeleteIncomeButton }
