import { request } from '../request'

export interface FileObject {
  id: string; bucket: string; objectKey: string; originalName: string
  contentType: string; size: number; bizType: string; bizId: string; createdAt: string
}

export function uploadFile(file: File, bizType: string, bizId: string): Promise<FileObject> {
  const form = new FormData()
  form.append('file', file)
  form.append('bizType', bizType)
  form.append('bizId', bizId)
  return request.post('/files/upload', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export function listFiles(bizType: string, bizId: string): Promise<FileObject[]> {
  return request.get('/files', { params: { bizType, bizId } })
}

export function deleteFile(id: string): Promise<void> {
  return request.delete(`/files/${id}`)
}

export function getDownloadUrl(id: string): string {
  return `${import.meta.env.VITE_API_BASE_URL || '/api/v1'}/files/${id}/download`
}
