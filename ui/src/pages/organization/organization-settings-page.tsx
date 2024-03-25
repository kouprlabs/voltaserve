import { useState } from 'react'
import { useParams } from 'react-router-dom'
import { Divider, IconButton, Text } from '@chakra-ui/react'
import {
  IconEdit,
  IconExit,
  IconTrash,
  IconUserPlus,
  SectionSpinner,
} from '@koupr/ui'
import classNames from 'classnames'
import { Helmet } from 'react-helmet-async'
import OrganizationAPI from '@/client/api/organization'
import { geEditorPermission, geOwnerPermission } from '@/client/api/permission'
import { swrConfig } from '@/client/options'
import OrganizationDelete from '@/components/organization/organization-delete'
import OrganizationEditName from '@/components/organization/organization-edit-name'
import OrganizationInviteMembers from '@/components/organization/organization-invite-members'
import OrganizationLeave from '@/components/organization/organization-leave'

const Spacer = () => <div className={classNames('grow')} />

const OrganizationSettingsPage = () => {
  const { id } = useParams()
  const { data: org, error } = OrganizationAPI.useGetById(id, swrConfig())
  const [isNameModalOpen, setIsNameModalOpen] = useState(false)
  const [isInviteMembersModalOpen, setIsInviteMembersModalOpen] =
    useState(false)
  const [isLeaveModalOpen, setIsLeaveModalOpen] = useState(false)
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false)
  const sectionClassName = classNames('flex', 'flex-col', 'gap-1', 'py-1.5')
  const rowClassName = classNames(
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
          <Text>Name</Text>
          <Spacer />
          <Text>{org.name}</Text>
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
          <Text>Invite members</Text>
          <Spacer />
          <IconButton
            icon={<IconUserPlus />}
            isDisabled={!geOwnerPermission(org.permission)}
            aria-label=""
            onClick={() => {
              setIsInviteMembersModalOpen(true)
            }}
          />
        </div>
        <div className={rowClassName}>
          <Text>Leave</Text>
          <Spacer />
          <IconButton
            icon={<IconExit />}
            variant="solid"
            colorScheme="red"
            aria-label=""
            onClick={() => setIsLeaveModalOpen(true)}
          />
        </div>
        <Divider />
        <div className={rowClassName}>
          <Text>Delete permanently</Text>
          <Spacer />
          <IconButton
            icon={<IconTrash />}
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
