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
import GroupAPI from '@/api/group'
import { swrConfig } from '@/api/options'
import { geEditorPermission } from '@/api/permission'
import GroupAddMember from '@/components/group/add-member'
import GroupDelete from '@/components/group/delete'
import GroupEditName from '@/components/group/edit-name'

const Spacer = () => <Box flexGrow={1} />

const GroupSettingsPage = () => {
  const params = useParams()
  const { data: group, error } = GroupAPI.useGetById(
    params.id as string,
    swrConfig()
  )
  const [isNameModalOpen, setIsNameModalOpen] = useState(false)
  const [isAddMembersModalOpen, setIsAddMembersModalOpen] = useState(false)
  const [deleteModalOpen, setDeleteModalOpen] = useState(false)
  const hasEditPermission = useMemo(
    () => group && geEditorPermission(group.permission),
    [group]
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
            disabled={!hasEditPermission}
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
            disabled={!hasEditPermission}
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
            disabled={!hasEditPermission}
            aria-label=""
            onClick={() => setDeleteModalOpen(true)}
          />
        </HStack>
        <GroupEditName
          open={isNameModalOpen}
          group={group}
          onClose={() => setIsNameModalOpen(false)}
        />
        <GroupAddMember
          open={isAddMembersModalOpen}
          group={group}
          onClose={() => setIsAddMembersModalOpen(false)}
        />
        <GroupDelete
          open={deleteModalOpen}
          group={group}
          onClose={() => setDeleteModalOpen(false)}
        />
      </Stack>
    </>
  )
}

export default GroupSettingsPage
