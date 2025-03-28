// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
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
  Avatar,
  Badge,
  Button,
} from '@chakra-ui/react'
import {
  DataTable,
  IconAdd,
  PagePagination,
  RelativeDate,
  SectionError,
  SectionPlaceholder,
  SectionSpinner,
  Text,
  usePageMonitor,
  usePagePagination,
} from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import {
  OrganizationAPI,
  OrganizationSortOrder,
} from '@/client/api/organization'
import { errorToString } from '@/client/error'
import { swrConfig } from '@/client/options'
import { organizationPaginationStorage } from '@/infra/pagination'
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
    error: listError,
    isLoading: listIsLoading,
    mutate,
  } = OrganizationAPI.useList(
    { query, page, size, sortOrder: OrganizationSortOrder.Desc },
    swrConfig(),
  )
  const { hasPagination } = usePageMonitor({
    totalPages: list?.totalPages ?? 1,
    totalElements: list?.totalElements ?? 0,
    steps,
  })
  const listIsEmpty = list && !listError && list.totalElements === 0
  const listIsReady = list && !listError && list.totalElements > 0

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
        {listIsLoading ? <SectionSpinner /> : null}
        {listError ? <SectionError text={errorToString(listError)} /> : null}
        {listIsEmpty ? (
          <SectionPlaceholder
            text="There are no items."
            content={
              <Button
                as={Link}
                to="/new/organization"
                leftIcon={<IconAdd />}
                variant="solid"
              >
                New Organization
              </Button>
            }
          />
        ) : null}
        {listIsReady ? (
          <DataTable
            items={list.data}
            columns={[
              {
                title: 'Name',
                renderCell: (o) => (
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
                ),
              },
              {
                title: 'Permission',
                renderCell: (o) => <Badge>{o.permission}</Badge>,
              },
              {
                title: 'Date',
                renderCell: (o) => (
                  <RelativeDate date={new Date(o.createTime)} />
                ),
              },
            ]}
            pagination={
              hasPagination ? (
                <PagePagination
                  totalElements={list.totalElements}
                  totalPages={list.totalPages}
                  page={page}
                  size={size}
                  steps={steps}
                  setPage={setPage}
                  setSize={setSize}
                />
              ) : undefined
            }
          />
        ) : null}
      </div>
    </>
  )
}

export default OrganizationListPage
