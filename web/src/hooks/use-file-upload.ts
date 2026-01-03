import { useCallback, useState } from 'react'

import { ImportedExpenseDraft } from '@/types/expense-import'

interface UseFileUploadReturn {
  uploadFile: (file: File) => Promise<ImportedExpenseDraft[]>
  isUploading: boolean
  error: string | undefined
  fileName: string | undefined
  resetUpload: () => void
}

export function useFileUpload(): UseFileUploadReturn {
  const [isUploading, setIsUploading] = useState(false)
  const [error, setError] = useState<string>()
  const [fileName, setFileName] = useState<string>()

  const uploadFile = useCallback(async (file: File): Promise<ImportedExpenseDraft[]> => {
    setIsUploading(true)
    setError(undefined)

    try {
      const formData = new FormData()
      formData.append('file', file)

      const response = await fetch('/api/expenses/import', {
        method: 'POST',
        body: formData,
      })

      if (!response.ok) {
        const message = (await response.json().catch(() => null))?.message ?? 'Falha ao processar o arquivo enviado.'
        throw new Error(message)
      }

      const data = (await response.json()) as { expenses: ImportedExpenseDraft[] }
      setFileName(file.name)

      return data.expenses
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Não foi possível ler o arquivo selecionado.'
      setError(message)
      throw error
    } finally {
      setIsUploading(false)
    }
  }, [])

  const resetUpload = useCallback(() => {
    setFileName(undefined)
    setError(undefined)
    setIsUploading(false)
  }, [])

  return {
    uploadFile,
    isUploading,
    error,
    fileName,
    resetUpload,
  }
}
