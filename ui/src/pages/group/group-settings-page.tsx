// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useMemo, useState } from 'react'
import { useParams } from 'react-router-dom'
import { IconButton } from '@chakra-ui/react'
import {
  IconEdit,
  IconDelete,
  IconPersonAdd,
  SectionSpinner,
  Form,
} from '@koupr/ui'
import { Helmet } from 'react-helmet-async'
import GroupAPI from '@/client/api/group'
import { geEditorPermission, geOwnerPermission } from '@/client/api/permission'
import { swrConfig } from '@/client/options'
import GroupAddMember from '@/components/group/group-add-member'
import GroupDelete from '@/components/group/group-delete'
import GroupEditName from '@/components/group/group-edit-name'
import { truncateEnd } from '@/lib/helpers/truncate-end'

const GroupSettingsPage = () => {
  const { id } = useParams()
  const { data: group, error } = GroupAPI.useGet(id, swrConfig())
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
      <Form
        sections={[
          {
            title: 'Basics',
            rows: [
              {
                label: 'Name',
                content: (
                  <>
                    <span>{truncateEnd(group.name, 50)}</span>
                    <IconButton
                      icon={<IconEdit />}
                      isDisabled={!hasEditPermission}
                      aria-label=""
                      onClick={() => {
                        setIsNameModalOpen(true)
                      }}
                    />
                  </>
                ),
              },
            ],
          },
          {
            title: 'Membership',
            rows: [
              {
                label: 'Add members',
                content: (
                  <IconButton
                    icon={<IconPersonAdd />}
                    isDisabled={!hasOwnerPermission}
                    aria-label=""
                    onClick={() => {
                      setIsAddMembersModalOpen(true)
                    }}
                  />
                ),
              },
            ],
          },
          {
            title: 'Advanced',
            rows: [
              {
                label: 'Delete group',
                content: (
                  <IconButton
                    icon={<IconDelete />}
                    variant="solid"
                    colorScheme="red"
                    isDisabled={!hasOwnerPermission}
                    aria-label=""
                    onClick={() => setDeleteModalOpen(true)}
                  />
                ),
              },
            ],
          },
        ]}
      />
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
    </>
  )
}

export default GroupSettingsPage
