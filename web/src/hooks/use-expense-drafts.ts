import { useCallback, useMemo, useState } from 'react'

import { DEFAULT_CATEGORY_ID } from '@/domain/category'
import { useCategories } from '@/hooks/use-categories'
import { useGroup } from '@/hooks/use-group'
import { usePredict } from '@/hooks/use-predict'
import { ExpenseDraftRow, ImportedExpenseDraft } from '@/types/expense-import'

interface UseExpenseDraftsReturn {
  drafts: ExpenseDraftRow[]
  selectedDrafts: ExpenseDraftRow[]
  draftsToPersist: ExpenseDraftRow[]
  predictingCategories: Set<string>
  initializeDrafts: (expenses: ImportedExpenseDraft[]) => void
  updateDraft: (draftId: string, updater: Partial<ExpenseDraftRow>) => void
  resetDrafts: () => void
}

export function useExpenseDrafts(): UseExpenseDraftsReturn {
  const [drafts, setDrafts] = useState<ExpenseDraftRow[]>([])
  const [predictingCategories, setPredictingCategories] = useState<Set<string>>(new Set())

  const { me, partner } = useGroup()
  const { categories } = useCategories()
  const { predictCategoryID } = usePredict()

  const flatCategories = useMemo(() => categories.flatMap((group) => group.categories), [categories])

  const initializeDrafts = useCallback(
    (expenses: ImportedExpenseDraft[]) => {
      const defaultCategory = flatCategories[0]?.id ?? DEFAULT_CATEGORY_ID

      const parsedDrafts = expenses.map<ExpenseDraftRow>((expense) => ({
        ...expense,
        include: true,
        categoryId: defaultCategory,
        payerId: me?.id ?? partner?.id,
        status: 'idle',
      }))

      setDrafts(parsedDrafts)

      // Predict categories for all imported expenses
      parsedDrafts.forEach((draft) => {
        setPredictingCategories((prev) => new Set(prev).add(draft.id))

        predictCategoryID({
          name: draft.description,
          amount_cents: parseInt(draft.amountInCents) || 0,
        })
          .then((response) => {
            const predictedCategoryId = response.data.category_id
            if (typeof predictedCategoryId === 'number') {
              updateDraft(draft.id, { categoryId: predictedCategoryId })
            }
          })
          .catch((error) => {
            console.error('Failed to predict category for:', draft.description, error)
          })
          .finally(() => {
            setPredictingCategories((prev) => {
              const next = new Set(prev)
              next.delete(draft.id)
              return next
            })
          })
      })
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [flatCategories, me?.id, partner?.id, predictCategoryID],
  )

  const updateDraft = useCallback((draftId: string, updater: Partial<ExpenseDraftRow>) => {
    setDrafts((previous) =>
      previous.map((draft) => {
        if (draft.id !== draftId) {
          return draft
        }

        const hasExplicitStatusUpdate = updater.status !== undefined
        const computedStatus =
          updater.status ?? (draft.status === 'success' || draft.status === 'error' ? 'idle' : draft.status)

        const nextDraft: ExpenseDraftRow = {
          ...draft,
          ...updater,
          status: computedStatus,
          errorMessage: updater.errorMessage !== undefined ? updater.errorMessage : draft.errorMessage,
        }

        if (!hasExplicitStatusUpdate && (draft.status === 'success' || draft.status === 'error')) {
          nextDraft.status = 'idle'
          nextDraft.errorMessage = updater.errorMessage !== undefined ? updater.errorMessage : undefined
        }

        return nextDraft
      }),
    )
  }, [])

  const resetDrafts = useCallback(() => {
    setDrafts([])
    setPredictingCategories(new Set())
  }, [])

  const selectedDrafts = useMemo(() => drafts.filter((draft) => draft.include), [drafts])
  const draftsToPersist = useMemo(() => selectedDrafts.filter((draft) => draft.status !== 'success'), [selectedDrafts])

  return {
    drafts,
    selectedDrafts,
    draftsToPersist,
    predictingCategories,
    initializeDrafts,
    updateDraft,
    resetDrafts,
  }
}
