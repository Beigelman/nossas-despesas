'use client'

import { UploadIcon } from 'lucide-react'
import { ChangeEvent, useCallback, useRef, useState } from 'react'

import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { useExpenseDrafts } from '@/hooks/use-expense-drafts'
import { useExpenses } from '@/hooks/use-expenses'
import { useFileUpload } from '@/hooks/use-file-upload'
import { useGroup } from '@/hooks/use-group'

import { FileUploadZone } from '../../../../components/file-upload-zone'
import { ExpenseDraftsTable } from './expense-drafts-table'

export function ImportExpensesButton() {
  const [open, setOpen] = useState(false)
  const [isSaving, setIsSaving] = useState(false)
  const [generalError, setGeneralError] = useState<string>()
  const inputRef = useRef<HTMLInputElement | null>(null)

  const { me, partner } = useGroup()
  const { createExpenseAsync } = useExpenses('')
  const { uploadFile, isUploading, fileName, resetUpload } = useFileUpload()
  const { drafts, selectedDrafts, draftsToPersist, predictingCategories, initializeDrafts, updateDraft, resetDrafts } =
    useExpenseDrafts()

  const handleFileChange = useCallback(
    async (event: ChangeEvent<HTMLInputElement>) => {
      const file = event.target.files?.[0]
      if (!file) {
        return
      }

      setGeneralError(undefined)

      try {
        const expenses = await uploadFile(file)
        initializeDrafts(expenses)
        inputRef.current && (inputRef.current.value = '')
      } catch (error) {
        // Error is already set in useFileUpload
        console.error('File upload failed:', error)
      }
    },
    [uploadFile, initializeDrafts],
  )

  const handleSaveSelected = useCallback(async () => {
    if (!draftsToPersist.length || !me?.id) {
      return
    }

    setGeneralError(undefined)
    setIsSaving(true)

    try {
      for (const draft of draftsToPersist) {
        if (!draft.description || !draft.amountInCents) {
          updateDraft(draft.id, {
            status: 'error',
            errorMessage: 'Revise a descrição e o valor antes de salvar.',
          })
          continue
        }

        if (!draft.payerId) {
          updateDraft(draft.id, {
            status: 'error',
            errorMessage: 'Escolha quem pagou esta despesa.',
          })
          continue
        }

        updateDraft(draft.id, { status: 'saving', errorMessage: undefined })

        try {
          const receiverId = draft.payerId === me.id ? partner?.id ?? 0 : me.id
          const amountInCents = parseInt(draft.amountInCents)
          await createExpenseAsync({
            name: draft.description,
            amount: amountInCents,
            categoryId: draft.categoryId,
            payerId: draft.payerId,
            receiverId,
            splitType: draft.splitType,
            createdAt: draft.date ? new Date(draft.date) : new Date(),
            metadata: { silent: true },
          })

          updateDraft(draft.id, { status: 'success' })
        } catch (error) {
          const message =
            error instanceof Error
              ? error.message
              : 'Erro inesperado ao salvar esta despesa. Tente novamente mais tarde.'
          updateDraft(draft.id, { status: 'error', errorMessage: message })
        }
      }
    } finally {
      setIsSaving(false)
    }
  }, [createExpenseAsync, draftsToPersist, me?.id, partner?.id, updateDraft])

  const resetState = useCallback(() => {
    resetDrafts()
    resetUpload()
    setGeneralError(undefined)
    setIsSaving(false)
    inputRef.current && (inputRef.current.value = '')
  }, [resetDrafts, resetUpload])

  const handleOpenChange = useCallback(
    (nextOpen: boolean) => {
      setOpen(nextOpen)
      if (!nextOpen) {
        resetState()
      }
    },
    [resetState],
  )

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogTrigger asChild>
        <Button variant="outline" className="gap-2">
          <UploadIcon size={16} />
          Importar despesas
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-5xl">
        <DialogHeader>
          <DialogTitle>Importar faturas ou extratos</DialogTitle>
          <DialogDescription>
            Suba um arquivo CSV ou PDF com as movimentações do seu cartão ou banco, revise as informações e salve
            somente o que fizer sentido.
          </DialogDescription>
        </DialogHeader>

        <section className="space-y-4">
          <FileUploadZone
            fileName={fileName}
            isUploading={isUploading}
            onFileChange={handleFileChange}
            inputRef={inputRef}
          />

          {generalError && (
            <div className="rounded-md border border-destructive/40 bg-destructive/10 px-4 py-2 text-sm text-destructive">
              {generalError}
            </div>
          )}

          <ExpenseDraftsTable
            drafts={drafts}
            selectedDrafts={selectedDrafts}
            predictingCategories={predictingCategories}
            onUpdateDraft={updateDraft}
          />
        </section>
        <DialogFooter className="items-center justify-between gap-2 sm:flex-row">
          <div className="text-sm text-muted-foreground">
            {selectedDrafts.length} despesas selecionadas • {draftsToPersist.length} pendentes para importar.
          </div>
          <div className="flex flex-col-reverse gap-2 sm:flex-row sm:gap-3">
            <Button variant="outline" onClick={resetState}>
              Limpar seleção
            </Button>
            <Button onClick={handleSaveSelected} disabled={!draftsToPersist.length || isUploading || isSaving || !me}>
              {isSaving ? 'Salvando...' : 'Salvar despesas selecionadas'}
            </Button>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
