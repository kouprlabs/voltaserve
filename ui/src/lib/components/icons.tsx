// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { HTMLAttributes } from 'react'
import { cx } from '@emotion/css'

export type IconBaseProps = {
  filled?: boolean
} & HTMLAttributes<HTMLSpanElement>

type GetClassNameOptions = {
  filled?: boolean
  className?: string
}

export function getClassName({ filled, className }: GetClassNameOptions) {
  return cx(
    'material-symbols-rounded',
    { 'material-symbols-rounded__filled': filled },
    'text-[16px]',
    className,
  )
}

export const IconPlayArrow = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    play_arrow
  </span>
)

export const IconUpload = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    upload
  </span>
)

export const IconAdmin = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    admin_panel_settings
  </span>
)

export const IconDatabase = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    database
  </span>
)

export const IconRemoveOperator = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    remove_moderator
  </span>
)

export const IconNotifications = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    notifications
  </span>
)

export const IconMoreVert = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    more_vert
  </span>
)

export const IconLogout = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    logout
  </span>
)

export const IconChevronLeft = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    chevron_left
  </span>
)

export const IconChevronRight = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    chevron_right
  </span>
)

export const IconChevronDown = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    keyboard_arrow_down
  </span>
)

export const IconChevronUp = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    keyboard_arrow_up
  </span>
)

export const IconAdd = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    add
  </span>
)

export const IconEdit = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    edit
  </span>
)

export const IconGroup = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    group
  </span>
)

export const IconDownload = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    download
  </span>
)

export const IconArrowTopRight = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    arrow_top_right
  </span>
)

export const IconFileCopy = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    file_copy
  </span>
)

export const IconDelete = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    delete
  </span>
)

export const IconSend = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    send
  </span>
)

export const IconPersonAdd = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    person_add
  </span>
)

export const IconPerson = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    person
  </span>
)

export const IconCheck = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    check
  </span>
)

export const IconLibraryAddCheck = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    library_add_check
  </span>
)

export const IconSelectCheckBox = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    select_check_box
  </span>
)

export const IconCheckBoxOutlineBlank = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    check_box_outline_blank
  </span>
)

export const IconCheckCircle = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    check_circle
  </span>
)

export const IconError = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    error
  </span>
)

export const IconWarning = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    warning
  </span>
)

export const IconInvitations = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    forward_to_inbox
  </span>
)

export const IconWorkspaces = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    workspaces
  </span>
)

export const IconFlag = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    flag
  </span>
)

export const IconClose = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    close
  </span>
)

export const IconSchedule = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    schedule
  </span>
)

export const IconClearAll = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    clear_all
  </span>
)

export const IconOpenInNew = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    open_in_new
  </span>
)

export const IconInfo = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    info
  </span>
)

export const IconSearch = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    search
  </span>
)

export const IconRefresh = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    refresh
  </span>
)

export const IconSync = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    sync
  </span>
)

export const IconGridView = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    grid_view
  </span>
)

export const IconArrowUpward = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    arrow_upward
  </span>
)

export const IconArrowDownward = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    arrow_downward
  </span>
)

export const IconExpandMore = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    expand_more
  </span>
)

export const IconList = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    list
  </span>
)

export const IconHourglass = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    hourglass
  </span>
)

export const IconKeyboardArrowLeft = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    keyboard_arrow_left
  </span>
)

export const IconKeyboardArrowRight = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    keyboard_arrow_right
  </span>
)

export const IconKeyboardDoubleArrowRight = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    keyboard_double_arrow_right
  </span>
)

export const IconKeyboardDoubleArrowLeft = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    keyboard_double_arrow_left
  </span>
)

export const IconFirstPage = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    first_page
  </span>
)

export const IconLastPage = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    last_page
  </span>
)

export const IconHistory = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    history
  </span>
)

export const IconModeHeat = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    mode_heat
  </span>
)

export const IconSecurity = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    security
  </span>
)

export const IconVisibility = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    visibility
  </span>
)

export const IconTune = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    tune
  </span>
)

export const IconHome = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    home
  </span>
)

export const IconStacks = ({ className, filled, ...props }: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    stacks
  </span>
)

export const IconCloudUpload = ({
  className,
  filled,
  ...props
}: IconBaseProps) => (
  <span className={getClassName({ filled, className })} {...props}>
    cloud_upload
  </span>
)
