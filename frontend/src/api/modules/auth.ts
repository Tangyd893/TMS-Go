import { request } from '../request'

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResult {
  accessToken: string
  refreshToken: string
  user: UserInfo
}

export interface UserInfo {
  id: string
  username: string
  realName: string
  avatarUrl: string
  permissions: string[]
  roles: string[]
}

export function login(data: LoginRequest): Promise<LoginResult> {
  return request.post('/auth/login', data) as Promise<LoginResult>
}

export function refreshToken(refreshToken: string): Promise<LoginResult> {
  return request.post('/auth/refresh', { refreshToken }) as Promise<LoginResult>
}

export function getCurrentUser(): Promise<UserInfo> {
  return request.get('/auth/me') as Promise<UserInfo>
}
