import { useMemo, useState } from 'react'
import { useParams } from 'react-router-dom'
import { Box, Divider, HStack, IconButton, Stack, Text } from '@chakra-ui/react'
import {
  IconEdit,
  IconTrash,
  IconUserPlus,
  SectionSpinner,
  variables,
} from '@koupr/ui'
import { Helmet } from 'react-helmet-async'
import GroupAPI from '@/client/api/group'
import { geEditorPermission, geOwnerPermission } from '@/client/api/permission'
import { swrConfig } from '@/client/options'
import AddMember from '@/components/group/add-member'
import Delete from '@/components/group/delete'
import EditName from '@/components/group/edit-name'

const Spacer = () => <Box flexGrow={1} />

const GroupSettingsPage = () => {
  const params = useParams()
  const { data: group, error } = GroupAPI.useGetById(
    params.id as string,
    swrConfig(),
  )
  const [isNameModalOpen, setIsNameModalOpen] = useState(false)
  const [isAddMembersModalOpen, setIsAddMembersModalOpen] = useState(false)
  const [deleteModalOpen, setDeleteModalOpen] = useState(false)
  const hasEditPermission = useMemo(
    () => group && geEditorPermission(group.permission),
    [group],
  )
  const hasOwnerPermission = useMemo(
    () => group && geOwnerPermission(group.permission),
    [group],
  )

  if (error) {
    return null
  }
  if (!group) {
    return <SectionSpinner />
  }

  return (
    <>
      <Helmet>
        <title>{group.name}</title>
      </Helmet>
      <Stack spacing={variables.spacing} w="100%">
        <HStack spacing={variables.spacing}>
          <Text>Name</Text>
          <Spacer />
          <Text>{group.name}</Text>
          <IconButton
            icon={<IconEdit />}
            isDisabled={!hasEditPermission}
            aria-label=""
            onClick={() => {
              setIsNameModalOpen(true)
            }}
          />
        </HStack>
        <Divider />
        <HStack spacing={variables.spacing}>
          <Text>Add members</Text>
          <Spacer />
          <IconButton
            icon={<IconUserPlus />}
            isDisabled={!hasOwnerPermission}
            aria-label=""
            onClick={() => {
              setIsAddMembersModalOpen(true)
            }}
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
            isDisabled={!hasOwnerPermission}
            aria-label=""
            onClick={() => setDeleteModalOpen(true)}
          />
        </HStack>
        <EditName
          open={isNameModalOpen}
          group={group}
          onClose={() => setIsNameModalOpen(false)}
        />
        <AddMember
          open={isAddMembersModalOpen}
          group={group}
          onClose={() => setIsAddMembersModalOpen(false)}
        />
        <Delete
          open={deleteModalOpen}
          group={group}
          onClose={() => setDeleteModalOpen(false)}
        />
      </Stack>
    </>
  )
}

export default GroupSettingsPage
