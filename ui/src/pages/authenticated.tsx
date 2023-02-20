import { Outlet } from 'react-router-dom'
import DrawerLayout from '@/components/layout/drawer'

const AuthenticatedPage = () => {
  return (
    <DrawerLayout>
      <Outlet />
    </DrawerLayout>
  )
}

export default AuthenticatedPage
