import { request } from '../request'

export interface DashboardData {
  pendingDispatch: number
  inTransit: number
  pendingException: number
  pendingSettlement: number
}

export interface OrderStats {
  totalOrders: number
  byStatus: StatusCount[]
}

export interface StatusCount {
  status: string
  count: number
}

export interface TransportEfficiency {
  totalTasks: number
  completedTasks: number
  avgTransitHours: number
}

export function getDashboard(): Promise<DashboardData> {
  return request.get('/reports/dashboard')
}

export function getOrderStats(): Promise<OrderStats> {
  return request.get('/reports/order-stats')
}

export function getTransportEfficiency(): Promise<TransportEfficiency> {
  return request.get('/reports/transport-efficiency')
}
