import { Query } from '@/client/api/file'

export type ListQueryParams = {
  file_id?: string
  page?: string
  size?: string
  sort_by?: string
  sort_order?: string
  query?: string
  type?: string
  organization_id?: string
  group_id?: string
  exclude_group_members?: string
}

export type ListOptions = {
  fileId?: string
  type?: FileType
  query?: string | Query
  organizationId?: string
  groupId?: string
  excludeGroupMembers?: boolean
  size?: number
  page?: number
  sortBy?: SortBy
  sortOrder?: SortOrder
}

export enum SortBy {
  Name = 'name',
  DateCreated = 'date_created',
  DateModified = 'date_modified',
  Frequency = 'frequency',
  Kind = 'kind',
  Size = 'size',
  Email = 'email',
  FullName = 'full_name',
  Version = 'version',
}

export enum SortOrder {
  Asc = 'asc',
  Desc = 'desc',
}

export enum FileType {
  File = 'file',
  Folder = 'folder',
}
