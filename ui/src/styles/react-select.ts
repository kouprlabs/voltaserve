import { SystemStyleObject } from '@chakra-ui/system'

const reactSelectStyles = {
  dropdownIndicator: (provided: SystemStyleObject) => ({
    ...provided,
    bg: 'transparent',
    cursor: 'inherit',
    position: 'absolute',
    right: '0px',
  }),
  indicatorSeparator: (provided: SystemStyleObject) => ({
    ...provided,
    display: 'none',
  }),
  placeholder: (provided: SystemStyleObject) => ({
    ...provided,
    textAlign: 'center',
  }),
  singleValue: (provided: SystemStyleObject) => ({
    ...provided,
    textAlign: 'center',
  }),
}

export default reactSelectStyles
