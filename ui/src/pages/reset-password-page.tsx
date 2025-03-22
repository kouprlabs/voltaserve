// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useCallback, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import {
  Button,
  FormControl,
  FormErrorMessage,
  Input,
  Link as ChakraLink,
} from '@chakra-ui/react'
import { Logo } from '@koupr/ui'
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
import { AccountAPI } from '@/client/idp/account'
import LayoutFull from '@/components/layout/layout-full'
import PasswordHints from '@/components/sign-up/password-hints'
import { YupFactory } from '@/lib/validation'

type FormValues = {
  newPassword: string
  newPasswordConfirmation: string
}

const ResetPasswordPage = () => {
  const params = useParams()
  const token = params.token as string
  const formSchema = Yup.object().shape({
    newPassword: YupFactory.password('New password'),
    newPasswordConfirmation: Yup.string()
      .oneOf([Yup.ref('newPassword'), undefined], 'Passwords do not match.')
      .required('Confirm your password.'),
  })
  const [isCompleted, setIsCompleted] = useState(false)
  const { data: passwordRequirements } = AccountAPI.useGetPasswordRequirements()

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
            <Logo type="voltaserve" size="md" isGlossy={true} />
          </div>
          {isCompleted ? (
            <div className={cx('flex', 'flex-col', 'items-center', 'gap-0.5')}>
              <span className={cx('text-center')}>
                Password successfully changed.
              </span>
              <ChakraLink as={Link} to="/sign-in">
                Sign in
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
                {({ errors, touched, isSubmitting, values }) => (
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
                            isInvalid={Boolean(
                              errors.newPassword && touched.newPassword,
                            )}
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
                            {passwordRequirements ? (
                              <div className="pt-1">
                                <PasswordHints
                                  value={values.newPassword}
                                  requirements={passwordRequirements}
                                />
                              </div>
                            ) : null}
                          </FormControl>
                        )}
                      </Field>
                      <Field name="newPasswordConfirmation">
                        {({ field }: FieldAttributes<FieldProps>) => (
                          <FormControl
                            isInvalid={Boolean(
                              errors.newPasswordConfirmation &&
                                touched.newPasswordConfirmation,
                            )}
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
                <span>Password already reset?</span>
                <ChakraLink as={Link} to="/sign-in">
                  Sign in
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
