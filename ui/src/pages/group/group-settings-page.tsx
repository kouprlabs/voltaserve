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
import { Divider, IconButton } from '@chakra-ui/react'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import GroupAPI from '@/client/api/group'
import { geEditorPermission, geOwnerPermission } from '@/client/api/permission'
import { swrConfig } from '@/client/options'
import GroupAddMember from '@/components/group/group-add-member'
import GroupDelete from '@/components/group/group-delete'
import GroupEditName from '@/components/group/group-edit-name'
import { IconEdit, IconDelete, IconPersonAdd } from '@/lib/components/icons'
import SectionSpinner from '@/lib/components/section-spinner'
import { truncateEnd } from '@/lib/helpers/truncate-end'

const Spacer = () => <div className={cx('grow')} />

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
          <span>Name</span>
          <Spacer />
          <span>{truncateEnd(group.name, 60)}</span>
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
          <span>Add members</span>
          <Spacer />
          <IconButton
            icon={<IconPersonAdd />}
            isDisabled={!hasOwnerPermission}
            aria-label=""
            onClick={() => {
              setIsAddMembersModalOpen(true)
            }}
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
