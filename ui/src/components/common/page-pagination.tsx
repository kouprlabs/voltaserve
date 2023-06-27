import { useCallback, useEffect, useState } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { Select } from '@chakra-ui/react'
import Pagination from '@/components/common/pagination'

export const usePagePagination = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const queryParams = new URLSearchParams(location.search)
  const page = Number(queryParams.get('page')) || 1
  const [size, setSize] = useState(5)

  useEffect(() => {
    if (!queryParams.has('page')) {
      queryParams.set('page', '1')
      navigate({ search: `?${queryParams.toString()}` })
    }
  }, [queryParams, navigate])

  const handlePageChange = useCallback(
    (page: number) => {
      queryParams.set('page', String(page))
      navigate({ search: `?${queryParams.toString()}` })
    },
    [queryParams, navigate]
  )

  return { page, size, onPageChange: handlePageChange, onSizeChange: setSize }
}

type PagePaginationProps = {
  totalPages: number
  page: number
  size: number
  onPageChange: (page: number) => void
  onSizeChange: (size: number) => void
}

const PagePagination = ({
  totalPages,
  page,
  size,
  onPageChange,
  onSizeChange,
}: PagePaginationProps) => {
  return (
    <>
      {totalPages > 1 ? (
        <Pagination
          page={page}
          totalPages={totalPages}
          onPageChange={onPageChange}
        />
      ) : null}
      <Select
        defaultValue={size}
        onChange={(event) => onSizeChange(parseInt(event.target.value))}
      >
        <option value="5">5 items</option>
        <option value="10">10 items</option>
        <option value="20">20 items</option>
        <option value="40">40 items</option>
        <option value="80">80 items</option>
        <option value="100">100 items</option>
      </Select>
    </>
  )
}

export default PagePagination
