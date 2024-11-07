// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useEffect, useState } from 'react'
import {
  useLocation,
  useNavigate,
  useSearchParams,
  Link,
} from 'react-router-dom'
import { Avatar, Heading, Link as ChakraLink } from '@chakra-ui/react'
import {
  DataTable,
  IconEdit,
  PagePagination,
  RelativeDate,
  SectionSpinner,
  Text,
  usePagePagination,
} from '@koupr/ui'
import * as Yup from 'yup'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import ConsoleAPI, {
  OrganizationManagementList,
} from '@/client/console/console'
import ConsoleRenameModal from '@/components/console/console-rename-modal'
import { consoleOrganizationsPaginationStorage } from '@/infra/pagination'
import { decodeQuery } from '@/lib/helpers/query'

const ConsolePanelOrganizations = () => {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const location = useLocation()
  const query = decodeQuery(searchParams.get('q') as string)
  const [list, setList] = useState<OrganizationManagementList>()
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigateFn: navigate,
    searchFn: () => location.search,
    storage: consoleOrganizationsPaginationStorage(),
  })
  const [confirmRenameWindowOpen, setConfirmRenameWindowOpen] = useState(false)
  const [isSubmitting, setSubmitting] = useState(false)
  const [currentName, setCurrentName] = useState<string>('')
  const [organizationId, setOrganizationId] = useState<string>()
  const formSchema = Yup.object().shape({
    name: Yup.string().required('Name is required').max(255),
  })

  const renameOrganization = async (
    id: string | null,
    currentName: string | null,
    newName: string | null,
    confirm: boolean = false,
  ) => {
    if (confirm && organizationId !== undefined && newName !== null) {
      try {
        setSubmitting(true)
        await ConsoleAPI.renameObject(
          { id: organizationId, name: newName },
          'organization',
        )
      } finally {
        closeConfirmationWindow()
      }
    } else if (id !== null && currentName !== null && currentName !== '') {
      setConfirmRenameWindowOpen(true)
      setCurrentName(currentName)
      setOrganizationId(id)
    }
  }

  const closeConfirmationWindow = () => {
    setConfirmRenameWindowOpen(false)
    setSubmitting(false)
    setCurrentName('')
    setOrganizationId(undefined)
  }

  useEffect(() => {
    if (query && query.length >= 3) {
      ConsoleAPI.searchObject('organization', {
        page: page,
        size: size,
        query: query,
      }).then((value) => setList(value))
    } else {
      ConsoleAPI.listOrganizations({ page: page, size: size }).then((value) =>
        setList(value),
      )
    }
  }, [page, size, isSubmitting, query])

  if (!list) {
    return <SectionSpinner />
  }

  return (
    <>
      <ConsoleRenameModal
        closeConfirmationWindow={closeConfirmationWindow}
        isOpen={confirmRenameWindowOpen}
        isSubmitting={isSubmitting}
        previousName={currentName}
        object="organization"
        formSchema={formSchema}
        request={renameOrganization}
      />
      <Helmet>
        <title>Organization Management</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>
          Organization Management
        </Heading>
        {list && list.data.length > 0 ? (
          <DataTable
            items={list.data}
            columns={[
              {
                title: 'Name',
                renderCell: (organization) => (
                  <div
                    className={cx(
                      'flex',
                      'flex-row',
                      'gap-1.5',
                      'items-center',
                    )}
                  >
                    <Avatar
                      name={organization.name}
                      size="sm"
                      className={cx('w-[40px]', 'h-[40px]')}
                    />

                    <ChakraLink
                      as={Link}
                      to={`/console/organizations/${organization.id}`}
                      className={cx('no-underline')}
                    >
                      <Text noOfLines={1}>{organization.name}</Text>
                    </ChakraLink>
                  </div>
                ),
              },
              {
                title: 'Created',
                renderCell: (organization) => (
                  <RelativeDate date={new Date(organization.createTime)} />
                ),
              },
              {
                title: 'Updated',
                renderCell: (organization) => (
                  <RelativeDate date={new Date(organization.updateTime)} />
                ),
              },
            ]}
            actions={[
              {
                label: 'Rename',
                icon: <IconEdit />,
                onClick: async (organization) => {
                  await renameOrganization(
                    organization.id,
                    organization.name,
                    null,
                  )
                },
              },
            ]}
          />
        ) : (
          <div>No organizations found.</div>
        )}
        {list ? (
          <div className={cx('self-end')}>
            <PagePagination
              totalElements={list.totalElements}
              totalPages={Math.ceil(list.totalElements / size)}
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

export default ConsolePanelOrganizations
