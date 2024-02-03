import { useNavigate } from 'react-router-dom'
import { usePagePagination } from '@koupr/ui'
import { filePaginationSteps, filesPaginationStorage } from '@/infra/pagination'
import { useAppSelector } from '@/store/hook'

export default function useFileListSearchParams() {
  const navigate = useNavigate()
  const sortBy = useAppSelector((state) => state.ui.files.sortBy)
  const sortOrder = useAppSelector((state) => state.ui.files.sortOrder)
  const { page, size } = usePagePagination({
    navigate,
    location,
    storage: filesPaginationStorage(),
    steps: filePaginationSteps(),
  })

  const params: any = {}
  if (page) {
    params.page = page.toString()
  }
  if (size) {
    params.size = size.toString()
  }
  if (sortBy) {
    params.sort_by = sortBy.toString()
  }
  if (sortOrder) {
    params.sort_order = sortOrder.toString()
  }
  return new URLSearchParams(params)
}
