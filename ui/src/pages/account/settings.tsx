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
  useColorMode,
} from '@chakra-ui/react'
import { variables, IconEdit, IconTrash, SectionSpinner } from '@koupr/ui'
import { Helmet } from 'react-helmet-async'
import { swrConfig } from '@/api/options'
import StorageAPI from '@/api/storage'
import UserAPI from '@/api/user'
import AccountChangePassword from '@/components/account/change-password'
import AccountDelete from '@/components/account/delete'
import AccountEditEmail from '@/components/account/edit-email'
import AccountEditFullName from '@/components/account/edit-full-name'
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
          <Text fontWeight="bold">Storage usage</Text>
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
              <Text>Calculating???</Text>
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
            <Text>{user.email}</Text>
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
          <Text fontWeight="bold">Delete account</Text>
          <HStack spacing={variables.spacing} h={ROW_HEIGHT}>
            <Text>Delete permanently</Text>
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
        <AccountEditFullName
          open={isFullNameModalOpen}
          user={user}
          onClose={() => setIsFullNameModalOpen(false)}
        />
        <AccountEditEmail
          open={isEmailModalOpen}
          user={user}
          onClose={() => setIsEmailModalOpen(false)}
        />
        <AccountChangePassword
          open={isPasswordModalOpen}
          user={user}
          onClose={() => setIsPasswordModalOpen(false)}
        />
        <AccountDelete
          open={isDeleteModalOpen}
          onClose={() => setIsDeleteModalOpen(false)}
        />
      </Stack>
    </>
  )
}

export default AccountSettingsPage
