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
import { PagePagination, SectionSpinner, usePagePagination } from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import OrganizationAPI, { SortOrder } from '@/client/api/organization'
import { swrConfig } from '@/client/options'
import { CreateOrganizationButton } from '@/components/app-bar/app-bar-buttons'
import { organizationPaginationStorage } from '@/infra/pagination'
import prettyDate from '@/lib/helpers/pretty-date'
import { decodeQuery } from '@/lib/helpers/query'
import { useAppDispatch } from '@/store/hook'
import { mutateUpdated } from '@/store/ui/organizations'

const OrganizationListPage = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const dispatch = useAppDispatch()
  const [searchParams] = useSearchParams()
  const query = decodeQuery(searchParams.get('q') as string)
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigateFn: navigate,
    searchFn: () => location.search,
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
                        <Text noOfLines={1}>{o.name}</Text>
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

export default OrganizationListPage
