type User = {
  id: number
  name: string
  email: string
  groupId?: number
  profileImage?: string
  flags?: string[]
  createdAt: Date
  updatedAt: Date
}

export type { User }
