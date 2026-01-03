import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { Group } from '@/domain/group'
import { User } from '@/domain/user'
import { privateHttpClient } from '@/http/private-client'
import { ApiResponse } from '@/lib/api'

import { useUser } from './use-user'

type GetGroupApiResponse = ApiResponse<{
  id: number
  name: string
  members: {
    id: number
    name: string
    group_id?: number
    email: string
    profile_picture?: string
    created_at: string
    updated_at: string
  }[]
  created_at: string
  updated_at: string
}>

function useGroup() {
  const { user } = useUser()
  const queryClient = useQueryClient()

  const { data, isError, isPending } = useQuery({
    enabled: user?.groupId !== undefined,
    queryKey: ['group'],
    queryFn: async (): Promise<Group> => {
      const {
        data: { data },
      } = await privateHttpClient.get<GetGroupApiResponse>(`/group`)

      const members = data.members.map<User>((member) => ({
        id: member.id,
        name: member.name,
        email: member.email,
        groupId: member.group_id,
        profileImage: member.profile_picture,
        createdAt: new Date(member.created_at),
        updatedAt: new Date(member.updated_at),
      }))

      return {
        id: data.id,
        name: data.name,
        members,
        createdAt: new Date(data.created_at),
        updatedAt: new Date(data.updated_at),
      }
    },
  })

  const { mutate: inviteUserToGroup } = useMutation({
    mutationFn: async (email: string) =>
      await privateHttpClient.post('/group/invite', { email, base_url: process.env.NEXT_PUBLIC_BASE_URL }),
    onSuccess: async () => toast.success('Convite enviado com sucesso!'),
    onError: (error) => toast.error(`Falha ao enviar convite: ${error.message}`),
  })

  const getMember = (userId: number): User | undefined => {
    return data?.members.find((member) => member.id === userId)
  }

  const me = user

  const partner = data?.members.find((member) => member.id !== user?.id)

  const refresh = () => queryClient.invalidateQueries({ queryKey: ['group'] })

  return {
    group: data,
    refresh,
    inviteUserToGroup,
    getMember,
    me,
    partner,
    isLoading: isPending,
    isError,
  }
}

export { useGroup }
