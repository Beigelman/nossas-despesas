import axios from 'axios'
import { getSession, signIn } from 'next-auth/react'

const privateHttpClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  headers: {
    Accept: 'application/json',
    'Content-Type': 'application/json',
  },
})

privateHttpClient.interceptors.request.use(
  async (config) => {
    const session = await getSession()
    const token = session?.token

    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }

    return config
  },
  (error) => Promise.reject(error),
)

privateHttpClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config
    if (error?.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      try {
        const session = await getSession()
        const resp = await signIn('refresh-token', { redirect: false, refreshToken: session?.refreshToken })
        if (resp?.error) {
          return Promise.reject(resp.error)
        }

        const newSession = await getSession()
        originalRequest.headers.Authorization = `Bearer ${newSession?.token}`

        return privateHttpClient(originalRequest)
      } catch (e) {
        return Promise.reject(e)
      }
    }

    return Promise.reject(error)
  },
)

export { privateHttpClient }
