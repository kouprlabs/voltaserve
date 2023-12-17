import { Select } from '@chakra-ui/react'
import Pagination from '@/components/common/pagination'

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
