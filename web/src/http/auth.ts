import { User } from '@/domain/user'

import { publicHttpClient } from './client'

type LoginResponse = {
  statusCode: number
  data: {
    user: {
      id: number
      name: string
      email: string
      profile_picture?: string
      group_id?: number
      flags?: string[]
      created_at: string
      updated_at: string
    }
    token: string
    refresh_token: string
  }
  date: string
}

type LoginUser = {
  user: User
  token: string
  refreshToken: string
}

async function signInWithCredentials(email: string, password: string): Promise<LoginUser | null> {
  try {
    const {
      data: { data },
    } = await publicHttpClient.post<LoginResponse>('/auth/sign-in/credentials', {
      email,
      password,
    })

    return {
      user: {
        id: data.user.id,
        name: data.user.name,
        email: data.user.email,
        profileImage: data.user.profile_picture,
        groupId: data.user.group_id,
        flags: data.user.flags,
        createdAt: new Date(data.user.created_at),
        updatedAt: new Date(data.user.updated_at),
      },
      token: data.token,
      refreshToken: data.refresh_token,
    }
  } catch (err) {
    return null
  }
}

async function signInWithGoogle(token: string): Promise<LoginUser | null> {
  try {
    const {
      data: { data },
    } = await publicHttpClient.post<LoginResponse>('/auth/sign-in/google', {
      token,
    })
    return {
      user: {
        id: data.user.id,
        name: data.user.name,
        email: data.user.email,
        profileImage: data.user.profile_picture,
        groupId: data.user.group_id,
        flags: data.user.flags,
        createdAt: new Date(data.user.created_at),
        updatedAt: new Date(data.user.updated_at),
      },
      token: data.token,
      refreshToken: data.refresh_token,
    }
  } catch (err) {
    return null
  }
}

type CreateUserRequest = {
  name: string
  email: string
  password: string
  passwordConfirmation: string
}

async function signUpWithCredentials(payload: CreateUserRequest): Promise<LoginUser | null> {
  try {
    const {
      data: { data },
    } = await publicHttpClient.post<LoginResponse>('/auth/sign-up/credentials', {
      name: payload.name,
      email: payload.email,
      password: payload.password,
      confirm_password: payload.passwordConfirmation,
    })

    return {
      user: {
        id: data.user.id,
        name: data.user.name,
        email: data.user.email,
        profileImage: data.user.profile_picture,
        groupId: data.user.group_id,
        flags: data.user.flags,
        createdAt: new Date(data.user.created_at),
        updatedAt: new Date(data.user.updated_at),
      },
      token: data.token,
      refreshToken: data.refresh_token,
    }
  } catch (err) {
    return null
  }
}

async function refreshToken(token: string): Promise<LoginUser | null> {
  try {
    const {
      data: { data },
    } = await publicHttpClient.post<LoginResponse>('/auth/refresh-token', {
      refresh_token: token,
    })

    return {
      user: {
        id: data.user.id,
        name: data.user.name,
        email: data.user.email,
        profileImage: data.user.profile_picture,
        groupId: data.user.group_id,
        createdAt: new Date(data.user.created_at),
        updatedAt: new Date(data.user.updated_at),
      },
      token: data.token,
      refreshToken: data.refresh_token,
    }
  } catch (err) {
    return null
  }
}

export type { LoginUser }
export { refreshToken, signInWithCredentials, signUpWithCredentials, signInWithGoogle }
