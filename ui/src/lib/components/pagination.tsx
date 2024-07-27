// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useCallback, useMemo } from 'react'
import { ButtonGroup, Button, IconButton } from '@chakra-ui/react'
import {
  IconKeyboardArrowLeft,
  IconKeyboardArrowRight,
  IconKeyboardDoubleArrowLeft,
  IconKeyboardDoubleArrowRight,
  IconFirstPage,
  IconLastPage,
} from '@/lib/components/icons'

type PaginationProps = {
  totalPages: number
  page: number
  maxButtons?: number
  uiSize?: string
  onPageChange?: (page: number) => void
}

const Pagination = ({
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
        icon={<IconFirstPage />}
        aria-label="First"
        onClick={() => handlePageChange(firstPage)}
      />
      <IconButton
        variant="outline"
        size={uiSize}
        isDisabled={rewindPage < 1}
        icon={<IconKeyboardDoubleArrowLeft />}
        aria-label="Rewind"
        onClick={() => handlePageChange(rewindPage)}
      />
      <IconButton
        variant="outline"
        size={uiSize}
        isDisabled={page === 1}
        icon={<IconKeyboardArrowLeft />}
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
        icon={<IconKeyboardArrowRight />}
        aria-label="Next"
        onClick={() => handlePageChange(nextPage)}
      />
      <IconButton
        variant="outline"
        size={uiSize}
        isDisabled={fastForwardPage > lastPage}
        icon={<IconKeyboardDoubleArrowRight />}
        aria-label="Fast Forward"
        onClick={() => handlePageChange(fastForwardPage)}
      />
      <IconButton
        variant="outline"
        size={uiSize}
        isDisabled={page === lastPage}
        icon={<IconLastPage />}
        aria-label="Last"
        onClick={() => handlePageChange(lastPage)}
      />
    </ButtonGroup>
  )
}

export default Pagination
