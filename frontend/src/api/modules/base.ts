import { request } from '../request'

export interface PageResult<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
}

export function listEntities<T>(prefix: string, params: Record<string, any>): Promise<PageResult<T>> {
  return request.get(prefix, { params }) as Promise<PageResult<T>>
}

export function getEntity<T>(prefix: string, id: string): Promise<T> {
  return request.get(`${prefix}/${id}`) as Promise<T>
}

export function createEntity<T>(prefix: string, data: any): Promise<T> {
  return request.post(prefix, data) as Promise<T>
}

export function updateEntity<T>(prefix: string, id: string, data: any): Promise<T> {
  return request.put(`${prefix}/${id}`, data) as Promise<T>
}

export function deleteEntity(prefix: string, id: string): Promise<void> {
  return request.delete(`${prefix}/${id}`)
}
