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
  Stack,
  Box,
  Text,
  Switch,
  HStack,
} from '@chakra-ui/react'
import { variables } from '@koupr/ui'

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
    [localStorageNamespace]
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
        JSON.stringify(event.target.checked)
      )
    },
    [localStorageActiveKey]
  )

  if (isCollapsed) {
    return (
      <Popover>
        <PopoverTrigger>
          <IconButton
            icon={icon}
            variant="outline"
            w="50px"
            h="50px"
            p={variables.spacing}
            borderRadius={variables.borderRadiusSm}
            aria-label={label}
            title={label}
          />
        </PopoverTrigger>
        <PopoverContent>
          <PopoverBody>{children}</PopoverBody>
        </PopoverContent>
      </Popover>
    )
  } else {
    return (
      <Stack
        border="1px solid"
        borderColor="gray.200"
        borderRadius={variables.borderRadiusSm}
        minW={expandedMinWidth}
        spacing={0}
      >
        <HStack
          direction="row"
          spacing={variables.spacingSm}
          h="50px"
          px={variables.spacing}
          flexShrink={0}
        >
          {icon}
          <Text flexGrow={1}>{label}</Text>
          <Switch isChecked={isActive} onChange={handleChange} />
        </HStack>
        {isActive && (
          <Box
            p={`0 ${variables.spacingSm} ${variables.spacingSm} ${variables.spacingSm}`}
          >
            {children}
          </Box>
        )}
      </Stack>
    )
  }
}

export default SwitchCard
