import { useCallback, useEffect, useMemo, useState } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { Select } from '@chakra-ui/react'
import Pagination from '@/components/common/pagination'

type PagePaginationHookOptions = {
  disableLocalStorage?: boolean
  localStoragePrefix?: string
  localStorageNamespace?: string
}

export const usePagePagination = ({
  disableLocalStorage = false,
  localStoragePrefix = 'app',
  localStorageNamespace = 'main',
}: PagePaginationHookOptions) => {
  const navigate = useNavigate()
  const location = useLocation()
  const queryParams = useMemo(
    () => new URLSearchParams(location.search),
    [location.search]
  )
  const page = Number(queryParams.get('page')) || 1
  const localStorageSizeKey = useMemo(
    () => `${localStoragePrefix}_${localStorageNamespace}_pagination_size`,
    [localStoragePrefix, localStorageNamespace]
  )
  const [size, setSize] = useState(
    localStorage.getItem(localStorageSizeKey) && !disableLocalStorage
      ? parseInt(localStorage.getItem(localStorageSizeKey) as string)
      : 5
  )

  useEffect(() => {
    if (size && !disableLocalStorage) {
      localStorage.setItem(localStorageSizeKey, JSON.stringify(size))
    }
  }, [size, localStorageSizeKey, disableLocalStorage])

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
  paginationSize?: string
  onPageChange: (page: number) => void
  onSizeChange: (size: number) => void
}

const PagePagination = ({
  totalPages,
  page,
  size,
  paginationSize = 'md',
  onPageChange,
  onSizeChange,
}: PagePaginationProps) => {
  return (
    <>
      {totalPages > 1 ? (
        <Pagination
          size={paginationSize}
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
