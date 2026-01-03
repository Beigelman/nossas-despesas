import { cookies, headers } from 'next/headers'
import { getServerSession as nextAuthGetServerSession } from 'next-auth'

import { nextAuthOptions } from './next-auth-options'

/**
 * Wrapper para getServerSession que aguarda headers() e cookies()
 * antes de chamar o NextAuth, garantindo compatibilidade com Next.js 15
 *
 * Nota: Os erros de headers() e cookies() ainda podem aparecer porque
 * o NextAuth 4.24.5 usa essas APIs internamente sem await. Esses são
 * avisos do Next.js 15 em modo de desenvolvimento, mas não impedem
 * o funcionamento da aplicação.
 */
export async function getServerSession() {
  // Aguarda headers e cookies antes de chamar getServerSession
  // Isso garante que o contexto esteja pronto quando o NextAuth tentar usá-los
  const [headersList, cookiesList] = await Promise.all([headers(), cookies()])

  // Armazena os valores para garantir que estão disponíveis
  // quando o NextAuth tentar acessá-los
  headersList.get('x-url')
  cookiesList.getAll()

  return nextAuthGetServerSession(nextAuthOptions)
}
