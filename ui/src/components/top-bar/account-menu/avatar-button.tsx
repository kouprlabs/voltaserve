import { Avatar, useColorModeValue, useToken } from '@chakra-ui/react'
import { forwardRef } from '@chakra-ui/system'
import classNames from 'classnames'
import { User } from '@/client/idp/user'
import { useAppSelector } from '@/store/hook'
import { NavType } from '@/store/ui/nav'
import ActiveCircle from './active-circle'

type AvatarButtonProps = {
  user: User
}

const AvatarButton = forwardRef<AvatarButtonProps, 'div'>(
  ({ user, ...props }, ref) => {
    const borderColor = useColorModeValue('gray.300', 'gray.700')
    const [borderColorDecoded] = useToken('colors', [borderColor])
    const activeNav = useAppSelector((state) => state.ui.nav.active)
    return (
      <div ref={ref} {...props} className={classNames('cursor-pointer')}>
        <ActiveCircle>
          <Avatar
            name={user.fullName}
            src={user.picture}
            size="sm"
            width="40px"
            height="40px"
            border={
              activeNav === NavType.Account
                ? 'none'
                : `1px solid ${borderColorDecoded}`
            }
          />
        </ActiveCircle>
      </div>
    )
  },
)

export default AvatarButton
