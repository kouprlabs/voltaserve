// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { ReactNode, useEffect, useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import cx from 'classnames'
import { StorageOptions } from '../../types'
import { IconChevronLeft, IconChevronRight } from '../icons'
import { DrawerContext } from './drawer-context'

type DrawerProps = {
  children?: ReactNode
  logo?: ReactNode
  storage?: StorageOptions
}

const Drawer = ({ children, storage, logo }: DrawerProps) => {
  const [isCollapsed, setIsCollapsed] = useState<boolean | undefined>(undefined)
  const [isTouched, setIsTouched] = useState(false)
  const localStorageCollapsedKey = useMemo(
    () =>
      `${storage?.prefix || 'app'}_${
        storage?.namespace || 'main'
      }_drawer_collapsed`,
    [storage],
  )

  useEffect(() => {
    let collapse = false
    if (typeof localStorage !== 'undefined') {
      const value = localStorage.getItem(localStorageCollapsedKey)
      if (value) {
        collapse = JSON.parse(value)
      } else {
        localStorage.setItem(localStorageCollapsedKey, JSON.stringify(true))
      }
    }
    setIsCollapsed(collapse)
  }, [localStorageCollapsedKey, setIsCollapsed])

  if (isCollapsed === undefined) {
    return null
  }

  return (
    <DrawerContext.Provider
      value={{
        isCollapsed,
        isTouched,
      }}
    >
      <div
        className={cx(
          'flex',
          'flex-col',
          'h-full',
          'border-r',
          'border-r-gray-200',
          'dark:border-r-gray-700',
          'shrink-0',
          'gap-0',
        )}
      >
        <div
          className={cx('flex', 'items-center', 'justify-center', 'h-[80px]')}
        >
          <Link
            to={
              location.pathname.startsWith('/console') ? '/console/dashboard' : '/'
            }
          >
            <div className={cx('flex', 'h-[40px]')}>
              <div
                className={cx(
                  'flex',
                  'items-center',
                  'justify-center',
                  'w-[40px]',
                  'h-[40px]',
                )}
              >
                {logo}
              </div>
            </div>
          </Link>
        </div>
        <div
          className={cx(
            'flex',
            'flex-col',
            'items-center',
            'gap-0.5',
            'pt-0',
            'pr-1.5',
            'pb-1.5',
            'pl-1.5',
          )}
        >
          {children}
        </div>
        <div className={cx('grow')} />
        <div
          className={cx(
            'flex',
            'flex-row',
            'items-center',
            'gap-0',
            { 'justify-center': isCollapsed, 'justify-end': !isCollapsed },
            'h-[50px]',
            'w-full',
            { 'px-0': isCollapsed, 'px-1.5': !isCollapsed },
            'cursor-pointer',
            'hover:bg-gray-100',
            'hover:dark:bg-gray-600',
            'active:bg-gray-200',
            'active:dark:bg-gray-700',
          )}
          onClick={() => {
            setIsCollapsed(!isCollapsed)
            setIsTouched(true)
            localStorage.setItem(
              localStorageCollapsedKey,
              JSON.stringify(!isCollapsed),
            )
          }}
        >
          {isCollapsed ? <IconChevronRight /> : <IconChevronLeft />}
        </div>
      </div>
    </DrawerContext.Provider>
  )
}

export default Drawer
export { DrawerContext } from './drawer-context'
export { DrawerItem } from './drawer-item'
