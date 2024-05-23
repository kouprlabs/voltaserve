import { extendTheme } from '@chakra-ui/react'
import breakpoints from './breakpoints'
import colors from './colors'
import Breadcrumb from './components/breadcrumb'
import Button from './components/button'
import Card from './components/card'
import Checkbox from './components/checkbox'
import Heading from './components/heading'
import Input from './components/input'
import Link from './components/link'
import Menu from './components/menu'
import Modal from './components/modal'
import Popover from './components/popover'
import Progress from './components/progress'
import Select from './components/select'
import Tabs from './components/tabs'
import Textarea from './components/textarea'
import Tooltip from './components/tooltip'
import styles from './styles'
import typography from './typography'

const overrides = {
  breakpoints,
  styles,
  colors,
  ...typography,
  components: {
    Button,
    Heading,
    Checkbox,
    Select,
    Input,
    Textarea,
    Modal,
    Link,
    Progress,
    Tabs,
    Tooltip,
    Popover,
    Breadcrumb,
    Menu,
    Card,
  },
}

export const theme = extendTheme(overrides)
