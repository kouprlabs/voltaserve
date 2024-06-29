import {
  ChangeEvent,
  ReactElement,
  ReactNode,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from 'react'
import {
  IconButton,
  Popover,
  PopoverBody,
  PopoverContent,
  PopoverTrigger,
  Switch,
  Tooltip,
} from '@chakra-ui/react'
import cx from 'classnames'

type SwitchCardProps = {
  children?: ReactNode
  icon: ReactElement
  label: string
  isCollapsed?: boolean
  localStorageNamespace: string
  expandedMinWidth?: string
}

const SwitchCard = ({
  children,
  icon,
  label,
  isCollapsed,
  localStorageNamespace,
  expandedMinWidth,
}: SwitchCardProps) => {
  const [isActive, setIsActive] = useState(false)
  const localStorageActiveKey = useMemo(
    () => `voltaserve_${localStorageNamespace}_switch_card_active`,
    [localStorageNamespace],
  )

  useEffect(() => {
    let active = false
    if (typeof localStorage !== 'undefined') {
      const value = localStorage.getItem(localStorageActiveKey)
      if (value) {
        active = JSON.parse(value)
      } else {
        localStorage.setItem(localStorageActiveKey, JSON.stringify(false))
      }
    }
    if (active) {
      setIsActive(true)
    } else {
      setIsActive(false)
    }
  }, [localStorageActiveKey, setIsActive])

  const handleChange = useCallback(
    (event: ChangeEvent<HTMLInputElement>) => {
      setIsActive(event.target.checked)
      localStorage.setItem(
        localStorageActiveKey,
        JSON.stringify(event.target.checked),
      )
    },
    [localStorageActiveKey],
  )

  if (isCollapsed) {
    return (
      <Popover>
        <PopoverTrigger>
          <div>
            <Tooltip label={label}>
              <IconButton
                icon={icon}
                variant="outline"
                className={cx('w-[50px]', 'h-[50px]', 'p-1.5', 'rounded-md')}
                aria-label={label}
                title={label}
              />
            </Tooltip>
          </div>
        </PopoverTrigger>
        <PopoverContent>
          <PopoverBody>{children}</PopoverBody>
        </PopoverContent>
      </Popover>
    )
  } else {
    return (
      <div
        className={cx(
          'flex',
          'flex-col',
          'gap-0',
          'border',
          'border-gray-200',
          'dark:border-gray-600',
          'rounded-md',
        )}
        style={{ minWidth: expandedMinWidth }}
      >
        <div
          className={cx(
            'flex',
            'flex-row',
            'items-center',
            'gap-1',
            'h-[50px]',
            'px-1',
            'shrink-0',
          )}
        >
          {icon}
          <span className={cx('grow')}>{label}</span>
          <Switch isChecked={isActive} onChange={handleChange} />
        </div>
        {isActive ? (
          <div className={cx('pt-0', 'pr-1', 'pb-1', 'pl-1')}>{children}</div>
        ) : null}
      </div>
    )
  }
}

export default SwitchCard
