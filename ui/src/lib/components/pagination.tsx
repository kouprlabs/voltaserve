import { useCallback, useMemo } from 'react'
import { ButtonGroup, Button, IconButton } from '@chakra-ui/react'
import {
  MdOutlineKeyboardArrowLeft,
  MdOutlineKeyboardArrowRight,
  MdOutlineKeyboardDoubleArrowLeft,
  MdOutlineKeyboardDoubleArrowRight,
  MdOutlineFirstPage,
  MdOutlineLastPage,
} from 'react-icons/md'

type PaginationProps = {
  totalPages: number
  page: number
  maxButtons?: number
  uiSize?: string
  onPageChange?: (page: number) => void
}

export const Pagination = ({
  totalPages,
  page,
  maxButtons: maxButtonsProp = 5,
  uiSize = 'md',
  onPageChange,
}: PaginationProps) => {
  const maxButtons = totalPages < maxButtonsProp ? totalPages : maxButtonsProp
  const pages = useMemo(() => {
    const end = Math.ceil(page / maxButtons) * maxButtons
    const start = end - maxButtons + 1
    return Array.from({ length: end - start + 1 }, (_, index) => start + index)
  }, [page, maxButtons])
  const firstPage = 1
  const lastPage = totalPages
  const fastForwardPage = pages[pages.length - 1] + 1
  const rewindPage = pages[0] - maxButtons
  const nextPage = page + 1
  const previousPage = page - 1

  const handlePageChange = useCallback(
    (value: number) => {
      if (value !== page) {
        onPageChange?.(value)
      }
    },
    [page, onPageChange],
  )

  return (
    <ButtonGroup>
      <IconButton
        variant="outline"
        size={uiSize}
        isDisabled={page === 1}
        icon={<MdOutlineFirstPage />}
        aria-label="First"
        onClick={() => handlePageChange(firstPage)}
      />
      <IconButton
        variant="outline"
        size={uiSize}
        isDisabled={rewindPage < 1}
        icon={<MdOutlineKeyboardDoubleArrowLeft />}
        aria-label="Rewind"
        onClick={() => handlePageChange(rewindPage)}
      />
      <IconButton
        variant="outline"
        size={uiSize}
        isDisabled={page === 1}
        icon={<MdOutlineKeyboardArrowLeft />}
        aria-label="Previous"
        onClick={() => handlePageChange(previousPage)}
      />
      {pages.map((index) => (
        <Button
          size={uiSize}
          key={index}
          isDisabled={index > totalPages}
          onClick={() => handlePageChange(index)}
          colorScheme={index === page ? 'blue' : undefined}
        >
          {index}
        </Button>
      ))}
      <IconButton
        variant="outline"
        size={uiSize}
        isDisabled={page === lastPage}
        icon={<MdOutlineKeyboardArrowRight />}
        aria-label="Next"
        onClick={() => handlePageChange(nextPage)}
      />
      <IconButton
        variant="outline"
        size={uiSize}
        isDisabled={fastForwardPage > lastPage}
        icon={<MdOutlineKeyboardDoubleArrowRight />}
        aria-label="Fast Forward"
        onClick={() => handlePageChange(fastForwardPage)}
      />
      <IconButton
        variant="outline"
        size={uiSize}
        isDisabled={page === lastPage}
        icon={<MdOutlineLastPage />}
        aria-label="Last"
        onClick={() => handlePageChange(lastPage)}
      />
    </ButtonGroup>
  )
}
