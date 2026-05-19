import { request, type PageResult } from '../request'

export interface Order {
  id: string; orderNo: string; customerId: string; customerName: string
  shipperName: string; shipperPhone: string; shipperAddress: string
  receiverName: string; receiverPhone: string; receiverAddress: string
  originName: string; destName: string
  planPickupTime: string; planDeliveryTime: string
  cargoName: string; cargoWeight: number; cargoVolume: number; cargoQuantity: number
  transportRequirement: string; status: string; remark: string
  cargoItems: CargoItem[]; createdAt: string
}

export interface CargoItem {
  id: string; cargoName: string; cargoType: string; quantity: number; weight: number; volume: number; unit: string
}

export function listOrders(params: Record<string, any>): Promise<PageResult<Order>> {
  return request.get('/orders', { params })
}

export function getOrder(id: string): Promise<Order> {
  return request.get(`/orders/${id}`)
}

export function createOrder(data: any): Promise<Order> {
  return request.post('/orders', data)
}

export function updateOrder(id: string, data: any): Promise<Order> {
  return request.put(`/orders/${id}`, data)
}

export function submitOrder(id: string): Promise<Order> {
  return request.post(`/orders/${id}/submit`)
}

export function cancelOrder(id: string): Promise<Order> {
  return request.post(`/orders/${id}/cancel`)
}

export interface DispatchPlan {
  id: string; planNo: string; status: string; remark: string
  details: DispatchDetail[]; createdAt: string
}

export interface DispatchDetail {
  id: string; orderId: string; carrierId: string; vehicleId: string; driverId: string; seq: number
}

export interface PendingOrder { id: string; orderNo: string; originName: string; destName: string }

export function listDispatchPlans(params: Record<string, any>): Promise<PageResult<DispatchPlan>> {
  return request.get('/dispatch/plans', { params })
}

export function getDispatchPlan(id: string): Promise<DispatchPlan> {
  return request.get(`/dispatch/plans/${id}`)
}

export function createDispatchPlan(orderIds: string[]): Promise<DispatchPlan> {
  return request.post('/dispatch/plans', { orderIds })
}

export function assignDispatch(id: string, data: any): Promise<void> {
  return request.post(`/dispatch/plans/${id}/assign`, data)
}

export function cancelDispatch(id: string): Promise<void> {
  return request.post(`/dispatch/plans/${id}/cancel`)
}

export function listPendingOrders(): Promise<PendingOrder[]> {
  return request.get('/dispatch/orders')
}
