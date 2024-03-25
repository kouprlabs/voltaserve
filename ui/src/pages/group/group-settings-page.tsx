import { useMemo, useState } from 'react'
import { useParams } from 'react-router-dom'
import { Divider, IconButton, Text } from '@chakra-ui/react'
import { IconEdit, IconTrash, IconUserPlus, SectionSpinner } from '@koupr/ui'
import classNames from 'classnames'
import { Helmet } from 'react-helmet-async'
import GroupAPI from '@/client/api/group'
import { geEditorPermission, geOwnerPermission } from '@/client/api/permission'
import { swrConfig } from '@/client/options'
import GroupAddMember from '@/components/group/group-add-member'
import GroupDelete from '@/components/group/group-delete'
import GroupEditName from '@/components/group/group-edit-name'

const Spacer = () => <div className={classNames('grow')} />

const GroupSettingsPage = () => {
  const { id } = useParams()
  const { data: group, error } = GroupAPI.useGetById(id, swrConfig())
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
  if (!group) {
    return <SectionSpinner />
  }

  return (
    <>
      <Helmet>
        <title>{group.name}</title>
      </Helmet>
      <div className={sectionClassName}>
        <div className={rowClassName}>
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
        </div>
        <Divider />
        <div className={rowClassName}>
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
        </div>
        <Divider />
        <div className={rowClassName}>
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
        </div>
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
      </div>
    </>
  )
}

export default GroupSettingsPage
