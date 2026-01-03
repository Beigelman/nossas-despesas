import { NextRequest, NextResponse } from 'next/server'

import { defaultParserRegistry } from '@/lib/parsers'

export const runtime = 'nodejs'

export async function POST(request: NextRequest) {
  try {
    const formData = await request.formData()
    const file = formData.get('file')

    if (!file || !(file instanceof Blob)) {
      return NextResponse.json({ message: 'Arquivo inválido' }, { status: 400 })
    }

    const filename = 'name' in file ? file.name : undefined
    const mimeType = file.type
    const arrayBuffer = await file.arrayBuffer()

    // Determine if it's a PDF or text-based file
    const isPdf = mimeType === 'application/pdf' || filename?.toLowerCase().endsWith('.pdf')
    let content: string | Uint8Array

    if (isPdf) {
      try {
        // Use unpdf - serverless-optimized PDF library
        const { extractText } = await import('unpdf')
        const { text } = await extractText(new Uint8Array(arrayBuffer))
        content = Array.isArray(text) ? text.join('\n') : text
      } catch (error) {
        console.error('PDF extraction error:', error)
        return NextResponse.json(
          {
            message:
              'Não foi possível processar o arquivo PDF. Por favor, tente converter para CSV ou entre em contato com o suporte.',
          },
          { status: 400 },
        )
      }
    } else {
      content = Buffer.from(arrayBuffer).toString('utf-8')
    }

    // Use the parser registry to parse the file
    const result = await defaultParserRegistry.parse({
      content,
      filename,
      mimeType,
    })

    if (result.expenses.length === 0) {
      return NextResponse.json(
        {
          message: 'Não foi possível extrair despesas do arquivo. Verifique se o formato está correto.',
        },
        { status: 400 },
      )
    }

    return NextResponse.json({ expenses: result.expenses })
  } catch (error) {
    console.error('Import error:', error)
    return NextResponse.json({ message: 'Não foi possível processar o arquivo enviado.' }, { status: 500 })
  }
}
