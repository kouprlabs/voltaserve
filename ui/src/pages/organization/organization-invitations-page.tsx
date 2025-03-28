// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useCallback, useEffect, useState } from 'react'
import { useLocation, useNavigate, useParams } from 'react-router-dom'
import { Button, useToast } from '@chakra-ui/react'
import {
  DataTable,
  IconDelete,
  IconPersonAdd,
  IconSend,
  Text,
  PagePagination,
  SectionSpinner,
  usePagePagination,
  RelativeDate,
  SectionError,
  SectionPlaceholder,
} from '@koupr/ui'
import {
  InvitationAPI,
  InvitationSortBy,
  InvitationSortOrder,
} from '@/client/api/invitation'
import { OrganizationAPI } from '@/client/api/organization'
import { geOwnerPermission } from '@/client/api/permission'
import { errorToString } from '@/client/error'
import { swrConfig } from '@/client/options'
import OrganizationInviteMembers from '@/components/organization/organization-invite-members'
import OrganizationStatus from '@/components/organization/organization-status'
import { outgoingInvitationPaginationStorage } from '@/infra/pagination'
import truncateMiddle from '@/lib/helpers/truncate-middle'
import { useAppDispatch } from '@/store/hook'
import { mutateUpdated } from '@/store/ui/outgoing-invitations'

const OrganizationInvitationsPage = () => {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const location = useLocation()
  const { id } = useParams()
  const toast = useToast()
  const {
    data: org,
    error: orgError,
    isLoading: orgIsLoading,
  } = OrganizationAPI.useGet(id, swrConfig())
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigateFn: navigate,
    searchFn: () => location.search,
    storage: outgoingInvitationPaginationStorage(),
  })
  const {
    data: list,
    error: listError,
    isLoading: listIsLoading,
    mutate,
  } = InvitationAPI.useGetOutgoing(
    {
      organizationId: id,
      page,
      size,
      sortBy: InvitationSortBy.DateCreated,
      sortOrder: InvitationSortOrder.Desc,
    },
    swrConfig(),
  )
  // prettier-ignore
  const [isInviteMembersModalOpen, setIsInviteMembersModalOpen] = useState(false)
  const orgIsReady = org && !orgError
  const listIsEmpty = list && !listError && list.totalElements === 0
  const listIsReady = list && !listError && list.totalElements > 0

  const handleResend = useCallback(
    async (invitationId: string) => {
      await InvitationAPI.resend(invitationId)
      toast({
        title: 'Invitation resent',
        status: 'success',
        isClosable: true,
      })
    },
    [toast],
  )

  const handleDelete = useCallback(
    async (invitationId: string) => {
      await InvitationAPI.delete(invitationId)
      await mutate()
    },
    [mutate],
  )

  useEffect(() => {
    if (mutate) {
      dispatch(mutateUpdated(mutate))
    }
  }, [mutate])

  return (
    <>
      {orgIsLoading ? <SectionSpinner /> : null}
      {orgError ? <SectionError text={errorToString(orgError)} /> : null}
      {orgIsReady ? (
        <>
          {listIsLoading ? <SectionSpinner /> : null}
          {listError ? <SectionError text={errorToString(listError)} /> : null}
          {listIsEmpty ? (
            <SectionPlaceholder
              text="There are no items."
              content={
                geOwnerPermission(org.permission) ? (
                  <Button
                    leftIcon={<IconPersonAdd />}
                    onClick={() => {
                      setIsInviteMembersModalOpen(true)
                    }}
                  >
                    Invite Members
                  </Button>
                ) : undefined
              }
            />
          ) : null}
          {listIsReady ? (
            <DataTable
              items={list.data}
              columns={[
                {
                  title: 'Email',
                  renderCell: (i) => <Text>{truncateMiddle(i.email, 50)}</Text>,
                },
                {
                  title: 'Status',
                  renderCell: (i) => <OrganizationStatus value={i.status} />,
                },
                {
                  title: 'Date',
                  renderCell: (i) => (
                    <RelativeDate date={new Date(i.createTime)} />
                  ),
                },
              ]}
              actions={[
                {
                  label: 'Resend',
                  icon: <IconSend />,
                  onClick: (i) => handleResend(i.id),
                },
                {
                  label: 'Delete',
                  icon: <IconDelete />,
                  isDestructive: true,
                  onClick: (i) => handleDelete(i.id),
                },
              ]}
              pagination={
                list.totalPages ? (
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
          <OrganizationInviteMembers
            open={isInviteMembersModalOpen}
            id={org.id}
            onClose={() => setIsInviteMembersModalOpen(false)}
          />
        </>
      ) : null}
    </>
  )
}

export default OrganizationInvitationsPage
