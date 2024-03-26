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
import OrganizationAPI, {
  Organization,
  SortOrder,
} from '@/client/api/organization'
import { swrConfig } from '@/client/options'

export type OrganizationSelectorProps = {
  onConfirm?: (organization: Organization) => void
}

const OrganizationSelector = ({ onConfirm }: OrganizationSelectorProps) => {
  const { isOpen, onOpen, onClose } = useDisclosure()
  const [page, setPage] = useState(1)
  const [query, setQuery] = useState('')
  const [selected, setSelected] = useState<Organization>()
  const [confirmed, setConfirmed] = useState<Organization>()
  const {
    data: list,
    error,
    mutate,
  } = OrganizationAPI.useList(
    { query, page, size: 5, sortOrder: SortOrder.Desc },
    swrConfig(),
  )
  const selectionColor = useToken(
    'colors',
    useColorModeValue('gray.100', 'gray.600'),
  )
  const dimmedButtonLabelColor = useColorModeValue('gray.500', 'gray.500')
  const normalButtonLabelColor = useColorModeValue('black', 'white')

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
      setConfirmed(selected)
      onConfirm?.(selected)
      onClose()
    }
  }, [selected, onConfirm, onClose])

  return (
    <>
      <Button
        variant="outline"
        w="100%"
        style={{
          color: confirmed ? normalButtonLabelColor : dimmedButtonLabelColor,
        }}
        onClick={onOpen}
      >
        {confirmed ? confirmed.name : 'Select Organization'}
      </Button>
      <Modal
        size="xl"
        isOpen={isOpen}
        onClose={onClose}
        closeOnOverlayClick={false}
      >
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Select Organization</ModalHeader>
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
                  <Text>Failed to load organizations.</Text>
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
                    {list.data.map((o) => (
                      <Tr
                        key={o.id}
                        cursor="pointer"
                        style={{
                          backgroundColor:
                            selected?.id === o.id
                              ? selectionColor
                              : 'transparent',
                        }}
                        onClick={() => setSelected(o)}
                      >
                        <Td className={classNames('px-0.5', 'text-center')}>
                          <Radio size="md" isChecked={selected?.id === o.id} />
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
                              name={o.name}
                              size="sm"
                              className={classNames('w-[40px]', 'h-[40px]')}
                            />
                            <Text fontSize={variables.bodyFontSize}>
                              {o.name}
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
              mr={variables.spacingSm}
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

export default OrganizationSelector
