import { redirect } from 'next/navigation'

import { getServerSession } from '@/lib/auth/get-server-session'

type ExpensesLayoutProps = {
  children: React.ReactNode
}

export default async function ExpensesLayout({ children }: ExpensesLayoutProps) {
  const session = await getServerSession()

  if (!session?.user.groupId) {
    redirect('/group')
  }

  return <>{children}</>
}
