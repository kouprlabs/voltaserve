import { ListOptions, ListQueryParams } from '@/client/api/types/queries'
import { encodeQuery } from '@/lib/helpers/query'

export const paramsFromListOptions = (
  options?: ListOptions,
): URLSearchParams => {
  const params: ListQueryParams = {}
  if (options?.query) {
    params.query = encodeQuery(JSON.stringify(options.query))
  }
  if (options?.page) {
    params.page = options.page.toString()
  }
  if (options?.size) {
    params.size = options.size.toString()
  }
  if (options?.sortBy) {
    params.sort_by = options.sortBy.toString()
  }
  if (options?.sortOrder) {
    params.sort_order = options.sortOrder.toString()
  }
  if (options?.organizationId) {
    params.organization_id = options.organizationId.toString()
  }
  if (options?.type) {
    params.type = options.type
  }
  if (options?.groupId) {
    params.group_id = options.groupId.toString()
  }
  if (options?.excludeGroupMembers) {
    params.exclude_group_members = options.excludeGroupMembers.toString()
  }
  return new URLSearchParams(params)
}
