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
} from '@koupr/ui'
import cx from 'classnames'
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
    error: listError,
    isLoading: isListLoading,
    mutate,
  } = OrganizationAPI.useList(
    { query, page, size: 5, sortOrder: SortOrder.Desc },
    swrConfig(),
  )
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
      setConfirmed(selected)
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
          { 'text-black': confirmed },
          { 'dark:text-white': confirmed },
          { 'text-gray-500': !confirmed },
          { 'dark:text-gray-500': !confirmed },
        )}
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
            <div className={cx('flex', 'flex-col', 'gap-1.5')}>
              <SearchInput
                placeholder="Search Organizations"
                query={query}
                onChange={handleSearchInputChange}
              />
              {isListLoading ? <SectionSpinner /> : null}
              {isListError ? (
                <SectionError text="Failed to load organizations." />
              ) : null}
              {isListEmpty ? (
                <SectionPlaceholder text="There are no organizations." />
              ) : null}
              {isListReady ? (
                <div
                  className={cx(
                    'flex',
                    'flex-col',
                    'justify-between',
                    'h-[320px]',
                  )}
                >
                  <Table variant="simple" size="sm">
                    <colgroup>
                      <col className={cx('w-[40px]')} />
                      <col className={cx('w-auto')} />
                    </colgroup>
                    <Tbody>
                      {list.data.map((o) => (
                        <Tr
                          key={o.id}
                          className={cx(
                            'cursor-pointer',
                            { 'bg-gray-100': selected?.id === o.id },
                            { 'dark:bg-gray-600': selected?.id === o.id },
                            { 'bg-transparent': selected?.id !== o.id },
                          )}
                          onClick={() => setSelected(o)}
                        >
                          <Td className={cx('px-0.5', 'text-center')}>
                            <Radio
                              size="md"
                              isChecked={selected?.id === o.id}
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
                                name={o.name}
                                size="sm"
                                className={cx('w-[40px]', 'h-[40px]')}
                              />
                              <span className={cx('text-base')}>{o.name}</span>
                            </div>
                          </Td>
                        </Tr>
                      ))}
                    </Tbody>
                  </Table>
                  <div className={cx('self-end')}>
                    {list.totalPages > 1 ? (
                      <Pagination
                        maxButtons={3}
                        page={page}
                        totalPages={list.totalPages}
                        onPageChange={(value) => setPage(value)}
                      />
                    ) : null}
                  </div>
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

export default OrganizationSelector
