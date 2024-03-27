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
import WorkspaceAPI, { SortOrder } from '@/client/api/workspace'
import { swrConfig } from '@/client/options'
import { CreateWorkspaceButton } from '@/components/top-bar/top-bar-buttons'
import prettyDate from '@/helpers/pretty-date'
import { decodeQuery } from '@/helpers/query'
import { workspacePaginationStorage } from '@/infra/pagination'

const WorkspaceListPage = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const [searchParams] = useSearchParams()
  const query = decodeQuery(searchParams.get('q') as string)
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigate,
    location,
    storage: workspacePaginationStorage(),
  })
  const {
    data: list,
    error,
    mutate,
  } = WorkspaceAPI.useList(
    { query, page, size, sortOrder: SortOrder.Desc },
    swrConfig(),
  )

  useEffect(() => {
    mutate()
  }, [query, page, size, mutate])

  return (
    <>
      <Helmet>
        <title>Workspaces</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading fontSize={variables.headingFontSize} pl={variables.spacingMd}>
          Workspaces
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
            <Text>Failed to load workspaces.</Text>
          </div>
        )}
        {!list && !error && <SectionSpinner />}
        {list && list.data.length === 0 && !error ? (
          <div
            className={cx(
              'flex',
              'items-center',
              'justify-center',
              'h-[300px]',
            )}
          >
            <div className={cx('flex', 'flex-col', 'gap-1.5', 'items-center')}>
              <Text>There are no workspaces.</Text>
              <CreateWorkspaceButton />
            </div>
          </div>
        ) : null}
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
              {list.data.map((w) => (
                <Tr key={w.id}>
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
                        name={w.name}
                        size="sm"
                        width="40px"
                        height="40px"
                      />
                      <ChakraLink
                        as={Link}
                        to={`/workspace/${w.id}/file/${w.rootId}`}
                        textDecoration="none"
                      >
                        <Text>{w.name}</Text>
                      </ChakraLink>
                    </div>
                  </Td>
                  <Td>
                    <ChakraLink
                      as={Link}
                      to={`/organization/${w.organization.id}/member`}
                      textDecoration="none"
                    >
                      {w.organization.name}
                    </ChakraLink>
                  </Td>
                  <Td>
                    <Badge>{w.permission}</Badge>
                  </Td>
                  <Td>{prettyDate(w.createTime)}</Td>
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

export default WorkspaceListPage
