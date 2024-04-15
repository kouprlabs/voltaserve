import { FileViewType } from '@/types/file'

export function computeScale(scale: number, viewType: FileViewType) {
  return viewType === 'list' ? scale * 0.5 : scale
}
