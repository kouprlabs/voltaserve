import { IconBaseProps } from 'react-icons'
import {
  BsFlag,
  BsCollection,
  BsBoxArrowUpRight,
  BsInfoCircle,
  BsCheck2,
  BsSearch,
  BsFillExclamationCircleFill,
  BsPlayFill,
  BsGridFill,
  BsSortDown,
  BsSortUp,
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
import { IoTimeOutline, IoCheckmarkCircle, IoRefresh } from 'react-icons/io5'
import { TbUsers } from 'react-icons/tb'
import { TbPlaylistX } from 'react-icons/tb'
import { VscClose } from 'react-icons/vsc'

const FONT_SIZE = '14px'
const FONT_SIZE_LG = '16px'

export const IconPlay = ({ fontSize, ...props }: IconBaseProps) => (
  <BsPlayFill fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconError = ({ fontSize, ...props }: IconBaseProps) => (
  <BsFillExclamationCircleFill fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconUpload = ({ fontSize, ...props }: IconBaseProps) => (
  <FiUpload fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconNotification = ({ fontSize, ...props }: IconBaseProps) => (
  <FiBell fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconDotsVertical = ({ fontSize, ...props }: IconBaseProps) => (
  <FiMoreVertical fontSize={fontSize || FONT_SIZE_LG} {...props} />
)

export const IconDotsVerticalSm = ({ fontSize, ...props }: IconBaseProps) => (
  <FiMoreVertical fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconExit = ({ fontSize, ...props }: IconBaseProps) => (
  <FiLogOut fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconChevronLeft = ({ fontSize, ...props }: IconBaseProps) => (
  <FiChevronLeft fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconChevronRight = ({ fontSize, ...props }: IconBaseProps) => (
  <FiChevronRight fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconAdd = ({ fontSize, ...props }: IconBaseProps) => (
  <FiPlus fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconEdit = ({ fontSize, ...props }: IconBaseProps) => (
  <FiEdit3 fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconShare = ({ fontSize, ...props }: IconBaseProps) => (
  <FiUsers fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconDownload = ({ fontSize, ...props }: IconBaseProps) => (
  <FiDownload fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconMove = ({ fontSize, ...props }: IconBaseProps) => (
  <FiCornerUpRight fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconCopy = ({ fontSize, ...props }: IconBaseProps) => (
  <FiCopy fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconTrash = ({ fontSize, ...props }: IconBaseProps) => (
  <FiTrash fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconSend = ({ fontSize, ...props }: IconBaseProps) => (
  <FiSend fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconUserPlus = ({ fontSize, ...props }: IconBaseProps) => (
  <FiUserPlus fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconCheck = ({ fontSize, ...props }: IconBaseProps) => (
  <BsCheck2 fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconCheckCircle = ({ fontSize, ...props }: IconBaseProps) => (
  <FiCheckCircle fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconAlertCircle = ({ fontSize, ...props }: IconBaseProps) => (
  <FiAlertCircle fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconCheckCircleFill = ({ fontSize, ...props }: IconBaseProps) => (
  <IoCheckmarkCircle fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconAlertCircleFill = ({ fontSize, ...props }: IconBaseProps) => (
  <IoMdAlert fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconCircle = ({ fontSize, ...props }: IconBaseProps) => (
  <FiCircle fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconWorkspace = ({ fontSize, ...props }: IconBaseProps) => (
  <BsCollection fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconGroup = ({ fontSize, ...props }: IconBaseProps) => (
  <TbUsers fontSize={fontSize || FONT_SIZE} {...props} />
)

export const IconOrganization = ({ fontSize, ...props }: IconBaseProps) => (
  <BsFlag fontSize={fontSize || FONT_SIZE} {...props} {...props} />
)

export const IconClose = ({ fontSize, ...props }: IconBaseProps) => (
  <VscClose fontSize={fontSize || FONT_SIZE} {...props} {...props} />
)

export const IconTime = ({ fontSize, ...props }: IconBaseProps) => (
  <IoTimeOutline fontSize={fontSize || FONT_SIZE} {...props} {...props} />
)

export const IconDeleteListItem = ({ fontSize, ...props }: IconBaseProps) => (
  <TbPlaylistX fontSize={fontSize || FONT_SIZE} {...props} {...props} />
)

export const IconExternalLink = ({ fontSize, ...props }: IconBaseProps) => (
  <BsBoxArrowUpRight fontSize={fontSize || FONT_SIZE} {...props} {...props} />
)

export const IconInfoCircle = ({ fontSize, ...props }: IconBaseProps) => (
  <BsInfoCircle fontSize={fontSize || FONT_SIZE} {...props} {...props} />
)

export const IconSearch = ({ fontSize, ...props }: IconBaseProps) => (
  <BsSearch fontSize={fontSize || FONT_SIZE} {...props} {...props} />
)

export const IconSearchBold = ({ fontSize, ...props }: IconBaseProps) => (
  <FaSearch fontSize={fontSize || FONT_SIZE} {...props} {...props} />
)

export const IconRefresh = ({ fontSize, ...props }: IconBaseProps) => (
  <IoRefresh fontSize={fontSize || FONT_SIZE_LG} {...props} {...props} />
)

export const IconGridFill = ({ fontSize, ...props }: IconBaseProps) => (
  <BsGridFill fontSize={fontSize || FONT_SIZE} {...props} {...props} />
)

export const IconSortDown = ({ fontSize, ...props }: IconBaseProps) => (
  <BsSortDown fontSize={fontSize || FONT_SIZE_LG} {...props} {...props} />
)

export const IconSortUp = ({ fontSize, ...props }: IconBaseProps) => (
  <BsSortUp fontSize={fontSize || FONT_SIZE_LG} {...props} {...props} />
)
