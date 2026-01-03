import { User } from './user'

type Group = {
  id: number
  name: string
  members: User[]
  createdAt: Date
  updatedAt: Date
}

type GroupBalance = {
  userId: number
  balance: number
}

export type { Group, GroupBalance }
