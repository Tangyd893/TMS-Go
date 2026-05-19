import { request, type PageResult } from '../request'

export interface TransportTask {
  id: string; taskNo: string; orderId: string; vehicleId: string; driverId: string
  driverName: string; plateNo: string; originName: string; destName: string
  status: string; actualDepartTime: string; actualArriveTime: string; createdAt: string
  remark: string
}

export function listTransportTasks(params: Record<string, any>): Promise<PageResult<TransportTask>> {
  return request.get('/transport/tasks', { params })
}

export function getTransportTask(id: string): Promise<TransportTask> {
  return request.get(`/transport/tasks/${id}`)
}

export function createTransportTask(data: any): Promise<TransportTask> {
  return request.post('/transport/tasks', data)
}

export function departTask(id: string): Promise<TransportTask> {
  return request.post(`/transport/tasks/${id}/depart`)
}

export function arriveTask(id: string): Promise<TransportTask> {
  return request.post(`/transport/tasks/${id}/arrive`)
}

export function signTask(id: string): Promise<TransportTask> {
  return request.post(`/transport/tasks/${id}/sign`)
}

export interface TransportNode {
  id: string; taskId: string; nodeType: string; nodeName: string
  locationName: string; arrivedAt: string; departedAt: string; remark: string
  createdAt: string
}

export function addNode(taskId: string, data: any): Promise<TransportNode> {
  return request.post(`/transport/tasks/${taskId}/nodes`, data)
}

export function getNodes(taskId: string): Promise<TransportNode[]> {
  return request.get(`/transport/tasks/${taskId}/nodes`)
}

export interface Receipt {
  id: string; taskId: string; orderId: string; receiptNo: string
  signBy: string; signAt: string; signImageUrl: string; remark: string
}

export function createReceipt(taskId: string, data: any): Promise<Receipt> {
  return request.post(`/transport/tasks/${taskId}/receipt`, data)
}

export function getReceipt(taskId: string): Promise<Receipt> {
  return request.get(`/transport/tasks/${taskId}/receipt`)
}
