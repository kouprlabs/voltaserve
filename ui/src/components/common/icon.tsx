import { IconBaseProps } from 'react-icons'
import {
  BsFlag,
  BsCollection,
  BsBoxArrowUpRight,
  BsInfoCircle,
  BsCheck2,
  BsSearch,
} from 'react-icons/bs'
import { FaSearch } from 'react-icons/fa'
import {
  FiBell,
  FiCheckCircle,
  FiChevronLeft,
  FiChevronRight,
  FiCircle,
  FiCopy,
  FiEdit3,
  FiLogOut,
  FiMoreVertical,
  FiCornerUpRight,
  FiPlus,
  FiSend,
  FiUsers,
  FiTrash,
  FiUpload,
  FiUserPlus,
  FiAlertCircle,
  FiDownload,
} from 'react-icons/fi'
import { IoMdAlert } from 'react-icons/io'
import { IoTimeOutline, IoCheckmarkCircle } from 'react-icons/io5'
import { TbUsers } from 'react-icons/tb'
import { TbPlaylistX } from 'react-icons/tb'
import { VscClose } from 'react-icons/vsc'

const DEFAULT_FONT_SIZE = '14px'

export const IconUpload = ({ fontSize, ...props }: IconBaseProps) => (
  <FiUpload fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconNotification = ({ fontSize, ...props }: IconBaseProps) => (
  <FiBell fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconDotsVertical = ({ fontSize, ...props }: IconBaseProps) => (
  <FiMoreVertical fontSize={fontSize || '16px'} {...props} />
)

export const IconDotsVerticalSm = ({ fontSize, ...props }: IconBaseProps) => (
  <FiMoreVertical fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconExit = ({ fontSize, ...props }: IconBaseProps) => (
  <FiLogOut fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconChevronLeft = ({ fontSize, ...props }: IconBaseProps) => (
  <FiChevronLeft fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconChevronRight = ({ fontSize, ...props }: IconBaseProps) => (
  <FiChevronRight fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconAdd = ({ fontSize, ...props }: IconBaseProps) => (
  <FiPlus fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconEdit = ({ fontSize, ...props }: IconBaseProps) => (
  <FiEdit3 fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconShare = ({ fontSize, ...props }: IconBaseProps) => (
  <FiUsers fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconDownload = ({ fontSize, ...props }: IconBaseProps) => (
  <FiDownload fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconMove = ({ fontSize, ...props }: IconBaseProps) => (
  <FiCornerUpRight fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconCopy = ({ fontSize, ...props }: IconBaseProps) => (
  <FiCopy fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconTrash = ({ fontSize, ...props }: IconBaseProps) => (
  <FiTrash fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconSend = ({ fontSize, ...props }: IconBaseProps) => (
  <FiSend fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconUserPlus = ({ fontSize, ...props }: IconBaseProps) => (
  <FiUserPlus fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconCheck = ({ fontSize, ...props }: IconBaseProps) => (
  <BsCheck2 fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconCheckCircle = ({ fontSize, ...props }: IconBaseProps) => (
  <FiCheckCircle fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconAlertCircle = ({ fontSize, ...props }: IconBaseProps) => (
  <FiAlertCircle fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconCheckCircleFill = ({ fontSize, ...props }: IconBaseProps) => (
  <IoCheckmarkCircle fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconAlertCircleFill = ({ fontSize, ...props }: IconBaseProps) => (
  <IoMdAlert fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconCircle = ({ fontSize, ...props }: IconBaseProps) => (
  <FiCircle fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconWorkspace = ({ fontSize, ...props }: IconBaseProps) => (
  <BsCollection fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconGroup = ({ fontSize, ...props }: IconBaseProps) => (
  <TbUsers fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} />
)

export const IconOrganization = ({ fontSize, ...props }: IconBaseProps) => (
  <BsFlag fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} {...props} />
)

export const IconClose = ({ fontSize, ...props }: IconBaseProps) => (
  <VscClose fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} {...props} />
)

export const IconTime = ({ fontSize, ...props }: IconBaseProps) => (
  <IoTimeOutline
    fontSize={fontSize || DEFAULT_FONT_SIZE}
    {...props}
    {...props}
  />
)

export const IconDeleteListItem = ({ fontSize, ...props }: IconBaseProps) => (
  <TbPlaylistX fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} {...props} />
)

export const IconExternalLink = ({ fontSize, ...props }: IconBaseProps) => (
  <BsBoxArrowUpRight
    fontSize={fontSize || DEFAULT_FONT_SIZE}
    {...props}
    {...props}
  />
)

export const IconInfoCircle = ({ fontSize, ...props }: IconBaseProps) => (
  <BsInfoCircle
    fontSize={fontSize || DEFAULT_FONT_SIZE}
    {...props}
    {...props}
  />
)

export const IconSearch = ({ fontSize, ...props }: IconBaseProps) => (
  <BsSearch fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} {...props} />
)

export const IconSearchBold = ({ fontSize, ...props }: IconBaseProps) => (
  <FaSearch fontSize={fontSize || DEFAULT_FONT_SIZE} {...props} {...props} />
)
