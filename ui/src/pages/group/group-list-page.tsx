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
import classNames from 'classnames'
import { Helmet } from 'react-helmet-async'
import GroupAPI, { SortOrder } from '@/client/api/group'
import { swrConfig } from '@/client/options'
import { CreateGroupButton } from '@/components/top-bar/top-bar-buttons'
import prettyDate from '@/helpers/pretty-date'
import { decodeQuery } from '@/helpers/query'
import { groupPaginationStorage } from '@/infra/pagination'

const GroupListPage = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const [searchParams] = useSearchParams()
  const query = decodeQuery(searchParams.get('q') as string)
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigate,
    location,
    storage: groupPaginationStorage(),
  })
  const {
    data: list,
    error,
    mutate,
  } = GroupAPI.useList(
    { query, page, size, sortOrder: SortOrder.Desc },
    swrConfig(),
  )

  useEffect(() => {
    mutate()
  }, [query, page, size, mutate])

  return (
    <>
      <Helmet>
        <title>Groups</title>
      </Helmet>
      <div className={classNames('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading fontSize={variables.headingFontSize} pl={variables.spacingMd}>
          Groups
        </Heading>
        {error && (
          <div
            className={classNames(
              'flex',
              'items-center',
              'justify-center',
              'h-[300px]',
            )}
          >
            <Text>Failed to load groups.</Text>
          </div>
        )}
        {!list && !error && <SectionSpinner />}
        {list && list.data.length === 0 && (
          <div
            className={classNames(
              'flex',
              'items-center',
              'justify-center',
              'h-[300px]',
            )}
          >
            <div className={classNames('flex', 'flex-col', 'gap-1.5')}>
              <Text>There are no groups.</Text>
              <CreateGroupButton />
            </div>
          </div>
        )}
        {list && list.data.length > 0 && (
          <Table variant="simple">
            <Thead>
              <Tr>
                <Th>Name</Th>
                <Th>Organization</Th>
                <Th>Permission</Th>
                <Th>Date</Th>
              </Tr>
            </Thead>
            <Tbody>
              {list.data.map((g) => (
                <Tr key={g.id}>
                  <Td>
                    <div
                      className={classNames(
                        'flex',
                        'flex-row',
                        'items-center',
                        'gap-1.5',
                      )}
                    >
                      <Avatar
                        name={g.name}
                        size="sm"
                        width="40px"
                        height="40px"
                      />
                      <ChakraLink
                        as={Link}
                        to={`/group/${g.id}/member`}
                        textDecoration="none"
                      >
                        {g.name}
                      </ChakraLink>
                    </div>
                  </Td>
                  <Td>
                    <ChakraLink
                      as={Link}
                      to={`/organization/${g.organization.id}/member`}
                      textDecoration="none"
                    >
                      {g.organization.name}
                    </ChakraLink>
                  </Td>
                  <Td>
                    <Badge>{g.permission}</Badge>
                  </Td>
                  <Td>{prettyDate(g.createTime)}</Td>
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

export default GroupListPage
