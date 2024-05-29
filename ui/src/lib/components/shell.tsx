import { MouseEvent, ReactElement } from 'react'
import cx from 'classnames'
import { StorageOptions } from '../types'
import Drawer, { DrawerItem } from './drawer'

type ShellItem = {
  href: string
  icon: ReactElement
  primaryText: string
  secondaryText: string
}

type ShellProps = {
  storage?: StorageOptions
  logo: ReactElement
  topBar: ReactElement
  items: ShellItem[]
  children?: ReactElement
  onContentClick?: (event: MouseEvent) => void
}

const Shell = ({
  logo,
  topBar,
  items,
  storage,
  children,
  onContentClick,
}: ShellProps) => (
  <div className={cx('flex', 'flex-row', 'items-center', 'gap-0', 'h-full')}>
    <Drawer storage={storage} logo={logo}>
      {items.map((item, index) => (
        <DrawerItem
          key={index}
          href={item.href}
          icon={item.icon}
          primaryText={item.primaryText}
          secondaryText={item.secondaryText}
        />
      ))}
    </Drawer>
    <div
      className={cx('flex', 'flex-col', 'items-center', 'h-full', 'w-full')}
      onClick={onContentClick}
    >
      {topBar}
      <div
        className={cx(
          'flex',
          'flex-col',
          'w-full',
          'lg:w-[1250px]',
          'px-2',
          'pt-2',
          'overflow-y-auto',
          'overflow-x-hidden',
          'grow',
        )}
      >
        {children}
      </div>
    </div>
  </div>
)

export default Shell
