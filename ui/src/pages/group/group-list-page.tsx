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
  Text,
} from '@chakra-ui/react'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import GroupAPI, { SortOrder } from '@/client/api/group'
import { swrConfig } from '@/client/options'
import { CreateGroupButton } from '@/components/app-bar/app-bar-buttons'
import { groupPaginationStorage } from '@/infra/pagination'
import PagePagination from '@/lib/components/page-pagination'
import SectionSpinner from '@/lib/components/section-spinner'
import prettyDate from '@/lib/helpers/pretty-date'
import { decodeQuery } from '@/lib/helpers/query'
import usePagePagination from '@/lib/hooks/page-pagination'
import { useAppDispatch } from '@/store/hook'
import { mutateUpdated } from '@/store/ui/groups'

const GroupListPage = () => {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
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

  useEffect(() => {
    if (mutate) {
      dispatch(mutateUpdated(mutate))
    }
  }, [mutate, dispatch])

  return (
    <>
      <Helmet>
        <title>Groups</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>Groups</Heading>
        {error ? (
          <div
            className={cx(
              'flex',
              'items-center',
              'justify-center',
              'h-[300px]',
            )}
          >
            <span>Failed to load groups.</span>
          </div>
        ) : null}
        {!list && !error ? <SectionSpinner /> : null}
        {list && list.data.length === 0 ? (
          <div
            className={cx(
              'flex',
              'items-center',
              'justify-center',
              'h-[300px]',
            )}
          >
            <div className={cx('flex', 'flex-col', 'gap-1.5')}>
              <span>There are no groups.</span>
              <CreateGroupButton />
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
              {list.data.map((g) => (
                <Tr key={g.id}>
                  <Td>
                    <div
                      className={cx(
                        'flex',
                        'flex-row',
                        'items-center',
                        'gap-1.5',
                      )}
                    >
                      <Avatar
                        name={g.name}
                        size="sm"
                        className={cx('w-[40px]', 'h-[40px]')}
                      />
                      <ChakraLink
                        as={Link}
                        to={`/group/${g.id}/member`}
                        className={cx('no-underline')}
                      >
                        <Text noOfLines={1}>{g.name}</Text>
                      </ChakraLink>
                    </div>
                  </Td>
                  <Td>
                    <ChakraLink
                      as={Link}
                      to={`/organization/${g.organization.id}/member`}
                      className={cx('no-underline')}
                    >
                      <Text noOfLines={1}>{g.organization.name}</Text>
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

export default GroupListPage
