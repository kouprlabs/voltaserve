// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useCallback, useState } from 'react'
import {
  Link,
  useLocation,
  useNavigate,
  useSearchParams,
} from 'react-router-dom'
import { Avatar, Link as ChakraLink } from '@chakra-ui/react'
import { Heading } from '@chakra-ui/react'
import {
  DataTable,
  IconEdit,
  PagePagination,
  RelativeDate,
  SectionError,
  SectionSpinner,
  Text,
  usePagePagination,
} from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import ConsoleAPI, { ConsoleGroup } from '@/client/console/console'
import { swrConfig } from '@/client/options'
import ConsoleRenameModal from '@/components/console/console-rename-modal'
import { consoleGroupsPaginationStorage } from '@/infra/pagination'
import { decodeQuery } from '@/lib/helpers/query'

const ConsolePanelGroups = () => {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const location = useLocation()
  const query = decodeQuery(searchParams.get('q') as string)
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigateFn: navigate,
    searchFn: () => location.search,
    storage: consoleGroupsPaginationStorage(),
  })
  const [isConfirmRenameOpen, setIsConfirmRenameOpen] = useState(false)
  const [currentName, setCurrentName] = useState<string>('')
  const [groupId, setGroupId] = useState<string>()
  const {
    data: list,
    error,
    mutate,
  } = ConsoleAPI.useListOrSearchObject<ConsoleGroup>(
    'group',
    { page, size, query },
    swrConfig(),
  )

  const renameRequest = useCallback(
    async (name: string) => {
      if (groupId) {
        await ConsoleAPI.renameObject({ id: groupId, name }, 'group')
        await mutate()
      }
    },
    [groupId],
  )

  return (
    <>
      <Helmet>
        <title>Groups</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>Groups</Heading>
        {!list && error ? <SectionError text="Failed to load groups." /> : null}
        {!list && !error ? <SectionSpinner /> : null}
        {list && list.totalElements > 0 ? (
          <>
            <DataTable
              items={list.data}
              columns={[
                {
                  title: 'Name',
                  renderCell: (group) => (
                    <div
                      className={cx(
                        'flex',
                        'flex-row',
                        'items-center',
                        'gap-1.5',
                      )}
                    >
                      <Avatar
                        name={group.name}
                        size="sm"
                        className={cx('w-[40px]', 'h-[40px]')}
                      />
                      <Text noOfLines={1}>{group.name}</Text>
                    </div>
                  ),
                },
                {
                  title: 'Organization',
                  renderCell: (group) => (
                    <ChakraLink
                      as={Link}
                      to={`/console/organizations/${group.organization.id}`}
                      className={cx('no-underline')}
                    >
                      <Text noOfLines={1}>{group.organization.name}</Text>
                    </ChakraLink>
                  ),
                },
                {
                  title: 'Created',
                  renderCell: (group) => (
                    <RelativeDate date={new Date(group.createTime)} />
                  ),
                },
                {
                  title: 'Updated',
                  renderCell: (group) => (
                    <RelativeDate date={new Date(group.updateTime)} />
                  ),
                },
              ]}
              actions={[
                {
                  label: 'Edit Name',
                  icon: <IconEdit />,
                  onClick: async (group) => {
                    setCurrentName(group.name)
                    setGroupId(group.id)
                    setIsConfirmRenameOpen(true)
                  },
                },
              ]}
              pagination={
                list.totalPages > 1 ? (
                  <PagePagination
                    totalElements={list.totalElements}
                    totalPages={Math.ceil(list.totalElements / size)}
                    page={page}
                    size={size}
                    steps={steps}
                    setPage={setPage}
                    setSize={setSize}
                  />
                ) : undefined
              }
            />
          </>
        ) : null}
      </div>
      <ConsoleRenameModal
        currentName={currentName}
        isOpen={isConfirmRenameOpen}
        onClose={() => setIsConfirmRenameOpen(false)}
        onRequest={renameRequest}
      />
    </>
  )
}

export default ConsolePanelGroups
