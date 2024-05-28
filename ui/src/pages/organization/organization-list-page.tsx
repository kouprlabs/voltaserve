import { useEffect } from 'react'
import {
  Link,
  useLocation,
  useNavigate,
  useSearchParams,
} from 'react-router-dom'
import {
  Heading,
  Link as ChakraLink,
  Table,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
  Avatar,
  Badge,
} from '@chakra-ui/react'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import OrganizationAPI, { SortOrder } from '@/client/api/organization'
import { swrConfig } from '@/client/options'
import { CreateOrganizationButton } from '@/components/top-bar/top-bar-buttons'
import prettyDate from '@/helpers/pretty-date'
import { decodeQuery } from '@/helpers/query'
import { organizationPaginationStorage } from '@/infra/pagination'
import { SectionSpinner, PagePagination, usePagePagination } from '@/lib'
import { useAppDispatch } from '@/store/hook'
import { mutateUpdated } from '@/store/ui/organizations'

const OrganizationListPage = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const dispatch = useAppDispatch()
  const [searchParams] = useSearchParams()
  const query = decodeQuery(searchParams.get('q') as string)
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigate,
    location,
    storage: organizationPaginationStorage(),
  })
  const {
    data: list,
    error,
    mutate,
  } = OrganizationAPI.useList(
    { query, page, size, sortOrder: SortOrder.Desc },
    swrConfig(),
  )

  useEffect(() => {
    mutate()
  }, [query, page, size, mutate])

  useEffect(() => {
    if (mutate) {
      dispatch(mutateUpdated(mutate))
    }
  }, [mutate, dispatch])

  return (
    <>
      <Helmet>
        <title>Organizations</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>Organizations</Heading>
        {!list && error ? (
          <div
            className={cx(
              'flex',
              'items-center',
              'justify-center',
              'h-[300px]',
            )}
          >
            <span>Failed to load organizations.</span>
          </div>
        ) : null}
        {!list && !error && <SectionSpinner />}
        {list && list.data.length === 0 ? (
          <div
            className={cx(
              'flex',
              'items-center',
              'justify-center',
              'h-[300px]',
            )}
          >
            <div className={cx('flex', 'flex-col', 'gap-1.5', 'items-center')}>
              <span>There are no organizations.</span>
              <CreateOrganizationButton />
            </div>
          </div>
        ) : null}
        {list && list.data.length > 0 ? (
          <Table variant="simple">
            <Thead>
              <Tr>
                <Th>Name</Th>
                <Th>Permission</Th>
                <Th>Date</Th>
              </Tr>
            </Thead>
            <Tbody>
              {list.data.map((o) => (
                <Tr key={o.id}>
                  <Td>
                    <div
                      className={cx(
                        'flex',
                        'flex-row',
                        'gap-1.5',
                        'items-center',
                      )}
                    >
                      <Avatar
                        name={o.name}
                        size="sm"
                        className={cx('w-[40px]', 'h-[40px]')}
                      />
                      <ChakraLink
                        as={Link}
                        to={`/organization/${o.id}/member`}
                        className={cx('no-underline')}
                      >
                        <span>{o.name}</span>
                      </ChakraLink>
                    </div>
                  </Td>
                  <Td>
                    <Badge>{o.permission}</Badge>
                  </Td>
                  <Td>{prettyDate(o.createTime)}</Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
        ) : null}
        {list ? (
          <PagePagination
            style={{ alignSelf: 'end' }}
            totalElements={list.totalElements}
            totalPages={list.totalPages}
            page={page}
            size={size}
            steps={steps}
            setPage={setPage}
            setSize={setSize}
          />
        ) : null}
      </div>
    </>
  )
}

export default OrganizationListPage
