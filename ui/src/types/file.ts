import { File } from '@/client/api/file'

export enum ViewType {
  Grid = 'grid',
  List = 'list',
}

export type CommonItemProps = {
  file: File
  scale: number
  viewType: ViewType
  isPresentational?: boolean
  isLoading?: boolean
  isSelectionMode?: boolean
}
