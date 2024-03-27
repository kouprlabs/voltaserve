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
  useColorModeValue,
  useToken,
} from '@chakra-ui/react'
import { SectionSpinner, Pagination, SearchInput, variables } from '@koupr/ui'
import classNames from 'classnames'
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
  const selectionColor = useToken(
    'colors',
    useColorModeValue('gray.100', 'gray.600'),
  )
  const dimmedButtonLabelColor = useToken(
    'colors',
    useColorModeValue('gray.500', 'gray.500'),
  )
  const normalButtonLabelColor = useToken(
    'colors',
    useColorModeValue('black', 'white'),
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
        className={classNames('w-full')}
        style={{
          color: value ? normalButtonLabelColor : dimmedButtonLabelColor,
        }}
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
            <div className={classNames('flex', 'flex-col', 'gap-1.5')}>
              <SearchInput
                query={query}
                onChange={(value) => setQuery(value)}
              />
              {!list && error && (
                <div
                  className={classNames(
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
                  className={classNames(
                    'flex',
                    'items-center',
                    'justify-center',
                    'h-[300px]',
                  )}
                >
                  <div
                    className={classNames(
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
                    <col className={classNames('w-[40px]')} />
                    <col className={classNames('w-[auto]')} />
                  </colgroup>
                  <Tbody>
                    {list.data.map((u) => (
                      <Tr
                        key={u.id}
                        className={classNames('cursor-pointer')}
                        bg={selected?.id === u.id ? selectionColor : 'auto'}
                        onClick={() => setSelected(u)}
                      >
                        <Td className={classNames('px-0.5', 'text-center')}>
                          <Radio size="md" isChecked={selected?.id === u.id} />
                        </Td>
                        <Td px={variables.spacingXs}>
                          <div
                            className={classNames(
                              'flex',
                              'flex-row',
                              'items-center',
                              'gap-1.5',
                            )}
                          >
                            <Avatar
                              name={u.fullName}
                              size="sm"
                              className={classNames('w-[40px]', 'h-[40px]')}
                            />
                            <Text className={classNames('text-base')}>
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
                <div className={classNames('self-end')}>
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
            <Button
              type="button"
              variant="outline"
              colorScheme="blue"
              className={classNames('mr-1')}
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
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  )
}

export default UserSelector
