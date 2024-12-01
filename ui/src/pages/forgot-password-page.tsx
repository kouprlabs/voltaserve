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
import { Button, FormControl, FormErrorMessage, Input, Link as ChakraLink, Heading } from '@chakra-ui/react'
import { Logo } from '@koupr/ui'
import { Field, FieldAttributes, FieldProps, Form, Formik, FormikHelpers } from 'formik'
import * as Yup from 'yup'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import AccountAPI from '@/client/idp/account'
import LayoutFull from '@/components/layout/layout-full'

type FormValues = {
  email: string
}

const ForgotPasswordPage = () => {
  const formSchema = Yup.object().shape({
    email: Yup.string().email('Email is not valid').required('Email is required'),
  })
  const [isCompleted, setIsCompleted] = useState(false)

  const handleSubmit = useCallback(async ({ email }: FormValues, { setSubmitting }: FormikHelpers<FormValues>) => {
    try {
      await AccountAPI.sendResetPasswordEmail({
        email,
      })
      setIsCompleted(true)
    } finally {
      setSubmitting(false)
    }
  }, [])

  return (
    <LayoutFull>
      <>
        <Helmet>
          <title>Forgot Password</title>
        </Helmet>
        <div className={cx('flex', 'flex-col', 'items-center', 'gap-2.5', 'w-full')}>
          <div className={cx('w-[64px]')}>
            <Logo type="voltaserve" size="md" isGlossy={true} />
          </div>
          <Heading className={cx('text-heading')}>Forgot Password</Heading>
          {isCompleted ? (
            <span className={cx('text-center')}>
              If your email belongs to an account, you will receive the recovery instructions in your inbox shortly.
            </span>
          ) : (
            <>
              <span className={cx('text-center')}>
                Please provide your account Email where we can send you the password recovery instructions.
              </span>
              <Formik
                initialValues={{
                  email: '',
                }}
                validationSchema={formSchema}
                validateOnBlur={false}
                onSubmit={handleSubmit}
              >
                {({ errors, touched, isSubmitting }) => (
                  <Form className={cx('w-full')}>
                    <div className={cx('flex', 'flex-col', 'items-center', 'gap-1.5')}>
                      <Field name="email">
                        {({ field }: FieldAttributes<FieldProps>) => (
                          <FormControl isInvalid={Boolean(errors.email && touched.email)}>
                            <Input {...field} id="email" placeholder="Email" disabled={isSubmitting} />
                            <FormErrorMessage>{errors.email}</FormErrorMessage>
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
                        Send Password Recovery Instructions
                      </Button>
                    </div>
                  </Form>
                )}
              </Formik>
              <div className={cx('flex', 'flex-row', 'items-center', 'gap-0.5')}>
                <span>Password recovered?</span>
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

export default ForgotPasswordPage
