import { useState } from 'react'
import { useParams } from 'react-router-dom'
import { Box, Divider, HStack, IconButton, Stack, Text } from '@chakra-ui/react'
import { Helmet } from 'react-helmet-async'
import { swrConfig } from '@/api/options'
import OrganizationAPI from '@/api/organization'
import { geEditorPermission } from '@/api/permission'
import {
  IconEdit,
  IconExit,
  IconTrash,
  IconUserPlus,
} from '@/components/common/icon'
import LoadingSpinner from '@/components/common/loading-spinner'
import OrganizationDelete from '@/components/organization/delete'
import OrganizationEditName from '@/components/organization/edit-name'
import OrganizationInviteMembers from '@/components/organization/invite-members'
import OrganizationLeave from '@/components/organization/leave'
import variables from '@/theme/variables'

const Spacer = () => <Box flexGrow={1} />

const OrganizationSettingsPage = () => {
  const params = useParams()
  const { data: org, error } = OrganizationAPI.useGetById(
    params.id as string,
    swrConfig()
  )
  const [isNameModalOpen, setIsNameModalOpen] = useState(false)
  const [isInviteMembersModalOpen, setIsInviteMembersModalOpen] =
    useState(false)
  const [isLeaveModalOpen, setIsLeaveModalOpen] = useState(false)
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false)

  if (error) {
    return null
  }

  if (!org) {
    return <LoadingSpinner />
  }

  return (
    <>
      <Helmet>
        <title>{org.name}</title>
      </Helmet>
      <Stack spacing={variables.spacing} w="100%">
        <HStack spacing={variables.spacing}>
          <Text>Name</Text>
          <Spacer />
          <Text>{org.name}</Text>
          <IconButton
            icon={<IconEdit />}
            disabled={!geEditorPermission(org.permission)}
            aria-label=""
            onClick={() => {
              setIsNameModalOpen(true)
            }}
          />
        </HStack>
        <Divider />
        <HStack spacing={variables.spacing}>
          <Text>Invite members</Text>
          <Spacer />
          <IconButton
            icon={<IconUserPlus />}
            disabled={!geEditorPermission(org.permission)}
            aria-label=""
            onClick={() => {
              setIsInviteMembersModalOpen(true)
            }}
          />
        </HStack>
        <HStack spacing={variables.spacing}>
          <Text>Leave</Text>
          <Spacer />
          <IconButton
            icon={<IconExit />}
            variant="solid"
            colorScheme="red"
            aria-label=""
            onClick={() => setIsLeaveModalOpen(true)}
          />
        </HStack>
        <Divider />
        <HStack spacing={variables.spacing}>
          <Text>Delete permanently</Text>
          <Spacer />
          <IconButton
            icon={<IconTrash />}
            variant="solid"
            colorScheme="red"
            disabled={!geEditorPermission(org.permission)}
            aria-label=""
            onClick={() => setIsDeleteModalOpen(true)}
          />
        </HStack>
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
      </Stack>
    </>
  )
}

export default OrganizationSettingsPage
