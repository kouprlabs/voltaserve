// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useEffect } from 'react'
import {
  Link,
  useLocation,
  useNavigate,
  useSearchParams,
} from 'react-router-dom'
import { Heading, Link as ChakraLink, Avatar, Badge } from '@chakra-ui/react'
import {
  DataTable,
  PagePagination,
  RelativeDate,
  SectionSpinner,
  Text,
  usePagePagination,
} from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import WorkspaceAPI, { SortOrder } from '@/client/api/workspace'
import { swrConfig } from '@/client/options'
import { CreateWorkspaceButton } from '@/components/app-bar/app-bar-buttons'
import { workspacePaginationStorage } from '@/infra/pagination'
import { decodeQuery } from '@/lib/helpers/query'
import { useAppDispatch } from '@/store/hook'
import { mutateUpdated } from '@/store/ui/workspaces'

const WorkspaceListPage = () => {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const location = useLocation()
  const [searchParams] = useSearchParams()
  const query = decodeQuery(searchParams.get('q') as string)
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigateFn: navigate,
    searchFn: () => location.search,
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
    mutate().then()
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
        {!list && !error ? <SectionSpinner /> : null}
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
          <DataTable
            items={list.data}
            columns={[
              {
                title: 'Name',
                renderCell: (w) => (
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
                      <Text noOfLines={1}>{w.name}</Text>
                    </ChakraLink>
                  </div>
                ),
              },
              {
                title: 'Organization',
                renderCell: (w) => (
                  <ChakraLink
                    as={Link}
                    to={`/organization/${w.organization.id}/member`}
                    className={cx('no-underline')}
                  >
                    <Text noOfLines={1}>{w.organization.name}</Text>
                  </ChakraLink>
                ),
              },
              {
                title: 'Permission',
                renderCell: (w) => <Badge>{w.permission}</Badge>,
              },
              {
                title: 'Date',
                renderCell: (w) => (
                  <RelativeDate date={new Date(w.createTime)} />
                ),
              },
            ]}
          />
        ) : null}
        {list ? (
          <div className={cx('self-end')}>
            <PagePagination
              totalElements={list.totalElements}
              totalPages={list.totalPages}
              page={page}
              size={size}
              steps={steps}
              setPage={setPage}
              setSize={setSize}
            />
          </div>
        ) : null}
      </div>
    </>
  )
}

export default WorkspaceListPage
