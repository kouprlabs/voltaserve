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
import AdminApi, { GroupManagementList } from '@/client/admin/admin'
import { adminGroupsPaginationStorage } from '@/infra/pagination'
import PagePagination from '@/lib/components/page-pagination'
import SectionSpinner from '@/lib/components/section-spinner'
import usePagePagination from '@/lib/hooks/page-pagination'
import AdminRenameModal from '@/pages/admin/admin-rename-modal'

const AdminPanelGroups = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const [list, setList] = useState<GroupManagementList | undefined>(undefined)
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigate,
    location,
    storage: adminGroupsPaginationStorage(),
  })
  const [confirmWindowOpen, setConfirmWindowOpen] = useState(false)
  const [isSubmitting, setSubmitting] = useState(false)
  const [currentName, setCurrentName] = useState<string | undefined>(undefined)
  const [groupId, setGroupId] = useState<string | undefined>(undefined)
  const formSchema = Yup.object().shape({
    name: Yup.string().required('Name is required').max(255),
  })

  const renameGroup = async (
    id: string | null,
    currentName: string | null,
    newName: string | null,
    confirm: boolean = false,
  ) => {
    if (confirm && groupId !== undefined && newName !== null) {
      try {
        setSubmitting(true)
        await AdminApi.renameObject({ id: groupId, name: newName }, 'group')
      } finally {
        closeConfirmationWindow()
      }
    } else if (id !== null && currentName !== null) {
      setConfirmWindowOpen(true)
      setCurrentName(currentName)
      setGroupId(id)
    }
  }

  const closeConfirmationWindow = () => {
    setConfirmWindowOpen(false)
    setSubmitting(false)
    setCurrentName(undefined)
    setGroupId(undefined)
  }

  useEffect(() => {
    AdminApi.listGroups({ page: page, size: size }).then((value) =>
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
        isOpen={confirmWindowOpen}
        isSubmitting={isSubmitting}
        previousName={currentName}
        object="group"
        formSchema={formSchema}
        request={renameGroup}
      />
      <Helmet>
        <title>Groups management</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>Groups management</Heading>
        {list && list.data.length > 0 ? (
          <Stack direction="column" spacing={2}>
            <Table variant="simple">
              <Thead>
                <Tr>
                  <Th>Group name</Th>
                  <Th>Organization</Th>
                  <Th>Create time</Th>
                  <Th>Update time</Th>
                  <Th>Actions</Th>
                </Tr>
              </Thead>
              <Tbody>
                {list.data.map((group) => (
                  <Tr key={group.id}>
                    <Td>
                      <Text>{group.name}</Text>
                    </Td>
                    <Td>
                      <Button
                        onClick={() => {
                          navigate(
                            `/admin/organizations/${group.organization.id}`,
                          )
                        }}
                      >
                        {group.organization.name}
                      </Button>
                    </Td>
                    <Td>
                      <Text>
                        {new Date(group.createTime).toLocaleDateString()}
                      </Text>
                    </Td>
                    <Td>
                      <Text>{new Date(group.updateTime).toLocaleString()}</Text>
                    </Td>
                    <Td>
                      <Button
                        onClick={async () => {
                          await renameGroup(group.id, group.name, null)
                        }}
                      >
                        Rename
                      </Button>
                    </Td>
                  </Tr>
                ))}
              </Tbody>
            </Table>
          </Stack>
        ) : (
          <div> No groups found </div>
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

export default AdminPanelGroups
