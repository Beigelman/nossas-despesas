import { AxiosError } from 'axios'

type ApiResponse<T> = {
  statusCode: number
  data: T
  date: string
}

type ApiError = AxiosError<{ error: string; message: string; status: number }>
