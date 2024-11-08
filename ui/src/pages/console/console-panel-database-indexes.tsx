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
  Code,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
} from '@chakra-ui/react'
import {
  DataTable,
  PagePagination,
  SectionSpinner,
  Text,
  usePagePagination,
} from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import ConsoleAPI from '@/client/console/console'
import { swrConfig } from '@/client/options'
import { organizationPaginationStorage } from '@/infra/pagination'
import { decodeQuery } from '@/lib/helpers/query'
import store from '@/store/configure-store'
import { errorOccurred } from '@/store/ui/error'

const ConsolePanelDatabaseIndexes = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const [searchParams] = useSearchParams()
  const [confirmWindowOpen, setConfirmWindowOpen] = useState<boolean>(false)
  const [isSubmitting, setSubmitting] = useState<boolean>(false)
  const [focusedIndex, setFocusedIndex] = useState<string | null>(null)
  const query = decodeQuery(searchParams.get('q') as string)
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigateFn: navigate,
    searchFn: () => location.search,
    storage: organizationPaginationStorage(),
  })
  const {
    data: list,
    error,
    mutate,
  } = ConsoleAPI.useListIndexes({ page, size }, swrConfig())

  const sendRebuildRequest = (
    indexName: string | null,
    confirmation: boolean | undefined,
  ) => {
    if (confirmation) {
      setSubmitting(true)
      setTimeout(() => {
        setSubmitting(false)
        setConfirmWindowOpen(false)
        setFocusedIndex(null)
      }, 2000)
    } else if (indexName) {
      setFocusedIndex(indexName)
      setConfirmWindowOpen(true)
    } else {
      const message = `Fatal error while dispatching rebuild index ${indexName}`
      store.dispatch(errorOccurred(message))
      console.error(message)
    }
  }

  useEffect(() => {
    mutate().then()
  }, [query, page, size, mutate])

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
          <ModalHeader>Rebuild Index</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            Are you sure you want to rebuild {focusedIndex}?
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
        <title>Index Management</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        {list && list.data.length > 0 ? (
          <DataTable
            items={list.data}
            columns={[
              {
                title: 'Name',
                renderCell: (item) => <Text>{item.indexName}</Text>,
              },
              {
                title: 'Table',
                renderCell: (item) => <Text>{item.tableName}</Text>,
              },
              {
                title: 'Syntax',
                renderCell: (item) => <Code>{item.indexDef}</Code>,
              },
            ]}
          />
        ) : (
          <div>No indexes to show</div>
        )}
        {list ? (
          <div className={cx('self-end')}>
            <PagePagination
              totalElements={list.totalElements}
              totalPages={Math.ceil(list.totalElements / size)}
              page={page}
              size={size}
              steps={steps}
              setPage={setPage}
              setSize={setSize}
            />
          </div>
        ) : null}
      </div>
    </>
  )
}

export default ConsolePanelDatabaseIndexes
