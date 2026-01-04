'use client'

import { Loader2, UploadCloudIcon } from 'lucide-react'
import { useTranslations } from 'next-intl'
import { type ChangeEvent, type RefObject } from 'react'

import { Label } from '@/components/ui/label'

interface FileUploadZoneProps {
  fileName?: string
  isUploading: boolean
  onFileChange: (event: ChangeEvent<HTMLInputElement>) => void
  inputRef: RefObject<HTMLInputElement | null>
}

const ACCEPTED_FILE_TYPES = '.csv,.pdf'

export function FileUploadZone({ fileName, isUploading, onFileChange, inputRef }: FileUploadZoneProps) {
  const t = useTranslations()

  return (
    <div className="rounded-lg border border-dashed border-muted-foreground/50 p-4">
      <div className="flex flex-col gap-3 text-sm md:flex-row md:items-center md:justify-between">
        <div className="flex items-center gap-3">
          <UploadCloudIcon className="h-5 w-5 text-muted-foreground" />
          <div>
            <p className="font-medium">{fileName ?? t('expenses.dragDropFile')}</p>
            <p className="text-muted-foreground">{t('expenses.fileUploadDescription')}</p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <input
            ref={inputRef}
            type="file"
            accept={ACCEPTED_FILE_TYPES}
            className="hidden"
            id="expense-import-file"
            onChange={onFileChange}
          />
          <Label
            htmlFor="expense-import-file"
            className="cursor-pointer rounded-md border border-input px-3 py-2 text-sm font-medium hover:bg-muted"
          >
            {t('expenses.selectFile')}
          </Label>
          {isUploading && <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />}
        </div>
      </div>
    </div>
  )
}
