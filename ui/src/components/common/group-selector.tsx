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
import { SectionSpinner, Pagination, SearchInput } from '@koupr/ui'
import classNames from 'classnames'
import GroupAPI, { Group, SortOrder } from '@/client/api/group'
import { swrConfig } from '@/client/options'

export type GroupSelectorProps = {
  value?: Group
  organizationId?: string
  onConfirm?: (group: Group) => void
}

const GroupSelector = ({
  value,
  organizationId,
  onConfirm,
}: GroupSelectorProps) => {
  const { isOpen, onOpen, onClose } = useDisclosure()
  const [page, setPage] = useState(1)
  const [query, setQuery] = useState('')
  const [selected, setSelected] = useState<Group>()
  const {
    data: list,
    error,
    mutate,
  } = GroupAPI.useList(
    { query, organizationId, page, size: 5, sortOrder: SortOrder.Desc },
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
        onClick={onOpen}
        style={{
          color: value ? normalButtonLabelColor : dimmedButtonLabelColor,
        }}
      >
        {value ? value.name : 'Select Group'}
      </Button>
      <Modal
        size="xl"
        isOpen={isOpen}
        onClose={onClose}
        closeOnOverlayClick={false}
      >
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Select Group</ModalHeader>
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
                  <Text>Failed to load groups.</Text>
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
                    <Text>There are no organizations.</Text>
                  </div>
                </div>
              )}
              {list && list.data.length > 0 && (
                <Table variant="simple" size="sm">
                  <colgroup>
                    <col style={{ width: '40px' }} />
                    <col style={{ width: 'auto' }} />
                  </colgroup>
                  <Tbody>
                    {list.data.map((g) => (
                      <Tr
                        key={g.id}
                        className={classNames('cursor-pointer')}
                        style={{
                          backgroundColor:
                            selected?.id === g.id
                              ? selectionColor
                              : 'transparent',
                        }}
                        onClick={() => setSelected(g)}
                      >
                        <Td className={classNames('px-0.5', 'text-center')}>
                          <Radio size="md" isChecked={selected?.id === g.id} />
                        </Td>
                        <Td className={classNames('px-0.5')}>
                          <div
                            className={classNames(
                              'flex',
                              'flex-row',
                              'items-center',
                              'gap-1.5',
                            )}
                          >
                            <Avatar
                              name={g.name}
                              size="sm"
                              className={classNames('w-[40px]', 'h-[40px]')}
                            />
                            <Text className={classNames('text-base')}>
                              {g.name}
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

export default GroupSelector
