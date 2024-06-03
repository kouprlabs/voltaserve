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
import WorkspaceAPI, { SortOrder } from '@/client/api/workspace'
import { swrConfig } from '@/client/options'
import { CreateWorkspaceButton } from '@/components/app-bar/app-bar-buttons'
import prettyDate from '@/helpers/pretty-date'
import { decodeQuery } from '@/helpers/query'
import { workspacePaginationStorage } from '@/infra/pagination'
import PagePagination from '@/lib/components/page-pagination'
import SectionSpinner from '@/lib/components/section-spinner'
import usePagePagination from '@/lib/hooks/page-pagination'
import { useAppDispatch } from '@/store/hook'
import { mutateUpdated } from '@/store/ui/workspaces'

const WorkspaceListPage = () => {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
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

  useEffect(() => {
    if (mutate) {
      dispatch(mutateUpdated(mutate))
    }
  }, [mutate, dispatch])

  return (
    <>
      <Helmet>
        <title>Workspaces</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>Workspaces</Heading>
        {!list && error ? (
          <div
            className={cx(
              'flex',
              'items-center',
              'justify-center',
              'h-[300px]',
            )}
          >
            <span>Failed to load workspaces.</span>
          </div>
        ) : null}
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
              <span>There are no workspaces.</span>
              <CreateWorkspaceButton />
            </div>
          </div>
        ) : null}
        {list && list.data.length > 0 ? (
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
                        className={cx('w-[40px]', 'h-[40px]')}
                      />
                      <ChakraLink
                        as={Link}
                        to={`/workspace/${w.id}/file/${w.rootId}`}
                        className={cx('no-underline')}
                      >
                        <span>{w.name}</span>
                      </ChakraLink>
                    </div>
                  </Td>
                  <Td>
                    <ChakraLink
                      as={Link}
                      to={`/organization/${w.organization.id}/member`}
                      className={cx('no-underline')}
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

export default WorkspaceListPage
