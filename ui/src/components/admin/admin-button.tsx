// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useRef, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  Button,
  IconButton,
  Input,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
} from '@chakra-ui/react'
import AdminApi from '@/client/admin/admin'
import { saveAdminToken } from '@/infra/admin-token'
import { IconAdmin } from '@/lib/components/icons'

const AdminButton = () => {
  const navigate = useNavigate()
  const buttonRef = useRef<HTMLButtonElement>(null)
  const [adminAuthOpen, setAdminAuthOpen] = useState<boolean>(false)
  const [isSubmitting, setSubmitting] = useState<boolean>(false)
  const [adminToken, setAdminToken] = useState<string>('')
  const authenticateAdmin = (token: string) => {
    saveAdminToken(token)
    AdminApi.adminAuthenticate().then((response) => {
      if (response) {
        setAdminAuthOpen(false)
        setAdminToken('')
        setSubmitting(false)
        setTimeout(() => navigate('/admin/dashboard'), 1000)
      } else {
        setAdminAuthOpen(false)
        setAdminToken('')
        setSubmitting(false)
      }
    })
  }
  return (
    <>
      <Modal
        isOpen={adminAuthOpen}
        onClose={() => {
          setAdminAuthOpen(false)
          setAdminToken('')
          setSubmitting(false)
        }}
      >
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Authorize yourself</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Input
              type="password"
              placeholder="Admin auth token"
              disabled={isSubmitting}
              autoFocus
              value={adminToken}
              onChange={(event) => setAdminToken(event.target.value)}
              autoComplete="adminToken"
            />
          </ModalBody>
          <ModalFooter>
            <Button
              type="button"
              variant="outline"
              isDisabled={isSubmitting}
              onClick={() => {
                setAdminAuthOpen(false)
                setAdminToken('')
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
                authenticateAdmin(adminToken)
                setSubmitting(true)
              }}
            >
              Go to Admin Panel
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
      <IconButton
        ref={buttonRef}
        icon={<IconAdmin />}
        aria-label=""
        onClick={() => {
          setAdminAuthOpen(true)
        }}
      />
    </>
  )
}

export default AdminButton
