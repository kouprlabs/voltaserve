import { useCallback, useEffect, useMemo, useState } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'

type PagePaginationHookOptions = {
  disableLocalStorage?: boolean
  localStoragePrefix?: string
  localStorageNamespace?: string
}

const usePagePagination = ({
  disableLocalStorage = false,
  localStoragePrefix = 'app',
  localStorageNamespace = 'main',
}: PagePaginationHookOptions) => {
  const navigate = useNavigate()
  const location = useLocation()
  const queryParams = useMemo(
    () => new URLSearchParams(location.search),
    [location.search],
  )
  const page = Number(queryParams.get('page')) || 1
  const localStorageSizeKey = useMemo(
    () => `${localStoragePrefix}_${localStorageNamespace}_pagination_size`,
    [localStoragePrefix, localStorageNamespace],
  )
  const [size, setSize] = useState(
    localStorage.getItem(localStorageSizeKey) && !disableLocalStorage
      ? parseInt(localStorage.getItem(localStorageSizeKey) as string)
      : 5,
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
    [queryParams, navigate],
  )

  return { page, size, onPageChange: handlePageChange, onSizeChange: setSize }
}

export default usePagePagination
