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
import { Link } from 'react-router-dom'
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

type FormValues = {
  fullName: string
  email: string
  password: string
  passwordConfirmation: string
}

const SignUpPage = () => {
  const [isConfirmationVisible, setIsConfirmationVisible] = useState(false)
  const formSchema = Yup.object().shape({
    fullName: Yup.string().required('Name is required.'),
    email: Yup.string()
      .email('Email is not valid.')
      .required('Email is required.'),
    password: Yup.string().required('Password is required.'),
    passwordConfirmation: Yup.string()
      .oneOf([Yup.ref('password'), undefined], 'Passwords must match.')
      .required('Confirm your password.'),
  })
  const { data: passwordRequirements } = AccountAPI.useGetPasswordRequirements()

  const handleSubmit = useCallback(
    async (
      { fullName, email, password }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>,
    ) => {
      try {
        await AccountAPI.create({
          fullName,
          email,
          password,
        })
        setIsConfirmationVisible(true)
      } finally {
        setSubmitting(false)
      }
    },
    [],
  )

  return (
    <LayoutFull>
      <>
        <Helmet>
          <title>Sign Up</title>
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
          {isConfirmationVisible ? (
            <span className={cx('text-center')}>
              Thanks! We just sent you a confirmation email. Just open your
              inbox, find the email, and click on the confirmation link.
            </span>
          ) : null}
          {!isConfirmationVisible ? (
            <>
              <Formik
                initialValues={{
                  fullName: '',
                  email: '',
                  password: '',
                  passwordConfirmation: '',
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
                      <Field name="fullName">
                        {({ field }: FieldAttributes<FieldProps>) => (
                          <FormControl
                            isInvalid={Boolean(
                              errors.fullName && touched.fullName,
                            )}
                          >
                            <Input
                              {...field}
                              id="fullName"
                              placeholder="Full name"
                              disabled={isSubmitting}
                            />
                            <FormErrorMessage>
                              {errors.fullName}
                            </FormErrorMessage>
                          </FormControl>
                        )}
                      </Field>
                      <Field name="email">
                        {({ field }: FieldAttributes<FieldProps>) => (
                          <FormControl
                            isInvalid={Boolean(errors.email && touched.email)}
                          >
                            <Input
                              {...field}
                              id="email"
                              placeholder="Email"
                              disabled={isSubmitting}
                            />
                            <FormErrorMessage>{errors.email}</FormErrorMessage>
                          </FormControl>
                        )}
                      </Field>
                      <Field name="password">
                        {({ field }: FieldAttributes<FieldProps>) => (
                          <FormControl
                            isInvalid={Boolean(
                              errors.password && touched.password,
                            )}
                          >
                            <Input
                              {...field}
                              id="password"
                              placeholder="Password"
                              type="password"
                              disabled={isSubmitting}
                            />
                            <FormErrorMessage>
                              {errors.password}
                            </FormErrorMessage>
                            {passwordRequirements ? (
                              <div className="pt-1">
                                <PasswordHints
                                  value={values.password}
                                  requirements={passwordRequirements}
                                />
                              </div>
                            ) : null}
                          </FormControl>
                        )}
                      </Field>
                      <Field name="passwordConfirmation">
                        {({ field }: FieldAttributes<FieldProps>) => (
                          <FormControl
                            isInvalid={Boolean(
                              errors.passwordConfirmation &&
                                touched.passwordConfirmation,
                            )}
                          >
                            <Input
                              {...field}
                              id="passwordConfirmation"
                              placeholder="Confirm password"
                              type="password"
                              disabled={isSubmitting}
                            />
                            <FormErrorMessage>
                              {errors.passwordConfirmation}
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
                        Sign Up
                      </Button>
                    </div>
                  </Form>
                )}
              </Formik>
              <div
                className={cx('flex', 'flex-row', 'items-center', 'gap-0.5')}
              >
                <span>Already a member?</span>
                <ChakraLink as={Link} to="/sign-in">
                  Sign in
                </ChakraLink>
              </div>
            </>
          ) : null}
        </div>
      </>
    </LayoutFull>
  )
}

export default SignUpPage
