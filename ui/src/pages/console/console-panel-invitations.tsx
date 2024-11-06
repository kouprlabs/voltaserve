// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useCallback, useEffect, useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { Badge, Button, Heading } from '@chakra-ui/react'
import {
  DataTable,
  IconFrontHand,
  PagePagination,
  RelativeDate,
  SectionSpinner,
  Text,
  usePagePagination,
} from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import ConsoleApi, { InvitationManagementList } from '@/client/console/console'
import ConsoleConfirmationModal from '@/components/console/console-confirmation-modal'
import { consoleInvitationsPaginationStorage } from '@/infra/pagination'
import relativeDate from '@/lib/helpers/relative-date'

const ConsolePanelInvitations = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const [list, setList] = useState<InvitationManagementList>()
  const [isSubmitting, setSubmitting] = useState(false)
  const [invitationId, setInvitationId] = useState<string>()
  const [invitationDetails, setInvitationDetails] = useState<string>()
  const [actionState, setActionState] = useState<boolean>()
  const [confirmInvitationWindowOpen, setConfirmInvitationWindowOpen] =
    useState(false)
  const [confirmWindowAction, setConfirmWindowAction] = useState<string>()
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigateFn: navigate,
    searchFn: () => location.search,
    storage: consoleInvitationsPaginationStorage(),
  })

  const changeInvitationStatus = useCallback(
    async (
      id: string | null,
      invitation: string | null,
      accept: boolean | null,
      confirm: boolean = false,
    ) => {
      if (confirm && invitationId && actionState !== undefined) {
        setSubmitting(true)
        try {
          await ConsoleApi.invitationChangeStatus({
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
    },
    [],
  )

  const closeConfirmationWindow = () => {
    setInvitationId(undefined)
    setInvitationDetails(undefined)
    setActionState(undefined)
    setConfirmInvitationWindowOpen(false)
    setSubmitting(false)
    setConfirmWindowAction(undefined)
  }

  useEffect(() => {
    ConsoleApi.listInvitations({ page: page, size: size }).then((value) =>
      setList(value),
    )
  }, [page, size, isSubmitting])

  if (!list) {
    return <SectionSpinner />
  }

  return (
    <>
      <ConsoleConfirmationModal
        isOpen={confirmInvitationWindowOpen}
        action={confirmWindowAction}
        target={invitationDetails}
        closeConfirmationWindow={closeConfirmationWindow}
        isSubmitting={isSubmitting}
        request={changeInvitationStatus}
      />
      <Helmet>
        <title>Invitation Management</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>Invitation Management</Heading>
        {list && list.data.length > 0 ? (
          <DataTable
            items={list.data}
            columns={[
              {
                title: 'Organization',
                renderCell: (invitation) => (
                  <Button
                    onClick={() => {
                      navigate(
                        `/console/organizations/${invitation.organization.id}`,
                      )
                    }}
                  >
                    {invitation.organization.name}
                  </Button>
                ),
              },
              {
                title: 'Invitee',
                renderCell: (invitation) => <Text>{invitation.email}</Text>,
              },
              {
                title: 'Status',
                renderCell: (invitation) => (
                  <>
                    {invitation.status === 'pending' ? (
                      <Badge colorScheme="yellow">Pending</Badge>
                    ) : invitation.status === 'declined' ? (
                      <Badge colorScheme="red">Declined</Badge>
                    ) : invitation.status === 'accepted' ? (
                      <Badge colorScheme="green">Accepted</Badge>
                    ) : (
                      <Badge colorScheme="gray">Unknown</Badge>
                    )}
                  </>
                ),
              },
              {
                title: 'Created',
                renderCell: (invitation) => (
                  <RelativeDate date={new Date(invitation.createTime)} />
                ),
              },
              {
                title: 'Updated',
                renderCell: (invitation) => (
                  <RelativeDate date={new Date(invitation.updateTime)} />
                ),
              },
            ]}
            actions={[
              {
                label: 'Deny',
                icon: <IconFrontHand />,
                isDestructive: true,
                isDisabledFn: (invitation) => invitation.status !== 'pending',
                onClick: async (invitation) => {
                  setConfirmWindowAction('deny invitation')
                  await changeInvitationStatus(
                    invitation.id,
                    `${invitation.email} to ${invitation.organization.name}`,
                    false,
                  )
                },
              },
            ]}
          />
        ) : (
          <div>No invitations found.</div>
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

export default ConsolePanelInvitations
