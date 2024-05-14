import { File } from '@/client/api/file'

export enum FileViewType {
  Grid = 'grid',
  List = 'list',
}

export type FileCommonProps = {
  file: File
  scale: number
  viewType: FileViewType
  isPresentational?: boolean
  isDragging?: boolean
  isLoading?: boolean
  isSelectionMode?: boolean
}
