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
  Text,
  Avatar,
  Badge,
} from '@chakra-ui/react'
import {
  SectionSpinner,
  PagePagination,
  variables,
  usePagePagination,
} from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import OrganizationAPI, { SortOrder } from '@/client/api/organization'
import { swrConfig } from '@/client/options'
import { CreateOrganizationButton } from '@/components/top-bar/top-bar-buttons'
import prettyDate from '@/helpers/pretty-date'
import { decodeQuery } from '@/helpers/query'
import { organizationPaginationStorage } from '@/infra/pagination'

const OrganizationListPage = () => {
  const navigate = useNavigate()
  const location = useLocation()
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

  return (
    <>
      <Helmet>
        <title>Organizations</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading fontSize={variables.headingFontSize} pl={variables.spacingMd}>
          Organizations
        </Heading>
        {!list && error && (
          <div
            className={cx(
              'flex',
              'items-center',
              'justify-center',
              'h-[300px]',
            )}
          >
            <Text>Failed to load organizations.</Text>
          </div>
        )}
        {!list && !error && <SectionSpinner />}
        {list && list.data.length === 0 && (
          <div
            className={cx(
              'flex',
              'items-center',
              'justify-center',
              'h-[300px]',
            )}
          >
            <div className={cx('flex', 'flex-col', 'gap-1.5', 'items-center')}>
              <Text>There are no organizations.</Text>
              <CreateOrganizationButton />
            </div>
          </div>
        )}
        {list && list.data.length > 0 && (
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
                        width="40px"
                        height="40px"
                      />
                      <ChakraLink
                        as={Link}
                        to={`/organization/${o.id}/member`}
                        className={cx('no-underline')}
                      >
                        <Text>{o.name}</Text>
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
        )}
        {list && (
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
        )}
      </div>
    </>
  )
}

export default OrganizationListPage
