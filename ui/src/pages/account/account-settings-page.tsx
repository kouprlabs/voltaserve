import { useState } from 'react'
import { Switch, Tooltip, useColorMode } from '@chakra-ui/react'
import { IconButton, Text, Progress, Separator } from '@radix-ui/themes'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import { IoWarning } from 'react-icons/io5'
import StorageAPI from '@/client/api/storage'
import UserAPI from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import AccountChangePassword from '@/components/account/account-change-password'
import AccountDelete from '@/components/account/account-delete'
import AccountEditEmail from '@/components/account/account-edit-email'
import AccountEditFullName from '@/components/account/account-edit-full-name'
import prettyBytes from '@/helpers/pretty-bytes'
import { IconEdit, IconTrash, SectionSpinner } from '@/lib'

interface IconButtonProps extends React.ComponentProps<typeof IconButton> {}

const EditButton = ({ ...props }: IconButtonProps) => {
  return (
    <IconButton radius="full" {...props}>
      <IconEdit />
    </IconButton>
  )
}

const Spacer = () => <div className={cx('grow')} />

const AccountSettingsPage = () => {
  const { colorMode, toggleColorMode } = useColorMode()
  const { data: user, error: userError } = UserAPI.useGet()
  const { data: storageUsage, error: storageUsageError } =
    StorageAPI.useGetAccountUsage(swrConfig())
  const [isFullNameModalOpen, setIsFullNameModalOpen] = useState(false)
  const [isEmailModalOpen, setIsEmailModalOpen] = useState(false)
  const [isPasswordModalOpen, setIsPasswordModalOpen] = useState(false)
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false)
  const sectionClassName = cx('flex', 'flex-col', 'gap-1', 'py-1.5')
  const rowClassName = cx(
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
      <div className={cx('flex', 'flex-col', 'gap-0')}>
        <div className={sectionClassName}>
          <Text size="2" weight="bold">
            Storage Usage
          </Text>
          {storageUsageError && <span>Failed to load storage usage.</span>}
          {storageUsage && !storageUsageError && (
            <>
              <Text size="2">
                {prettyBytes(storageUsage.bytes)} of{' '}
                {prettyBytes(storageUsage.maxBytes)} used
              </Text>
              <Progress value={storageUsage.percentage} />
            </>
          )}
          {!storageUsage && !storageUsageError && (
            <>
              <Text size="2">Calculating…</Text>
              <Progress />
            </>
          )}
        </div>
        <Separator size="4" />
        <div className={sectionClassName}>
          <Text size="2" weight="bold">
            Basics
          </Text>
          <div className={cx(rowClassName)}>
            <Text size="2">Full name</Text>
            <Spacer />
            <Text size="2">{user.fullName}</Text>
            <EditButton
              aria-label=""
              onClick={() => {
                setIsFullNameModalOpen(true)
              }}
            />
          </div>
        </div>
        <Separator size="4" />
        <div className={sectionClassName}>
          <Text size="2" weight="bold">
            Credentials
          </Text>
          <div className={cx(rowClassName)}>
            <Text size="2">Email</Text>
            <Spacer />
            {user.pendingEmail && (
              <div className={cx('flex', 'flex-row', 'items-center')}>
                <Tooltip label="Please check your inbox to confirm your email.">
                  <span>
                    <IoWarning
                      className={cx('text-yelow-400', 'text-[20px]')}
                    />
                  </span>
                </Tooltip>
                <Text size="2">{user.pendingEmail}</Text>
              </div>
            )}
            {!user.pendingEmail && (
              <Text size="2">{user.pendingEmail || user.email}</Text>
            )}
            <EditButton
              aria-label=""
              onClick={() => {
                setIsEmailModalOpen(true)
              }}
            />
          </div>
          <div className={cx(rowClassName)}>
            <Text size="2">Password</Text>
            <Spacer />
            <EditButton
              aria-label=""
              onClick={() => {
                setIsPasswordModalOpen(true)
              }}
            />
          </div>
        </div>
        <Separator size="4" />
        <div className={sectionClassName}>
          <Text size="2" weight="bold">
            Theme
          </Text>
          <div className={cx(rowClassName)}>
            <Text size="2">Dark mode</Text>
            <Spacer />
            <Switch
              isChecked={colorMode === 'dark'}
              onChange={() => toggleColorMode()}
            />
          </div>
        </div>
        <Separator size="4" />
        <div className={sectionClassName}>
          <Text size="2" weight="bold">
            Advanced
          </Text>
          <div className={cx(rowClassName)}>
            <Text size="2">Delete account</Text>
            <Spacer />
            <IconButton
              variant="solid"
              color="crimson"
              radius="full"
              aria-label=""
              onClick={() => setIsDeleteModalOpen(true)}
            >
              <IconTrash />
            </IconButton>
          </div>
        </div>
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
      </div>
    </>
  )
}

export default AccountSettingsPage
