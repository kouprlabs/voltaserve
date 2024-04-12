import { useCallback, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import {
  Button,
  FormControl,
  FormErrorMessage,
  Input,
  Link as ChakraLink,
  Text,
  Heading,
} from '@chakra-ui/react'
import {
  Field,
  FieldAttributes,
  FieldProps,
  Form,
  Formik,
  FormikHelpers,
} from 'formik'
import * as Yup from 'yup'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import AccountAPI from '@/client/idp/account'
import Logo from '@/components/common/logo'
import LayoutFull from '@/components/layout/layout-full'

type FormValues = {
  newPassword: string
  newPasswordConfirmation: string
}

const ResetPasswordPage = () => {
  const params = useParams()
  const token = params.token as string
  const formSchema = Yup.object().shape({
    newPassword: Yup.string()
      .required('Password is required')
      .matches(
        /^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#$%^&*])(?=.{8,})/,
        'Must contain at least 8 characters, one Uppercase, one Lowercase, one number and one special character',
      ),
    newPasswordConfirmation: Yup.string()
      .oneOf([Yup.ref('newPassword'), undefined], 'Passwords do not match')
      .required('Confirm your password'),
  })
  const [isCompleted, setIsCompleted] = useState(false)

  const handleSubmit = useCallback(
    async (
      { newPassword }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>,
    ) => {
      try {
        await AccountAPI.resetPassword({
          newPassword,
          token: token,
        })
        setIsCompleted(true)
      } finally {
        setSubmitting(false)
      }
    },
    [token],
  )

  return (
    <LayoutFull>
      <>
        <Helmet>
          <title>Reset Password</title>
        </Helmet>
        <div
          className={cx(
            'flex',
            'flex-col',
            'items-center',
            'gap-2.5',
            'w-full',
          )}
        >
          <div className={cx('w-[64px]')}>
            <Logo isGlossy={true} />
          </div>
          <Heading className={cx('text-heading')}>Reset Password</Heading>
          {isCompleted ? (
            <div className={cx('flex', 'flex-row', 'items-center', 'gap-0.5')}>
              <Text align="center">Password successfully changed.</Text>
              <ChakraLink as={Link} to="/sign-in">
                Sign In
              </ChakraLink>
            </div>
          ) : (
            <>
              <Formik
                initialValues={{
                  newPassword: '',
                  newPasswordConfirmation: '',
                }}
                validationSchema={formSchema}
                validateOnBlur={false}
                onSubmit={handleSubmit}
              >
                {({ errors, touched, isSubmitting }) => (
                  <Form className={cx('w-full')}>
                    <div
                      className={cx(
                        'flex',
                        'flex-col',
                        'items-center',
                        'gap-1.5',
                      )}
                    >
                      <Field name="newPassword">
                        {({ field }: FieldAttributes<FieldProps>) => (
                          <FormControl
                            isInvalid={
                              errors.newPassword && touched.newPassword
                                ? true
                                : false
                            }
                          >
                            <Input
                              {...field}
                              id="newPassword"
                              placeholder="New password"
                              type="password"
                              disabled={isSubmitting}
                            />
                            <FormErrorMessage>
                              {errors.newPassword}
                            </FormErrorMessage>
                          </FormControl>
                        )}
                      </Field>
                      <Field name="newPasswordConfirmation">
                        {({ field }: FieldAttributes<FieldProps>) => (
                          <FormControl
                            isInvalid={
                              errors.newPasswordConfirmation &&
                              touched.newPasswordConfirmation
                                ? true
                                : false
                            }
                          >
                            <Input
                              {...field}
                              id="newPasswordConfirmation"
                              placeholder="Confirm new password"
                              type="password"
                              disabled={isSubmitting}
                            />
                            <FormErrorMessage>
                              {errors.newPasswordConfirmation}
                            </FormErrorMessage>
                          </FormControl>
                        )}
                      </Field>
                      <Button
                        className={cx('w-full')}
                        variant="solid"
                        colorScheme="blue"
                        type="submit"
                        isLoading={isSubmitting}
                      >
                        Reset Password
                      </Button>
                    </div>
                  </Form>
                )}
              </Formik>
              <div
                className={cx('flex', 'flex-row', 'items-center', 'gap-0.5')}
              >
                <Text>Password already reset?</Text>
                <ChakraLink as={Link} to="/sign-in">
                  Sign In
                </ChakraLink>
              </div>
            </>
          )}
        </div>
      </>
    </LayoutFull>
  )
}

export default ResetPasswordPage
