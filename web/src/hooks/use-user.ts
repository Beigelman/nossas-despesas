import { useMutation } from '@tanstack/react-query'
import { signIn, useSession } from 'next-auth/react'
import { toast } from 'sonner'

import { User } from '@/domain/user'
import { privateHttpClient } from '@/http/private-client'

function useUser() {
  const { data: session, status, update } = useSession()

  const updateUser = async (user: User) => await update({ user })

  const { mutate: createGroup } = useMutation({
    mutationFn: async (name: string) => await privateHttpClient.post('/group', { name }),
    onSuccess: async () => {
      const resp = await signIn('refresh-token', { redirect: false, refreshToken: session?.refreshToken })
      if (resp?.error) {
        return toast.error(`Falha criar o grupo: ${resp.error}`)
      } else {
        toast.success('Grupo criado com sucesso')
      }
    },
    onError: (error) => toast.error(`Falha criar o grupo: ${error.message}`),
  })

  return {
    user: session?.user,
    createGroup,
    update: updateUser,
    token: session?.token,
    isLoading: status === 'loading',
  }
}

export { useUser }
