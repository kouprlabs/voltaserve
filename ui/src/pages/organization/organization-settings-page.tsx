import { useState } from 'react'
import { useParams } from 'react-router-dom'
import { Divider, IconButton } from '@chakra-ui/react'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import OrganizationAPI from '@/client/api/organization'
import { geEditorPermission, geOwnerPermission } from '@/client/api/permission'
import { swrConfig } from '@/client/options'
import OrganizationDelete from '@/components/organization/organization-delete'
import OrganizationEditName from '@/components/organization/organization-edit-name'
import OrganizationInviteMembers from '@/components/organization/organization-invite-members'
import OrganizationLeave from '@/components/organization/organization-leave'
import {
  IconEdit,
  IconLogout,
  IconDelete,
  IconPersonAdd,
  SectionSpinner,
} from '@/lib'

const Spacer = () => <div className={cx('grow')} />

const OrganizationSettingsPage = () => {
  const { id } = useParams()
  const { data: org, error } = OrganizationAPI.useGet(id, swrConfig())
  const [isNameModalOpen, setIsNameModalOpen] = useState(false)
  const [isInviteMembersModalOpen, setIsInviteMembersModalOpen] =
    useState(false)
  const [isLeaveModalOpen, setIsLeaveModalOpen] = useState(false)
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false)
  const sectionClassName = cx('flex', 'flex-col', 'gap-1', 'py-1.5')
  const rowClassName = cx(
    'flex',
    'flex-row',
    'items-center',
    'gap-1',
    `h-[40px]`,
  )

  if (error) {
    return null
  }

  if (!org) {
    return <SectionSpinner />
  }

  return (
    <>
      <Helmet>
        <title>{org.name}</title>
      </Helmet>
      <div className={sectionClassName}>
        <div className={rowClassName}>
          <span>Name</span>
          <Spacer />
          <span>{org.name}</span>
          <IconButton
            icon={<IconEdit />}
            isDisabled={!geEditorPermission(org.permission)}
            aria-label=""
            onClick={() => {
              setIsNameModalOpen(true)
            }}
          />
        </div>
        <Divider />
        <div className={rowClassName}>
          <span>Invite members</span>
          <Spacer />
          <IconButton
            icon={<IconPersonAdd />}
            isDisabled={!geOwnerPermission(org.permission)}
            aria-label=""
            onClick={() => {
              setIsInviteMembersModalOpen(true)
            }}
          />
        </div>
        <div className={rowClassName}>
          <span>Leave</span>
          <Spacer />
          <IconButton
            icon={<IconLogout />}
            variant="solid"
            colorScheme="red"
            aria-label=""
            onClick={() => setIsLeaveModalOpen(true)}
          />
        </div>
        <Divider />
        <div className={rowClassName}>
          <span>Delete permanently</span>
          <Spacer />
          <IconButton
            icon={<IconDelete />}
            variant="solid"
            colorScheme="red"
            isDisabled={!geEditorPermission(org.permission)}
            aria-label=""
            onClick={() => setIsDeleteModalOpen(true)}
          />
        </div>
        <OrganizationEditName
          open={isNameModalOpen}
          organization={org}
          onClose={() => setIsNameModalOpen(false)}
        />
        <OrganizationInviteMembers
          open={isInviteMembersModalOpen}
          id={org.id}
          onClose={() => setIsInviteMembersModalOpen(false)}
        />
        <OrganizationLeave
          open={isLeaveModalOpen}
          id={org.id}
          onClose={() => setIsLeaveModalOpen(false)}
        />
        <OrganizationDelete
          open={isDeleteModalOpen}
          organization={org}
          onClose={() => setIsDeleteModalOpen(false)}
        />
      </div>
    </>
  )
}

export default OrganizationSettingsPage
