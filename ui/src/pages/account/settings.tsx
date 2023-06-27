import { useState } from 'react'
import {
  Box,
  Divider,
  HStack,
  IconButton,
  IconButtonProps,
  Progress,
  Stack,
  Switch,
  Text,
  Tooltip,
  useColorMode,
} from '@chakra-ui/react'
import { variables, IconEdit, IconTrash, SectionSpinner } from '@koupr/ui'
import { Helmet } from 'react-helmet-async'
import { IoWarning } from 'react-icons/io5'
import StorageAPI from '@/client/api/storage'
import UserAPI from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import ChangePassword from '@/components/account/change-password'
import Delete from '@/components/account/delete'
import EditEmail from '@/components/account/edit-email'
import EditFullName from '@/components/account/edit-full-name'
import prettyBytes from '@/helpers/pretty-bytes'

const EditButton = (props: IconButtonProps) => (
  <IconButton icon={<IconEdit />} w="40px" h="40px" {...props} />
)

const Spacer = () => <Box flexGrow={1} />

const ROW_HEIGHT = '40px'
const SECTION_SPACING = variables.spacing

const AccountSettingsPage = () => {
  const { colorMode, toggleColorMode } = useColorMode()
  const { data: user, error: userError } = UserAPI.useGet()
  const { data: storageUsage, error: storageUsageError } =
    StorageAPI.useGetAccountUsage(swrConfig())
  const [isFullNameModalOpen, setIsFullNameModalOpen] = useState(false)
  const [isEmailModalOpen, setIsEmailModalOpen] = useState(false)
  const [isPasswordModalOpen, setIsPasswordModalOpen] = useState(false)
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false)

  if (userError) {
    return null
  }
  if (!user) {
    return <SectionSpinner />
  }

  return (
    <>
      <Helmet>
        <title>{user.fullName}</title>
      </Helmet>
      <Stack direction="column" spacing={0}>
        <Stack direction="column" py={SECTION_SPACING}>
          <Text fontWeight="bold">Storage Usage</Text>
          {storageUsageError && <Text>Failed to load storage usage.</Text>}
          {storageUsage && !storageUsageError && (
            <>
              <Text>
                {prettyBytes(storageUsage.bytes)} of{' '}
                {prettyBytes(storageUsage.maxBytes)} used
              </Text>
              <Progress value={storageUsage.percentage} hasStripe />
            </>
          )}
          {!storageUsage && !storageUsageError && (
            <>
              <Text>Calculating…</Text>
              <Progress value={0} hasStripe />
            </>
          )}
        </Stack>
        <Divider />
        <Stack direction="column" py={SECTION_SPACING}>
          <Text fontWeight="bold">Basics</Text>
          <HStack spacing={variables.spacing} h={ROW_HEIGHT}>
            <Text>Full name</Text>
            <Spacer />
            <Text>{user.fullName}</Text>
            <EditButton
              aria-label=""
              onClick={() => {
                setIsFullNameModalOpen(true)
              }}
            />
          </HStack>
        </Stack>
        <Divider />
        <Stack direction="column" py={SECTION_SPACING}>
          <Text fontWeight="bold">Credentials</Text>
          <HStack spacing={variables.spacing} h={ROW_HEIGHT}>
            <Text>Email</Text>
            <Spacer />
            {user.pendingEmail && (
              <HStack>
                <Tooltip label="Please check your inbox to confirm your email.">
                  <Box>
                    <IoWarning fontSize="20px" color="gold" />
                  </Box>
                </Tooltip>
                <Text>{user.pendingEmail}</Text>
              </HStack>
            )}
            {!user.pendingEmail && (
              <Text>{user.pendingEmail || user.email}</Text>
            )}
            <EditButton
              aria-label=""
              onClick={() => {
                setIsEmailModalOpen(true)
              }}
            />
          </HStack>
          <HStack spacing={variables.spacing} h={ROW_HEIGHT}>
            <Text>Password</Text>
            <Spacer />
            <EditButton
              aria-label=""
              onClick={() => {
                setIsPasswordModalOpen(true)
              }}
            />
          </HStack>
        </Stack>
        <Divider />
        <Stack direction="column" py={SECTION_SPACING}>
          <Text fontWeight="bold">Theme</Text>
          <HStack spacing={variables.spacing} h={ROW_HEIGHT}>
            <Text>Dark mode</Text>
            <Spacer />
            <Switch
              isChecked={colorMode === 'dark'}
              onChange={() => toggleColorMode()}
            />
          </HStack>
        </Stack>
        <Divider />
        <Stack direction="column" py={SECTION_SPACING}>
          <Text fontWeight="bold">Advanced</Text>
          <HStack spacing={variables.spacing} h={ROW_HEIGHT}>
            <Text>Delete account</Text>
            <Spacer />
            <IconButton
              icon={<IconTrash />}
              variant="solid"
              colorScheme="red"
              aria-label=""
              onClick={() => setIsDeleteModalOpen(true)}
            />
          </HStack>
        </Stack>
        <EditFullName
          open={isFullNameModalOpen}
          user={user}
          onClose={() => setIsFullNameModalOpen(false)}
        />
        <EditEmail
          open={isEmailModalOpen}
          user={user}
          onClose={() => setIsEmailModalOpen(false)}
        />
        <ChangePassword
          open={isPasswordModalOpen}
          user={user}
          onClose={() => setIsPasswordModalOpen(false)}
        />
        <Delete
          open={isDeleteModalOpen}
          onClose={() => setIsDeleteModalOpen(false)}
        />
      </Stack>
    </>
  )
}

export default AccountSettingsPage
