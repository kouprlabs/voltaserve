// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useCallback, useEffect, useState } from 'react'
import {
  Button,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  useDisclosure,
  Table,
  Tr,
  Tbody,
  Td,
  Avatar,
  Radio,
} from '@chakra-ui/react'
import {
  Pagination,
  SearchInput,
  SectionError,
  SectionPlaceholder,
  SectionSpinner,
  usePageMonitor,
} from '@koupr/ui'
import cx from 'classnames'
import UserAPI, { SortOrder, User } from '@/client/api/user'
import { swrConfig } from '@/client/options'
import { getPictureUrlById } from '@/lib/helpers/picture'
import userToString from '@/lib/helpers/user-to-string'

export type UserSelectorProps = {
  value?: User
  organizationId?: string
  groupId?: string
  excludeGroupMembers?: boolean
  onConfirm?: (user: User) => void
}

const UserSelector = ({
  value,
  organizationId,
  groupId,
  excludeGroupMembers,
  onConfirm,
}: UserSelectorProps) => {
  const { isOpen, onOpen, onClose } = useDisclosure()
  const [page, setPage] = useState(1)
  const [query, setQuery] = useState('')
  const [selected, setSelected] = useState<User>()
  const size = 5
  const {
    data: list,
    error: listError,
    isLoading: isListLoading,
    mutate,
  } = UserAPI.useList(
    {
      query,
      organizationId,
      groupId,
      excludeGroupMembers,
      page,
      size,
      sortOrder: SortOrder.Desc,
    },
    swrConfig(),
  )
  const { hasPageSwitcher } = usePageMonitor({
    totalPages: list?.totalPages ?? 1,
    totalElements: list?.totalElements ?? 0,
    steps: [size],
  })
  const isListError = !list && listError
  const isListEmpty = list && !listError && list.totalElements === 0
  const isListReady = list && !listError && list.totalElements > 0

  useEffect(() => {
    mutate().then()
  }, [page, query, mutate])

  useEffect(() => {
    if (!isOpen) {
      setPage(1)
      setSelected(undefined)
      setQuery('')
    }
  }, [isOpen])

  const handleConfirm = useCallback(() => {
    if (selected) {
      onConfirm?.(selected)
      onClose()
    }
  }, [selected, onConfirm, onClose])

  const handleSearchInputChange = useCallback((value: string) => {
    setPage(1)
    setQuery(value)
  }, [])

  return (
    <>
      <Button
        variant="outline"
        className={cx(
          'w-full',
          { 'text-black': value },
          { 'dark:text-white': value },
          { 'text-gray-500': !value },
          { 'dark:text-gray-500': !value },
        )}
        onClick={onOpen}
      >
        {value ? userToString(value) : 'Select User'}
      </Button>
      <Modal
        size="xl"
        isOpen={isOpen}
        onClose={onClose}
        closeOnOverlayClick={false}
      >
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Select User</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <div className={cx('flex', 'flex-col', 'gap-1.5')}>
              <SearchInput
                placeholder="Search Users"
                query={query}
                onChange={handleSearchInputChange}
              />
              {isListError ? (
                <SectionError text="Failed to load users." />
              ) : null}
              {isListLoading ? <SectionSpinner /> : null}
              {isListEmpty ? (
                <SectionPlaceholder text="There are no users." />
              ) : null}
              {isListReady ? (
                <div
                  className={cx(
                    'flex',
                    'flex-col',
                    'justify-between',
                    'gap-1.5',
                    'h-[320px]',
                  )}
                >
                  <Table variant="simple" size="sm">
                    <colgroup>
                      <col className={cx('w-[40px]')} />
                      <col className={cx('w-[auto]')} />
                    </colgroup>
                    <Tbody>
                      {list.data.map((u) => (
                        <Tr
                          key={u.id}
                          className={cx(
                            'cursor-pointer',
                            'h-[52px]',
                            { 'bg-gray-100': selected?.id === u.id },
                            { 'dark:bg-gray-600': selected?.id === u.id },
                            { 'bg-transparent': selected?.id !== u.id },
                          )}
                          onClick={() => setSelected(u)}
                        >
                          <Td className={cx('px-0.5', 'text-center')}>
                            <Radio
                              size="md"
                              isChecked={selected?.id === u.id}
                            />
                          </Td>
                          <Td className={cx('p-0.5')}>
                            <div
                              className={cx(
                                'flex',
                                'flex-row',
                                'items-center',
                                'gap-1.5',
                              )}
                            >
                              <Avatar
                                name={u.fullName}
                                src={
                                  u.picture
                                    ? getPictureUrlById(u.id, u.picture, {
                                        organizationId: organizationId,
                                        groupId: groupId,
                                      })
                                    : undefined
                                }
                                size="sm"
                                className={cx(
                                  'w-[40px]',
                                  'h-[40px]',
                                  'border',
                                  'border-gray-300',
                                  'dark:border-gray-700',
                                )}
                              />
                              <span className={cx('text-base')}>
                                {userToString(u)}
                              </span>
                            </div>
                          </Td>
                        </Tr>
                      ))}
                    </Tbody>
                  </Table>
                  {hasPageSwitcher ? (
                    <div className={cx('self-end')}>
                      <Pagination
                        maxButtons={3}
                        page={page}
                        totalPages={list.totalPages}
                        onPageChange={(value) => setPage(value)}
                      />
                    </div>
                  ) : null}
                </div>
              ) : null}
            </div>
          </ModalBody>
          <ModalFooter>
            <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
              <Button
                type="button"
                variant="outline"
                colorScheme="blue"
                onClick={onClose}
              >
                Cancel
              </Button>
              <Button
                variant="solid"
                colorScheme="blue"
                isDisabled={!selected}
                onClick={handleConfirm}
              >
                Confirm
              </Button>
            </div>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  )
}

export default UserSelector
