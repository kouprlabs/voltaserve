// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useEffect, useState } from 'react'
import { useLocation, useNavigate, useSearchParams } from 'react-router-dom'
import {
  Button,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Stack,
  Table,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
} from '@chakra-ui/react'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import AdminApi from '@/client/admin/admin'
import { swrConfig } from '@/client/options'
import { organizationPaginationStorage } from '@/infra/pagination'
import PagePagination from '@/lib/components/page-pagination'
import SectionSpinner from '@/lib/components/section-spinner'
import { decodeQuery } from '@/lib/helpers/query'
import usePagePagination from '@/lib/hooks/page-pagination'
import { useAppDispatch } from '@/store/hook'
import { mutateUpdated } from '@/store/ui/indexes'

const AdminPanelDatabaseIndexes = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const dispatch = useAppDispatch()
  const [searchParams] = useSearchParams()
  const [confirmWindowOpen, setConfirmWindowOpen] = useState<boolean>(false)
  const [isSubmitting, setSubmitting] = useState<boolean>(false)
  const [focusedIndex, setFocusedIndex] = useState<string | null>(null)
  const query = decodeQuery(searchParams.get('q') as string)
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigate,
    location,
    storage: organizationPaginationStorage(),
  })
  const {
    data: list,
    error,
    mutate,
  } = AdminApi.useListIndexes({ page, size }, swrConfig())

  const sendRebuildRequest = (
    indexname: string | null,
    confirmation: boolean | undefined,
  ) => {
    if (confirmation) {
      setSubmitting(true)
      setTimeout(() => {
        setSubmitting(false)
        setConfirmWindowOpen(false)
        setFocusedIndex(null)
      }, 2000)
    } else if (indexname) {
      setFocusedIndex(indexname)
      setConfirmWindowOpen(true)
    } else {
      console.log('Fatal error while dispatching rebuild index', indexname)
    }
  }

  useEffect(() => {
    mutate()
  }, [query, page, size, mutate])

  useEffect(() => {
    if (mutate) {
      dispatch(mutateUpdated(mutate))
    }
  }, [mutate, dispatch])

  if (error) {
    return null
  }

  if (!list) {
    return <SectionSpinner />
  }

  return (
    <>
      <Modal
        isOpen={confirmWindowOpen}
        onClose={() => {
          setFocusedIndex(null)
          setConfirmWindowOpen(false)
          setSubmitting(false)
        }}
      >
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Are You sure?</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            You are going to rebuild {focusedIndex}, please confirm this action
          </ModalBody>
          <ModalFooter>
            <Button
              type="button"
              variant="outline"
              colorScheme="blue"
              disabled={isSubmitting}
              onClick={() => {
                setFocusedIndex(null)
                setConfirmWindowOpen(false)
                setSubmitting(false)
              }}
            >
              Cancel
            </Button>
            <Button
              type="button"
              variant="solid"
              colorScheme="blue"
              isLoading={isSubmitting}
              onClick={() => {
                sendRebuildRequest(focusedIndex, true)
              }}
            >
              Confirm
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
      <Helmet>
        <title>Indexes management</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        {list && list.data.length > 0 ? (
          <Stack direction="column" spacing={2}>
            <Table variant="simple">
              <Thead>
                <Tr>
                  <Th>Name</Th>
                  <Th>Table</Th>
                  <Th>Syntax</Th>
                </Tr>
              </Thead>
              <Tbody>
                {list.data.map((item) => (
                  <Tr key={item.indexname}>
                    <Td>{item.indexname}</Td>
                    <Td>{item.tablename}</Td>
                    <Td>{item.indexdef}</Td>
                    <Td>
                      <Button
                        colorScheme="green"
                        onClick={() =>
                          sendRebuildRequest(item.indexname, false)
                        }
                      >
                        Rebuild
                      </Button>
                    </Td>
                  </Tr>
                ))}
              </Tbody>
            </Table>
          </Stack>
        ) : (
          <div>No indexes to show</div>
        )}
        {list ? (
          <PagePagination
            style={{ alignSelf: 'end' }}
            totalElements={list.totalElements}
            totalPages={Math.ceil(list.totalElements / size)}
            page={page}
            size={size}
            steps={steps}
            setPage={setPage}
            setSize={setSize}
          />
        ) : null}
      </div>
    </>
  )
}

export default AdminPanelDatabaseIndexes
