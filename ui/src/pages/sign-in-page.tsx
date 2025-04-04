// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useCallback } from 'react'
import { Link, useNavigate } from 'react-router-dom'
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
import { TokenAPI } from '@/client/idp/token'
import LayoutFull from '@/components/layout/layout-full'
import { saveToken } from '@/infra/token'

type FormValues = {
  email: string
  password: string
}

const SignInPage = () => {
  const navigate = useNavigate()
  const formSchema = Yup.object().shape({
    email: Yup.string()
      .email('Email is not valid.')
      .required('Email is required.'),
    password: Yup.string().required('Password is required.'),
  })

  const handleSignIn = useCallback(
    async (
      { email: username, password }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>,
    ) => {
      try {
        const token = await TokenAPI.exchange({
          username,
          password,
          grant_type: 'password',
        })
        saveToken(token).then()
        navigate('/workspace')
      } catch (error) {
        console.error(error)
      } finally {
        setSubmitting(false)
      }
    },
    [navigate],
  )

  return (
    <LayoutFull>
      <>
        <Helmet>
          <title>Sign In</title>
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
          <Formik
            initialValues={{
              email: '',
              password: '',
            }}
            validationSchema={formSchema}
            validateOnBlur={false}
            onSubmit={handleSignIn}
          >
            {({ errors, touched, isSubmitting }) => (
              <Form className={cx('w-full')}>
                <div
                  className={cx('flex', 'flex-col', 'items-center', 'gap-1.5')}
                >
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
                        isInvalid={Boolean(errors.password && touched.password)}
                      >
                        <Input
                          {...field}
                          id="password"
                          placeholder="Password"
                          type="password"
                          disabled={isSubmitting}
                        />
                        <FormErrorMessage>{errors.password}</FormErrorMessage>
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
                    Sign In
                  </Button>
                </div>
              </Form>
            )}
          </Formik>
          <div
            className={cx('flex', 'flex-col', 'items-center', 'max-w-[60ch]')}
          >
            <div className={cx('flex', 'flex-row', 'items-center', 'gap-0.5')}>
              <span>{"Don't have an account yet?"}</span>
              <ChakraLink as={Link} to="/sign-up">
                Sign up
              </ChakraLink>
            </div>
            <div className={cx('flex', 'flex-row', 'items-center', 'gap-0.5')}>
              <span>Cannot sign in?</span>
              <ChakraLink as={Link} to="/forgot-password">
                Reset password
              </ChakraLink>
            </div>
          </div>
        </div>
      </>
    </LayoutFull>
  )
}

export default SignInPage
