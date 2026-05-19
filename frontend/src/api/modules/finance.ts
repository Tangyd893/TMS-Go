import { request, type PageResult } from '../request'

export interface Receivable {
  id: string; receivableNo: string; orderId: string; taskId: string
  customerId: string; customerName: string
  feeItemId: string; feeItemName: string
  amount: number; hasTax: boolean; taxRate: number; taxAmount: number; totalAmount: number
  status: string; remark: string; createdAt: string
}

export interface Payable {
  id: string; payableNo: string; orderId: string; taskId: string
  carrierId: string; carrierName: string
  feeItemId: string; feeItemName: string
  amount: number; hasTax: boolean; taxRate: number; taxAmount: number; totalAmount: number
  status: string; remark: string; createdAt: string
}

export interface FeeItem {
  id: string; name: string; code: string
  unitPrice: number; unit: string; remark: string
  createdAt: string
}

export interface Statement {
  id: string; statementNo: string; partnerId: string; partnerName: string
  partnerType: string; statementPeriodStart: string; statementPeriodEnd: string
  totalReceivable: number; totalPayable: number; status: string
  confirmedAt: string; confirmedBy: string; remark: string; createdAt: string
}

export interface Settlement {
  id: string; settlementNo: string; statementId: string; partnerId: string
  partnerName: string; partnerType: string
  settlementAmount: number; settlementMethod: string; status: string
  settledAt: string; remark: string; createdAt: string
}

export function listReceivables(params: Record<string, any>): Promise<PageResult<Receivable>> {
  return request.get('/finance/receivables', { params })
}

export function listPayables(params: Record<string, any>): Promise<PageResult<Payable>> {
  return request.get('/finance/payables', { params })
}

export function listFeeItems(params: Record<string, any>): Promise<PageResult<FeeItem>> {
  return request.get('/finance/fee-items', { params })
}

export function createFeeItem(data: any): Promise<FeeItem> {
  return request.post('/finance/fee-items', data)
}

export function generateFees(data: { orderId: string; customerId: string; carrierId: string; customerName: string; carrierName: string; amount: number }): Promise<void> {
  return request.post('/finance/fees/generate', data)
}

export function listStatements(params: Record<string, any>): Promise<PageResult<Statement>> {
  return request.get('/finance/statements', { params })
}

export function createStatement(data: { partnerId: string; partnerName: string; partnerType: string; periodStart: string; periodEnd: string }): Promise<Statement> {
  return request.post('/finance/statements', data)
}

export function confirmStatement(id: string): Promise<void> {
  return request.post(`/finance/statements/${id}/confirm`)
}

export function listSettlements(params: Record<string, any>): Promise<PageResult<Settlement>> {
  return request.get('/finance/settlements', { params })
}

export function createSettlement(statementId: string): Promise<Settlement> {
  return request.post('/finance/settlements', { statementId })
}

export function completeSettlement(id: string): Promise<void> {
  return request.post(`/finance/settlements/${id}/complete`)
}
