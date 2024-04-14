import { ReactElement, useContext, useEffect, useState } from 'react'
import { Link, useLocation } from 'react-router-dom'
import cx from 'classnames'
import { DrawerContext } from './drawer-context'

type DrawerItemProps = {
  icon: ReactElement
  href: string
  primaryText: string
  secondaryText: string
  isActive?: boolean
}

export const DrawerItem = ({
  icon,
  href,
  primaryText,
  secondaryText,
}: DrawerItemProps) => {
  const location = useLocation()
  const [isActive, setIsActive] = useState<boolean>()
  const { isCollapsed } = useContext(DrawerContext)

  useEffect(() => {
    if (
      (href === '/' && location.pathname === '/') ||
      (href !== '/' && location.pathname.startsWith(href))
    ) {
      setIsActive(true)
    } else {
      setIsActive(false)
    }
  }, [location.pathname, href])

  return (
    <Link
      to={href}
      title={isCollapsed ? `${primaryText}: ${secondaryText}` : secondaryText}
      className={cx('w-full')}
    >
      <div
        className={cx(
          'flex',
          'flex-row',
          'items-center',
          'gap-1.5',
          'p-1.5',
          'h-[50px]',
          'rounded-md',
          {
            'bg-black': isActive,
            'dark:bg-white': isActive,
          },
          {
            'hover:bg-gray-100': !isActive,
            'dark:hover:bg-gray-600': !isActive,
          },
          {
            'hover:bg-gray-200': !isActive,
            'dark:hover:bg-gray-700': !isActive,
          },
        )}
      >
        <div
          className={cx(
            'flex',
            'items-center',
            'justify-center',
            'shrink-0',
            'w-[24px]',
            'h-[24px]',
            {
              'text-white': isActive,
              'dark:text-gray-800': isActive,
            },
          )}
        >
          {icon}
        </div>
        {!isCollapsed && (
          <span
            className={cx({
              'text-white': isActive,
              'dark:text-gray-800': isActive,
            })}
          >
            {primaryText}
          </span>
        )}
      </div>
    </Link>
  )
}
