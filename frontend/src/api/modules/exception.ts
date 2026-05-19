import { request, type PageResult } from '../request'

export interface Exception {
  id: string; taskId: string; orderId: string; exceptionNo: string
  exceptionType: string; description: string; severity: string
  reportBy: string; reportByName: string; status: string
  handlerName: string; handleResult: string; createdAt: string
}

export function listExceptions(params: Record<string, any>): Promise<PageResult<Exception>> {
  return request.get('/transport/exceptions', { params })
}

export function getException(id: string): Promise<Exception> {
  return request.get(`/transport/exceptions/${id}`)
}

export function createException(data: any): Promise<Exception> {
  return request.post('/transport/exceptions', data)
}

export function handleException(id: string, data: any): Promise<void> {
  return request.put(`/transport/exceptions/${id}/handle`, data)
}

export function closeException(id: string): Promise<void> {
  return request.post(`/transport/exceptions/${id}/close`)
}

export interface ExceptionLog {
  id: string; exceptionId: string; action: string; content: string
  operatorName: string; createdAt: string
}

export function getExceptionLogs(id: string): Promise<ExceptionLog[]> {
  return request.get(`/transport/exceptions/${id}/logs`)
}
