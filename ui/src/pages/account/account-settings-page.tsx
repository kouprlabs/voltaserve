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
  IconButton,
  IconButtonProps,
  Progress,
  Switch,
  Tooltip,
  useColorMode,
} from '@chakra-ui/react'
import {
  IconEdit,
  IconDelete,
  IconWarning,
  SectionSpinner,
  Form,
} from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import StorageAPI from '@/client/api/storage'
import UserAPI from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import AccountChangePassword from '@/components/account/account-change-password'
import AccountDelete from '@/components/account/account-delete'
import AccountEditEmail from '@/components/account/account-edit-email'
import AccountEditFullName from '@/components/account/account-edit-full-name'
import prettyBytes from '@/lib/helpers/pretty-bytes'
import { truncateEnd } from '@/lib/helpers/truncate-end'
import truncateMiddle from '@/lib/helpers/truncate-middle'

const EditButton = (props: IconButtonProps) => (
  <IconButton
    icon={<IconEdit />}
    className={cx('h-[40px]', 'w-[40px]')}
    {...props}
  />
)

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
      <Form
        sections={[
          {
            title: 'Storage',
            content: (
              <>
                {storageUsageError ? (
                  <span>Failed to load storage usage.</span>
                ) : null}
                {storageUsage && !storageUsageError ? (
                  <>
                    <span>
                      {prettyBytes(storageUsage.bytes)} of{' '}
                      {prettyBytes(storageUsage.maxBytes)} used
                    </span>
                    <Progress value={storageUsage.percentage} hasStripe />
                  </>
                ) : null}
                {!(storageUsage && !storageUsageError) ? (
                  <>
                    <span>Calculating…</span>
                    <Progress value={0} hasStripe />
                  </>
                ) : null}
              </>
            ),
          },
          {
            title: 'Basics',
            rows: [
              {
                label: 'Full name',
                content: (
                  <>
                    <span>{truncateEnd(user.fullName, 50)}</span>
                    <EditButton
                      aria-label="Edit name"
                      onClick={() => {
                        setIsFullNameModalOpen(true)
                      }}
                    />
                  </>
                ),
              },
            ],
          },
          {
            title: 'Credentials',
            rows: [
              {
                label: 'Email',
                content: (
                  <>
                    {user.pendingEmail ? (
                      <div
                        className={cx(
                          'flex',
                          'flex-row',
                          'gap-0.5',
                          'items-center',
                        )}
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
                        <span>{truncateMiddle(user.pendingEmail, 50)}</span>
                      </div>
                    ) : null}
                    {!user.pendingEmail ? (
                      <span>
                        {truncateMiddle(user.pendingEmail || user.email, 50)}
                      </span>
                    ) : null}
                    <EditButton
                      aria-label="Edit email"
                      onClick={() => {
                        setIsEmailModalOpen(true)
                      }}
                    />
                  </>
                ),
              },
              {
                label: 'Password',
                content: (
                  <EditButton
                    aria-label="Change password"
                    onClick={() => {
                      setIsPasswordModalOpen(true)
                    }}
                  />
                ),
              },
            ],
          },
          {
            title: 'Theme',
            rows: [
              {
                label: 'Dark mode',
                content: (
                  <Switch
                    isChecked={colorMode === 'dark'}
                    onChange={() => toggleColorMode()}
                  />
                ),
              },
            ],
          },
          {
            title: 'Advanced',
            rows: [
              {
                label: 'Delete account',
                content: (
                  <IconButton
                    icon={<IconDelete />}
                    variant="solid"
                    colorScheme="red"
                    aria-label="Delete account"
                    onClick={() => setIsDeleteModalOpen(true)}
                  />
                ),
              },
            ],
          },
        ]}
      />
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
    </>
  )
}

export default AccountSettingsPage
