// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { ReactElement, useState } from 'react'
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
import ConsoleAPI, { InvitationManagement } from '@/client/console/console'
import { swrConfig } from '@/client/options'
import ConsoleConfirmationModal, {
  ConsoleConfirmationModalRequest,
} from '@/components/console/console-confirmation-modal'
import { consoleInvitationsPaginationStorage } from '@/infra/pagination'

const ConsolePanelInvitations = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const [isConfirmationOpen, setIsConfirmationOpen] = useState(false)
  const [isConfirmationDestructive, setIsConfirmationDestructive] =
    useState(false)
  const [confirmationHeader, setConfirmationHeader] = useState<ReactElement>()
  const [confirmationBody, setConfirmationBody] = useState<ReactElement>()
  const [confirmationRequest, setConfirmationRequest] =
    useState<ConsoleConfirmationModalRequest>()
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigateFn: navigate,
    searchFn: () => location.search,
    storage: consoleInvitationsPaginationStorage(),
  })
  const { data: list, mutate } = ConsoleAPI.useListObject<InvitationManagement>(
    'invitation',
    { page, size },
    swrConfig(),
  )

  if (!list) {
    return <SectionSpinner />
  }

  return (
    <>
      {confirmationHeader && confirmationBody && confirmationRequest ? (
        <ConsoleConfirmationModal
          header={confirmationHeader}
          body={confirmationBody}
          isDestructive={isConfirmationDestructive}
          isOpen={isConfirmationOpen}
          onClose={() => setIsConfirmationOpen(false)}
          onRequest={confirmationRequest}
        />
      ) : null}
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
                  setConfirmationHeader(<>Deny Invitation</>)
                  setConfirmationBody(
                    <>Are you sure you want to deny this invitation?</>,
                  )
                  setConfirmationRequest(() => async () => {
                    await ConsoleAPI.invitationChangeStatus({
                      id: invitation.id,
                      accept: false,
                    })
                    await mutate()
                  })
                  setIsConfirmationDestructive(true)
                  setIsConfirmationOpen(true)
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
