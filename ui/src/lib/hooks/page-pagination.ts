import { useCallback, useEffect, useMemo, useState } from 'react'

export type NavigateArgs = {
  search: string
}

export type NavigateFunction = (args: NavigateArgs) => void

export type LocationObject = {
  search: string
}

export type StorageOptions = {
  enabled?: boolean
  prefix?: string
  namespace?: string
}

export type UsePagePaginationOptions = {
  navigate: NavigateFunction
  location: LocationObject
  storage?: StorageOptions
  steps?: number[]
}

const usePagePagination = ({
  navigate,
  location,
  storage = {
    enabled: false,
    prefix: 'app',
    namespace: 'main',
  },
  steps = [5, 10, 20, 40, 80, 100],
}: UsePagePaginationOptions) => {
  const queryParams = useMemo(
    () => new URLSearchParams(location.search),
    [location.search],
  )
  const page = Number(queryParams.get('page')) || 1
  const storageSizeKey = useMemo(
    () => `${storage.prefix}_${storage.namespace}_pagination_size`,
    [storage],
  )
  const [size, setSize] = useState(
    localStorage.getItem(storageSizeKey) && storage.enabled
      ? parseInt(localStorage.getItem(storageSizeKey) as string)
      : steps[0],
  )

  const _setSize = useCallback(
    (size: number) => {
      setSize(size)
      if (size && storage.enabled) {
        localStorage.setItem(storageSizeKey, JSON.stringify(size))
      }
    },
    [storageSizeKey, storage],
  )

  useEffect(() => {
    if (!queryParams.has('page')) {
      queryParams.set('page', '1')
      navigate({ search: `?${queryParams.toString()}` })
    }
  }, [queryParams, navigate])

  const setPage = useCallback(
    (page: number) => {
      queryParams.set('page', String(page))
      navigate({ search: `?${queryParams.toString()}` })
    },
    [queryParams, navigate],
  )

  return { page, size, steps, setPage, setSize: _setSize }
}

export default usePagePagination
