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
import { useLocation, useNavigate, useSearchParams } from 'react-router-dom'
import {
  Avatar,
  Center,
  Heading,
  IconButton,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  Stack,
  Table,
  Tbody,
  Td,
  Text,
  Th,
  Thead,
  Tr,
} from '@chakra-ui/react'
import * as Yup from 'yup'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import ConsoleApi, {
  OrganizationManagementList,
} from '@/client/console/console'
import ConsoleRenameModal from '@/components/console/console-rename-modal'
import { consoleOrganizationsPaginationStorage } from '@/infra/pagination'
import { IconMoreVert } from '@/lib/components/icons'
import PagePagination from '@/lib/components/page-pagination'
import SectionSpinner from '@/lib/components/section-spinner'
import { decodeQuery } from '@/lib/helpers/query'
import usePagePagination from '@/lib/hooks/page-pagination'

const ConsolePanelOrganizations = () => {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const location = useLocation()
  const query = decodeQuery(searchParams.get('q') as string)
  const [list, setList] = useState<OrganizationManagementList>()
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigate,
    location,
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
        await ConsoleApi.renameObject(
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
      ConsoleApi.searchObject('organization', {
        page: page,
        size: size,
        query: query,
      }).then((value) => setList(value))
    } else {
      ConsoleApi.listOrganizations({ page: page, size: size }).then((value) =>
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
          <Stack direction="column" spacing={2}>
            <Table variant="simple">
              <Thead>
                <Tr>
                  <Th>Organization name</Th>
                  <Th>Create time</Th>
                  <Th>Update time</Th>
                  <Th></Th>
                </Tr>
              </Thead>
              <Tbody>
                {list.data.map((organization) => (
                  <Tr
                    style={{ cursor: 'pointer' }}
                    key={organization.id}
                    onClick={(event) => {
                      if (
                        !(event.target instanceof HTMLButtonElement) &&
                        !(event.target instanceof HTMLSpanElement)
                      ) {
                        navigate(`/console/organizations/${organization.id}`)
                      }
                    }}
                  >
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
                          name={organization.name}
                          size="sm"
                          className={cx('w-[40px]', 'h-[40px]')}
                        />
                        <Text noOfLines={1}>{organization.name}</Text>
                      </div>
                    </Td>
                    <Td>
                      <Text>
                        {new Date(organization.createTime).toLocaleDateString()}
                      </Text>
                    </Td>
                    <Td>
                      <Text>
                        {new Date(organization.updateTime).toLocaleString()}
                      </Text>
                    </Td>
                    <Td>
                      <Center>
                        <Menu>
                          <MenuButton
                            as={IconButton}
                            icon={<IconMoreVert />}
                            variant="ghost"
                            aria-label=""
                          />
                          <MenuList>
                            <MenuItem
                              onClick={async () => {
                                await renameOrganization(
                                  organization.id,
                                  organization.name,
                                  null,
                                )
                              }}
                            >
                              Rename
                            </MenuItem>
                          </MenuList>
                        </Menu>
                      </Center>
                    </Td>
                  </Tr>
                ))}
              </Tbody>
            </Table>
          </Stack>
        ) : (
          <div>No organizations found.</div>
        )}
        {list ? (
          <PagePagination
            style={{ alignSelf: 'end' }}
            totalElements={list.totalElements}
            totalPages={Math.ceil(list.totalElements / size)}
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

export default ConsolePanelOrganizations
