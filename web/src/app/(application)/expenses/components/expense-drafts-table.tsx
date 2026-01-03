import { FileTextIcon } from 'lucide-react'

import { Table, TableBody, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { formatCurrency } from '@/lib/utils'
import { ExpenseDraftRow } from '@/types/expense-import'

import { ExpenseDraftRowComponent } from './expense-draft-row'

interface ExpenseDraftsTableProps {
  drafts: ExpenseDraftRow[]
  selectedDrafts: ExpenseDraftRow[]
  predictingCategories: Set<string>
  onUpdateDraft: (draftId: string, updater: Partial<ExpenseDraftRow>) => void
}

export function ExpenseDraftsTable({
  drafts,
  selectedDrafts,
  predictingCategories,
  onUpdateDraft,
}: ExpenseDraftsTableProps) {
  if (drafts.length === 0) {
    return (
      <div className="flex flex-col items-center gap-3 rounded-md border border-dashed border-muted-foreground/40 px-6 py-10 text-center text-sm text-muted-foreground">
        <FileTextIcon className="h-12 w-12 text-muted-foreground/70" />
        <p>
          Selecione um arquivo CSV ou PDF para começarmos a extrair as transações da sua fatura ou extrato bancário.
        </p>
        <p>Após o processamento, você poderá ajustar qualquer valor antes de salvar como novas despesas.</p>
      </div>
    )
  }

  return (
    <div className="space-y-2">
      <div className="flex flex-wrap items-center justify-between gap-2 text-sm text-muted-foreground">
        <p>
          {selectedDrafts.length} de {drafts.length} transações selecionadas para importação.
        </p>
        <p>
          Totais identificados:{' '}
          {formatCurrency(selectedDrafts.reduce((acc, draft) => acc + parseInt(draft.amountInCents), 0))}
        </p>
      </div>
      <div className="max-h-[420px] overflow-auto rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-12 text-center">Usar</TableHead>
              <TableHead>Descrição</TableHead>
              <TableHead className="w-32">Data</TableHead>
              <TableHead className="w-40">Valor</TableHead>
              <TableHead className="w-12">Categoria</TableHead>
              <TableHead className="w-40">Pago por</TableHead>
              <TableHead className="w-12">Distribuir</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {drafts.map((draft) => (
              <ExpenseDraftRowComponent
                key={draft.id}
                draft={draft}
                onUpdate={onUpdateDraft}
                isPredicting={predictingCategories.has(draft.id)}
                isPredicted={!predictingCategories.has(draft.id)}
              />
            ))}
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
