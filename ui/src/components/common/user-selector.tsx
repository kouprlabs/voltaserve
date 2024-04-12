import { useCallback, useEffect, useState } from 'react'
import {
  Button,
  Text,
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
import { SectionSpinner, Pagination, SearchInput } from '@koupr/ui'
import cx from 'classnames'
import UserAPI, { SortOrder, User } from '@/client/api/user'
import { swrConfig } from '@/client/options'
import userToString from '@/helpers/user-to-string'

export type UserSelectorProps = {
  value?: User
  organizationId?: string
  groupId?: string
  nonGroupMembersOnly?: boolean
  onConfirm?: (user: User) => void
}

const UserSelector = ({
  value,
  organizationId,
  groupId,
  nonGroupMembersOnly,
  onConfirm,
}: UserSelectorProps) => {
  const { isOpen, onOpen, onClose } = useDisclosure()
  const [page, setPage] = useState(1)
  const [query, setQuery] = useState('')
  const [selected, setSelected] = useState<User>()
  const {
    data: list,
    error,
    mutate,
  } = UserAPI.useList(
    {
      query,
      organizationId,
      groupId,
      nonGroupMembersOnly,
      page,
      size: 5,
      sortOrder: SortOrder.Desc,
    },
    swrConfig(),
  )

  useEffect(() => {
    mutate()
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
                query={query}
                onChange={(value) => setQuery(value)}
              />
              {!list && error && (
                <div
                  className={cx(
                    'flex',
                    'items-center',
                    'justify-center',
                    'h-[300px]',
                  )}
                >
                  <Text>Failed to load users.</Text>
                </div>
              )}
              {!list && !error && <SectionSpinner />}
              {list && list.data.length === 0 && (
                <div
                  className={cx(
                    'flex',
                    'items-center',
                    'justify-center',
                    'h-[300px]',
                  )}
                >
                  <div
                    className={cx(
                      'flex',
                      'flex-col',
                      'items-center',
                      'gap-1.5',
                    )}
                  >
                    <Text>There are no users.</Text>
                  </div>
                </div>
              )}
              {list && list.data.length > 0 && (
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
                          { 'bg-gray-100': selected?.id === u.id },
                          { 'dark:bg-gray-600': selected?.id === u.id },
                          { 'bg-transparent': selected?.id !== u.id },
                        )}
                        onClick={() => setSelected(u)}
                      >
                        <Td className={cx('px-0.5', 'text-center')}>
                          <Radio size="md" isChecked={selected?.id === u.id} />
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
                              size="sm"
                              className={cx('w-[40px]', 'h-[40px]')}
                            />
                            <Text className={cx('text-base')}>
                              {userToString(u)}
                            </Text>
                          </div>
                        </Td>
                      </Tr>
                    ))}
                  </Tbody>
                </Table>
              )}
              {list && (
                <div className={cx('self-end')}>
                  {list.totalPages > 1 ? (
                    <Pagination
                      uiSize="md"
                      maxButtons={3}
                      page={page}
                      totalPages={list.totalPages}
                      onPageChange={(value) => setPage(value)}
                    />
                  ) : null}
                </div>
              )}
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
