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
  Badge,
  Button,
  Heading,
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
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import AdminApi, { InvitationManagementList } from '@/client/admin/admin'
import { adminInvitationsPaginationStorage } from '@/infra/pagination'
import { IconChevronDown, IconChevronUp } from '@/lib/components/icons'
import PagePagination from '@/lib/components/page-pagination'
import SectionSpinner from '@/lib/components/section-spinner'
import usePagePagination from '@/lib/hooks/page-pagination'
import AdminConfirmationModal from '@/pages/admin/admin-confirmation-modal'

const AdminPanelInvitations = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const [list, setList] = useState<InvitationManagementList | undefined>(
    undefined,
  )
  const [isSubmitting, setSubmitting] = useState(false)
  const [invitationId, setInvitationId] = useState<string | undefined>(
    undefined,
  )
  const [invitationDetails, setInvitationDetails] = useState<
    string | undefined
  >(undefined)
  const [actionState, setActionState] = useState<boolean | undefined>(undefined)
  const [confirmInvitationWindowOpen, setConfirmInvitationWindowOpen] =
    useState(false)
  const [confirmWindowAction, setConfirmWindowAction] = useState<
    string | undefined
  >(undefined)
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigate,
    location,
    storage: adminInvitationsPaginationStorage(),
  })

  const changeInvitationStatus = async (
    id: string | null,
    invitation: string | null,
    accept: boolean | null,
    confirm: boolean = false,
  ) => {
    if (confirm && invitationId && actionState !== undefined) {
      setSubmitting(true)
      try {
        await AdminApi.invitationChangeStatus({
          id: invitationId,
          accept: actionState,
        })
      } finally {
        closeConfirmationWindow()
      }
    } else if (id && accept !== null && invitation) {
      setConfirmInvitationWindowOpen(true)
      setActionState(accept)
      setInvitationDetails(invitation)
      setInvitationId(id)
    }
  }

  const closeConfirmationWindow = () => {
    setInvitationId(undefined)
    setInvitationDetails(undefined)
    setActionState(undefined)
    setConfirmInvitationWindowOpen(false)
    setSubmitting(false)
    setConfirmWindowAction(undefined)
  }

  useEffect(() => {
    AdminApi.listInvitations({ page: page, size: size }).then((value) =>
      setList(value),
    )
  }, [page, size, isSubmitting])

  if (!list) {
    return <SectionSpinner />
  }

  return (
    <>
      <AdminConfirmationModal
        isOpen={confirmInvitationWindowOpen}
        action={confirmWindowAction}
        target={invitationDetails}
        closeConfirmationWindow={closeConfirmationWindow}
        isSubmitting={isSubmitting}
        request={changeInvitationStatus}
      />
      <Helmet>
        <title>Invitations management</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>Invitations management</Heading>
        {list && list.data.length > 0 ? (
          <Stack direction="column" spacing={2}>
            <Table variant="simple">
              <Thead>
                <Tr>
                  <Th>Organization</Th>
                  <Th>Invitee</Th>
                  <Th>Status</Th>
                  <Th>Create time</Th>
                  <Th>Update time</Th>
                  <Th>Actions</Th>
                </Tr>
              </Thead>
              <Tbody>
                {list.data.map((invitation) => (
                  <Tr key={invitation.id}>
                    <Td>
                      <Button
                        onClick={() => {
                          navigate(
                            `/admin/organizations/${invitation.organization.id}`,
                          )
                        }}
                      >
                        {invitation.organization.name}
                      </Button>
                    </Td>
                    <Td>
                      <Text>{invitation.email}</Text>
                    </Td>
                    <Td>
                      {invitation.status === 'pending' ? (
                        <Badge colorScheme="yellow">Pending</Badge>
                      ) : invitation.status === 'declined' ? (
                        <Badge colorScheme="red">Declined</Badge>
                      ) : invitation.status === 'accepted' ? (
                        <Badge colorScheme="green">Accepted</Badge>
                      ) : (
                        <Badge colorScheme="gray">Unknown</Badge>
                      )}
                    </Td>
                    <Td>
                      <Text>
                        {new Date(invitation.createTime).toLocaleDateString()}
                      </Text>
                    </Td>
                    <Td>
                      <Text>
                        {new Date(invitation.updateTime).toLocaleString()}
                      </Text>
                    </Td>
                    <Td>
                      {invitation.status === 'pending' ? (
                        <Menu>
                          {({ isOpen }) => (
                            <>
                              <MenuButton
                                isActive={isOpen}
                                as={Button}
                                rightIcon={
                                  isOpen ? (
                                    <IconChevronUp />
                                  ) : (
                                    <IconChevronDown />
                                  )
                                }
                              >
                                Actions
                              </MenuButton>
                              <MenuList>
                                {/*<MenuItem*/}
                                {/*  onClick={async () => {*/}
                                {/*    setConfirmWindowAction('deny invitation')*/}
                                {/*    await changeInvitationStatus(*/}
                                {/*      invitation.id,*/}
                                {/*      `${invitation.email} to ${invitation.organization.name}`,*/}
                                {/*      false,*/}
                                {/*    )*/}
                                {/*  }}*/}
                                {/*>*/}
                                {/*  Accept*/}
                                {/*</MenuItem>*/}
                                <MenuItem
                                  onClick={async () => {
                                    setConfirmWindowAction('deny invitation')
                                    await changeInvitationStatus(
                                      invitation.id,
                                      `${invitation.email} to ${invitation.organization.name}`,
                                      false,
                                    )
                                  }}
                                >
                                  Deny
                                </MenuItem>
                              </MenuList>
                            </>
                          )}
                        </Menu>
                      ) : null}
                    </Td>
                  </Tr>
                ))}
              </Tbody>
            </Table>
          </Stack>
        ) : (
          <div> No invitations found </div>
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

export default AdminPanelInvitations
