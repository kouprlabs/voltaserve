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
import { useLocation, useNavigate } from 'react-router-dom'
import {
  Button,
  Heading,
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
import AdminApi, { OrganizationManagementList } from '@/client/admin/admin'
import AdminHighlightableTr from '@/components/admin/admin-highlightable-tr'
import AdminRenameModal from '@/components/admin/admin-rename-modal'
import { adminOrganizationsPaginationStorage } from '@/infra/pagination'
import PagePagination from '@/lib/components/page-pagination'
import SectionSpinner from '@/lib/components/section-spinner'
import usePagePagination from '@/lib/hooks/page-pagination'

const AdminPanelOrganizations = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const [list, setList] = useState<OrganizationManagementList | undefined>(
    undefined,
  )
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigate,
    location,
    storage: adminOrganizationsPaginationStorage(),
  })
  const [confirmRenameWindowOpen, setConfirmRenameWindowOpen] = useState(false)
  const [isSubmitting, setSubmitting] = useState(false)
  const [currentName, setCurrentName] = useState<string | undefined>(undefined)
  const [organizationId, setOrganizationId] = useState<string | undefined>(
    undefined,
  )
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
        await AdminApi.renameObject(
          { id: organizationId, name: newName },
          'organization',
        )
      } finally {
        closeConfirmationWindow()
      }
    } else if (id !== null && currentName !== null) {
      setConfirmRenameWindowOpen(true)
      setCurrentName(currentName)
      setOrganizationId(id)
    }
  }

  const closeConfirmationWindow = () => {
    setConfirmRenameWindowOpen(false)
    setSubmitting(false)
    setCurrentName(undefined)
    setOrganizationId(undefined)
  }

  useEffect(() => {
    AdminApi.listOrganizations({ page: page, size: size }).then((value) =>
      setList(value),
    )
  }, [page, size, isSubmitting])

  if (!list) {
    return <SectionSpinner />
  }

  return (
    <>
      <AdminRenameModal
        closeConfirmationWindow={closeConfirmationWindow}
        isOpen={confirmRenameWindowOpen}
        isSubmitting={isSubmitting}
        previousName={currentName}
        object="organization"
        formSchema={formSchema}
        request={renameOrganization}
      />
      <Helmet>
        <title>Organizations management</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>
          Organizations management
        </Heading>
        {list && list.data.length > 0 ? (
          <Stack direction="column" spacing={2}>
            <Table variant="simple">
              <Thead>
                <Tr>
                  <Th>Workspace name</Th>
                  <Th>Create time</Th>
                  <Th>Update time</Th>
                  <Th>Actions</Th>
                </Tr>
              </Thead>
              <Tbody>
                {list.data.map((organization) => (
                  <AdminHighlightableTr
                    key={organization.id}
                    onClick={(event) => {
                      if (
                        !(event.target instanceof HTMLButtonElement) &&
                        !(event.target instanceof HTMLSpanElement) &&
                        !(event.target instanceof HTMLParagraphElement)
                      ) {
                        navigate(`/admin/organizations/${organization.id}`)
                      }
                    }}
                  >
                    <Td>
                      <Text>{organization.name}</Text>
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
                      <Button
                        onClick={async () => {
                          await renameOrganization(
                            organization.id,
                            organization.name,
                            null,
                          )
                        }}
                      >
                        Rename
                      </Button>
                    </Td>
                  </AdminHighlightableTr>
                ))}
              </Tbody>
            </Table>
          </Stack>
        ) : (
          <div> No organizations found </div>
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

export default AdminPanelOrganizations
