import { request } from '@/api/request'

export interface HealthStatus {
  status: string
  service: string
}

export function getHealth(): Promise<HealthStatus> {
  return request.get('/health') as unknown as Promise<HealthStatus>
}
