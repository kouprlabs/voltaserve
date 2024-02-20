import { useState } from 'react'
import {
  Divider,
  IconButton,
  IconButtonProps,
  Progress,
  Switch,
  Text,
  Tooltip,
  useColorMode,
} from '@chakra-ui/react'
import { IconEdit, IconTrash, SectionSpinner } from '@koupr/ui'
import classNames from 'classnames'
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

const Spacer = () => <div className={classNames('grow')} />

const AccountSettingsPage = () => {
  const { colorMode, toggleColorMode } = useColorMode()
  const { data: user, error: userError } = UserAPI.useGet()
  const { data: storageUsage, error: storageUsageError } =
    StorageAPI.useGetAccountUsage(swrConfig())
  const [isFullNameModalOpen, setIsFullNameModalOpen] = useState(false)
  const [isEmailModalOpen, setIsEmailModalOpen] = useState(false)
  const [isPasswordModalOpen, setIsPasswordModalOpen] = useState(false)
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false)
  const sectionClassName = classNames('flex', 'flex-col', 'gap-1', 'py-1.5')
  const rowClassName = classNames(
    'flex',
    'flex-row',
    'items-center',
    'gap-1',
    `h-[40px]`,
  )

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
      <div className={classNames('flex', 'flex-col', 'gap-0')}>
        <div className={sectionClassName}>
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
        </div>
        <Divider />
        <div className={sectionClassName}>
          <Text fontWeight="bold">Basics</Text>
          <div className={classNames(rowClassName)}>
            <Text>Full name</Text>
            <Spacer />
            <Text>{user.fullName}</Text>
            <EditButton
              aria-label=""
              onClick={() => {
                setIsFullNameModalOpen(true)
              }}
            />
          </div>
        </div>
        <Divider />
        <div className={sectionClassName}>
          <Text fontWeight="bold">Credentials</Text>
          <div className={classNames(rowClassName)}>
            <Text>Email</Text>
            <Spacer />
            {user.pendingEmail && (
              <div className={classNames('flex', 'flex-row', 'items-center')}>
                <Tooltip label="Please check your inbox to confirm your email.">
                  <span>
                    <IoWarning fontSize="20px" color="gold" />
                  </span>
                </Tooltip>
                <Text>{user.pendingEmail}</Text>
              </div>
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
          </div>
          <div className={classNames(rowClassName)}>
            <Text>Password</Text>
            <Spacer />
            <EditButton
              aria-label=""
              onClick={() => {
                setIsPasswordModalOpen(true)
              }}
            />
          </div>
        </div>
        <Divider />
        <div className={sectionClassName}>
          <Text fontWeight="bold">Theme</Text>
          <div className={classNames(rowClassName)}>
            <Text>Dark mode</Text>
            <Spacer />
            <Switch
              isChecked={colorMode === 'dark'}
              onChange={() => toggleColorMode()}
            />
          </div>
        </div>
        <Divider />
        <div className={sectionClassName}>
          <Text fontWeight="bold">Advanced</Text>
          <div className={classNames(rowClassName)}>
            <Text>Delete account</Text>
            <Spacer />
            <IconButton
              icon={<IconTrash />}
              variant="solid"
              colorScheme="red"
              aria-label=""
              onClick={() => setIsDeleteModalOpen(true)}
            />
          </div>
        </div>
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
      </div>
    </>
  )
}

export default AccountSettingsPage
