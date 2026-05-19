import { request } from '../request'

export interface DictType {
  id: string
  code: string
  name: string
  status: string
  remark: string
  items: DictItem[]
}

export interface DictItem {
  id: string
  typeCode: string
  itemCode: string
  itemName: string
  sortNo: number
  status: string
}

export function listDictTypes(): Promise<DictType[]> {
  return request.get('/dict/types') as Promise<DictType[]>
}

export function getDictType(code: string): Promise<DictType> {
  return request.get(`/dict/types/${code}`) as Promise<DictType>
}

export function createDictType(data: any): Promise<DictType> {
  return request.post('/dict/types', data) as Promise<DictType>
}

export function updateDictType(id: string, data: any): Promise<DictType> {
  return request.put(`/dict/types/${id}`, data) as Promise<DictType>
}

export function deleteDictType(id: string): Promise<void> {
  return request.delete(`/dict/types/${id}`)
}

export function listDictItems(typeCode: string): Promise<DictItem[]> {
  return request.get('/dict/items', { params: { typeCode } }) as Promise<DictItem[]>
}

export function createDictItem(data: any): Promise<DictItem> {
  return request.post('/dict/items', data) as Promise<DictItem>
}

export function updateDictItem(id: string, data: any): Promise<DictItem> {
  return request.put(`/dict/items/${id}`, data) as Promise<DictItem>
}

export function deleteDictItem(id: string): Promise<void> {
  return request.delete(`/dict/items/${id}`)
}
