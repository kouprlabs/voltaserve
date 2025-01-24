import { useCallback, useEffect, useMemo } from 'react'
import { useColorMode } from '@chakra-ui/react'
import { Select } from '@koupr/ui'
import { OptionBase, SingleValue } from 'chakra-react-select'
import cx from 'classnames'
import { useMedia } from 'react-use'

interface ThemeOption extends OptionBase {
  value: string
  label: string
}

const LOCAL_STORAGE_KEY = 'voltaserve-ui-color-mode'

const AccountThemeSwitcher = () => {
  const { colorMode, setColorMode } = useColorMode()
  const isSystemDark = useMedia('(prefers-color-scheme: dark)')
  const options = [
    { value: 'system', label: 'System' },
    { value: 'light', label: 'Light' },
    { value: 'dark', label: 'Dark' },
  ]
  const defaultValue = useMemo(() => {
    let option = options.find((e) => e.value === colorMode)
    if (localStorage.getItem(LOCAL_STORAGE_KEY) === 'system') {
      option = options[0]
    }
    return option
  }, [options])

  useEffect(() => {
    if (localStorage.getItem(LOCAL_STORAGE_KEY) === 'system') {
      setColorMode(isSystemDark ? 'dark' : 'light')
    }
  }, [isSystemDark])

  const handleChange = useCallback(
    (value: SingleValue<ThemeOption>) => {
      if (value!.value === 'system') {
        localStorage.setItem(LOCAL_STORAGE_KEY, 'system')
        setColorMode(isSystemDark ? 'dark' : 'light')
      } else {
        localStorage.setItem(LOCAL_STORAGE_KEY, value!.value)
        setColorMode(value!.value)
      }
    },
    [isSystemDark, setColorMode],
  )

  return (
    <Select<ThemeOption, false>
      className={cx('min-w-[150px]')}
      defaultValue={defaultValue}
      options={options}
      selectedOptionStyle="check"
      onChange={handleChange}
    />
  )
}

export default AccountThemeSwitcher
