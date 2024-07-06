// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import { useState } from 'react'
import {
  Divider,
  IconButton,
  IconButtonProps,
  Progress,
  Switch,
  Tooltip,
  useColorMode,
} from '@chakra-ui/react'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import StorageAPI from '@/client/api/storage'
import UserAPI from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import AccountChangePassword from '@/components/account/account-change-password'
import AccountDelete from '@/components/account/account-delete'
import AccountEditEmail from '@/components/account/account-edit-email'
import AccountEditFullName from '@/components/account/account-edit-full-name'
import { IconEdit, IconDelete, IconWarning } from '@/lib/components/icons'
import SectionSpinner from '@/lib/components/section-spinner'
import prettyBytes from '@/lib/helpers/pretty-bytes'

const EditButton = (props: IconButtonProps) => (
  <IconButton
    icon={<IconEdit />}
    className={cx('h-[40px]', 'w-[40px]')}
    {...props}
  />
)

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
          <span className={cx('font-bold')}>Storage Usage</span>
          {storageUsageError && <span>Failed to load storage usage.</span>}
          {storageUsage && !storageUsageError ? (
            <>
              <span>
                {prettyBytes(storageUsage.bytes)} of{' '}
                {prettyBytes(storageUsage.maxBytes)} used
              </span>
              <Progress value={storageUsage.percentage} hasStripe />
            </>
          ) : null}
          {!storageUsage && !storageUsageError ? (
            <>
              <span>Calculating…</span>
              <Progress value={0} hasStripe />
            </>
          ) : null}
        </div>
        <Divider />
        <div className={sectionClassName}>
          <span className={cx('font-bold')}>Basics</span>
          <div className={cx(rowClassName)}>
            <span>Full name</span>
            <Spacer />
            <span>{user.fullName}</span>
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
          <span className={cx('font-bold')}>Credentials</span>
          <div className={cx(rowClassName)}>
            <span>Email</span>
            <Spacer />
            {user.pendingEmail ? (
              <div
                className={cx('flex', 'flex-row', 'gap-0.5', 'items-center')}
              >
                <Tooltip label="Please check your inbox to confirm your email.">
                  <div
                    className={cx(
                      'flex',
                      'items-center',
                      'justify-center',
                      'cursor-default',
                    )}
                  >
                    <IconWarning className={cx('text-yellow-400')} />
                  </div>
                </Tooltip>
                <span>{user.pendingEmail}</span>
              </div>
            ) : null}
            {!user.pendingEmail ? (
              <span>{user.pendingEmail || user.email}</span>
            ) : null}
            <EditButton
              aria-label=""
              onClick={() => {
                setIsEmailModalOpen(true)
              }}
            />
          </div>
          <div className={cx(rowClassName)}>
            <span>Password</span>
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
          <span className={cx('font-bold')}>Theme</span>
          <div className={cx(rowClassName)}>
            <span>Dark mode</span>
            <Spacer />
            <Switch
              isChecked={colorMode === 'dark'}
              onChange={() => toggleColorMode()}
            />
          </div>
        </div>
        <Divider />
        <div className={sectionClassName}>
          <span className={cx('font-bold')}>Advanced</span>
          <div className={cx(rowClassName)}>
            <span>Delete account</span>
            <Spacer />
            <IconButton
              icon={<IconDelete />}
              variant="solid"
              colorScheme="red"
              aria-label=""
              onClick={() => setIsDeleteModalOpen(true)}
            />
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
