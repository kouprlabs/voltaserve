// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
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
import { GroupAPI, Group, GroupSortOrder } from '@/client/api/group'
import { errorToString } from '@/client/error'
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
  const size = 5
  const {
    data: list,
    error: listError,
    isLoading: listIsLoading,
    mutate,
  } = GroupAPI.useList(
    { query, organizationId, page, size, sortOrder: GroupSortOrder.Desc },
    swrConfig(),
  )
  const { hasPageSwitcher } = usePageMonitor({
    totalPages: list?.totalPages ?? 1,
    totalElements: list?.totalElements ?? 0,
    steps: [size],
  })
  const listIsEmpty = list && !listError && list.totalElements === 0
  const listIsReady = list && !listError && list.totalElements > 0

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
            <div className={cx('flex', 'flex-col', 'gap-1.5')}>
              <SearchInput
                placeholder="Search Groups"
                query={query}
                onChange={handleSearchInputChange}
              />
              {listIsLoading ? <SectionSpinner /> : null}
              {listError ? (
                <SectionError text={errorToString(listError)} />
              ) : null}
              {listIsEmpty ? (
                <SectionPlaceholder text="There are no items." />
              ) : null}
              {listIsReady ? (
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
                      {list.data.map((g) => (
                        <Tr
                          key={g.id}
                          className={cx(
                            'cursor-pointer',
                            'h-[52px]',
                            { 'bg-gray-100': selected?.id === g.id },
                            { 'dark:bg-gray-600': selected?.id === g.id },
                            { 'bg-transparent': selected?.id !== g.id },
                          )}
                          onClick={() => setSelected(g)}
                        >
                          <Td className={cx('px-0.5', 'text-center')}>
                            <Radio
                              size="md"
                              isChecked={selected?.id === g.id}
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
                                name={g.name}
                                size="sm"
                                className={cx('w-[40px]', 'h-[40px]')}
                              />
                              <span className={cx('text-base')}>{g.name}</span>
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
                        size="sm"
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

export default GroupSelector
