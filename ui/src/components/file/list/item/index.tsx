import { ChangeEvent, MouseEvent, useEffect } from 'react'
import { useCallback, useMemo, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import {
  Stack,
  Center,
  Link as ChakraLink,
  useColorModeValue,
  Checkbox,
  Box,
  Text,
  VStack,
} from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { File, SnapshotStatus } from '@/client/api/file'
import relativeDate from '@/helpers/relative-date'
import store from '@/store/configure-store'
import { useAppDispatch } from '@/store/hook'
import {
  selectedItemAdded,
  selectedItemRemoved,
  selectedItemsUpdated,
} from '@/store/ui/files'
import Icon from './icon'
import { performMultiSelect, performRangeSelect } from './perform-select'

type ItemProps = {
  file: File
  scale: number
  isPresentational?: boolean
  isLoading?: boolean
  onContextMenu?: (event: MouseEvent) => void
}

const WIDTH = 147
const MIN_HEIGHT = 110

const Item = ({
  file,
  scale,
  isPresentational,
  isLoading,
  onContextMenu,
}: ItemProps) => {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const width = useMemo(() => `${WIDTH * scale}px`, [scale])
  const minHeight = useMemo(() => `${MIN_HEIGHT * scale}px`, [scale])
  const hoverColor = useColorModeValue('gray.100', 'gray.700')
  const activeColor = useColorModeValue('gray.200', 'gray.600')
  const [isCheckboxVisible, setIsCheckboxVisible] = useState(false)
  const [isHovered, setIsHovered] = useState(false)
  const [isChecked, setIsChecked] = useState(false)
  const [isSelected, setIsSelected] = useState(false)
  const date = useMemo(
    () => relativeDate(new Date(file.createTime)),
    [file.createTime],
  )

  useEffect(() => {
    const unsubscribe = store.subscribe(() => {
      if (store.getState().ui.files.selectedItems.includes(file.id)) {
        setIsSelected(true)
        setIsChecked(true)
        setIsCheckboxVisible(true)
      } else {
        setIsSelected(false)
        setIsChecked(false)
        if (!isHovered) {
          setIsCheckboxVisible(false)
        }
      }
    })
    return () => unsubscribe()
  }, [file.id, isHovered])

  const handleClick = useCallback(
    (event: MouseEvent) => event.stopPropagation(),
    [],
  )

  const handleMouseEnter = useCallback(() => {
    setIsCheckboxVisible(true)
    setIsHovered(true)
  }, [])

  const handleMouseLeave = useCallback(() => {
    if (!isChecked) {
      setIsCheckboxVisible(false)
    }
    setIsHovered(false)
  }, [isChecked])

  const handleIconClick = useCallback(
    (event: MouseEvent) => {
      event.stopPropagation()
      if (store.getState().ui.files.isMultiSelectActive) {
        performMultiSelect(file, isSelected)
      } else if (store.getState().ui.files.isRangeSelectActive) {
        performRangeSelect(file)
      } else {
        dispatch(selectedItemsUpdated([file.id]))
      }
    },
    [file, isSelected, dispatch],
  )

  const handleFileLinkClick = useCallback(
    (event: MouseEvent) => {
      handleIconClick(event)
      window.open(`/file/${file.id}`, '_blank')?.focus()
    },
    [file.id, handleIconClick],
  )

  const handleIconDoubleClick = useCallback(() => {
    dispatch(selectedItemsUpdated([]))
    if (file.type === 'folder') {
      navigate(`/workspace/${file.workspaceId}/file/${file.id}`)
    } else if (file.type === 'file' && file.status === SnapshotStatus.Ready) {
      window.open(`/file/${file.id}`, '_blank')?.focus()
    }
  }, [file, navigate, dispatch])

  const handleCheckboxChange = useCallback(
    (event: ChangeEvent<HTMLInputElement>) => {
      event.stopPropagation()
      if (event.target.checked) {
        setIsChecked(true)
        dispatch(selectedItemAdded(file.id))
      } else {
        setIsChecked(false)
        dispatch(selectedItemRemoved(file.id))
      }
    },
    [file.id, dispatch],
  )

  const handleContextMenu = useCallback(
    (event: MouseEvent) => {
      if (event) {
        event.preventDefault()
        onContextMenu?.(event)
        if (!isSelected) {
          handleIconClick(event)
        }
      }
    },
    [isSelected, handleIconClick, onContextMenu],
  )

  return (
    <Stack
      position="relative"
      spacing={variables.spacingXs}
      w={width}
      py={variables.spacingSm}
      _hover={{ bg: hoverColor }}
      _active={{ bg: activeColor }}
      transition="background-color 0.4s ease"
      bg={isChecked ? hoverColor : 'transparent'}
      borderRadius={variables.borderRadiusSm}
      userSelect="none"
      cursor="default"
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
      onClick={handleClick}
      onContextMenu={handleContextMenu}
    >
      {(isCheckboxVisible || isSelected) && !isPresentational ? (
        <Checkbox
          position="absolute"
          top={variables.spacingSm}
          left={variables.spacingSm}
          isChecked={isChecked}
          zIndex={1}
          size="lg"
          onChange={handleCheckboxChange}
        />
      ) : null}
      <Center
        w={width}
        minH={minHeight}
        onDoubleClick={handleIconDoubleClick}
        onClick={handleIconClick}
      >
        <Icon file={file} scale={scale} isLoading={isLoading} />
      </Center>
      <Box w={width} title={file.name} px={variables.spacingXs}>
        {file.type === 'folder' && (
          <ChakraLink
            as={Link}
            to={`/workspace/${file.workspaceId}/file/${file.id}`}
            textAlign="center"
            noOfLines={3}
            textDecoration="none"
            _hover={{ textDecoration: 'underline' }}
          >
            {file.name}
          </ChakraLink>
        )}
        {file.type === 'file' && file.status === SnapshotStatus.Ready ? (
          <ChakraLink
            textAlign="center"
            noOfLines={3}
            textDecoration="none"
            _hover={{ textDecoration: 'underline' }}
            onClick={handleFileLinkClick}
          >
            {file.name}
          </ChakraLink>
        ) : null}
        {file.type === 'file' && file.status !== SnapshotStatus.Ready ? (
          <Text textAlign="center" noOfLines={3} onClick={handleIconClick}>
            {file.name}
          </Text>
        ) : null}
      </Box>
      <VStack spacing={0}>
        <Text textAlign="center" noOfLines={3} color="gray.500">
          {date}
        </Text>
      </VStack>
    </Stack>
  )
}

export default Item
