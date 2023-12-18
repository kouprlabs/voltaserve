import * as Yup from 'yup'

export default function parseEmailList(value: string): string[] {
  return [...new Set(value.split(',').map((e: string) => e.trim()))].filter(
    (e) => {
      if (e.length === 0) {
        return false
      }
      try {
        Yup.string()
          .email()
          .matches(
            /.+(\.[A-Za-z]{2,})$/,
            'Email must end with a valid top-level domain',
          )
          .validateSync(e)
        return true
      } catch {
        return false
      }
    },
  )
}
