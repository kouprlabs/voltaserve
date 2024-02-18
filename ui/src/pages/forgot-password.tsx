import { useCallback, useState } from 'react'
import { Link } from 'react-router-dom'
import {
  Button,
  FormControl,
  FormErrorMessage,
  Input,
  Link as ChakraLink,
  Text,
  Heading,
} from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import {
  Field,
  FieldAttributes,
  FieldProps,
  Form,
  Formik,
  FormikHelpers,
} from 'formik'
import * as Yup from 'yup'
import classNames from 'classnames'
import { Helmet } from 'react-helmet-async'
import AccountAPI from '@/client/idp/account'
import Logo from '@/components/common/logo'
import FullLayout from '@/components/layout/full'

type FormValues = {
  email: string
}

const ForgotPasswordPage = () => {
  const formSchema = Yup.object().shape({
    email: Yup.string()
      .email('Email is not valid')
      .required('Email is required'),
  })
  const [isCompleted, setIsCompleted] = useState(false)

  const handleSubmit = useCallback(
    async (
      { email }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>,
    ) => {
      try {
        await AccountAPI.sendResetPasswordEmail({
          email,
        })
        setIsCompleted(true)
      } finally {
        setSubmitting(false)
      }
    },
    [],
  )

  return (
    <FullLayout>
      <>
        <Helmet>
          <title>Forgot Password</title>
        </Helmet>
        <div
          className={classNames(
            'flex',
            'flex-col',
            'items-center',
            'gap-2.5',
            'w-full',
          )}
        >
          <Logo className="w-16" isGlossy={true} />
          <Heading fontSize={variables.headingFontSize}>
            Forgot Password
          </Heading>
          {isCompleted ? (
            <Text align="center">
              If your email belongs to an account, you will receive the recovery
              instructions in your inbox shortly.
            </Text>
          ) : (
            <>
              <Text align="center">
                Please provide your account Email where we can send you the
                password recovery instructions.
              </Text>
              <Formik
                initialValues={{
                  email: '',
                }}
                validationSchema={formSchema}
                validateOnBlur={false}
                onSubmit={handleSubmit}
              >
                {({ errors, touched, isSubmitting }) => (
                  <Form className="w-full">
                    <div
                      className={classNames(
                        'flex',
                        'flex-col',
                        'items-center',
                        'gap-1.5',
                      )}
                    >
                      <Field name="email">
                        {({ field }: FieldAttributes<FieldProps>) => (
                          <FormControl
                            isInvalid={
                              errors.email && touched.email ? true : false
                            }
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
                      <Button
                        variant="solid"
                        colorScheme="blue"
                        w="100%"
                        type="submit"
                        isLoading={isSubmitting}
                      >
                        Send Password Recovery Instructions
                      </Button>
                    </div>
                  </Form>
                )}
              </Formik>
              <div
                className={classNames(
                  'flex',
                  'flex-row',
                  'items-center',
                  'gap-0.5',
                )}
              >
                <Text>Password recovered?</Text>
                <ChakraLink as={Link} to="/sign-in">
                  Sign In
                </ChakraLink>
              </div>
            </>
          )}
        </div>
      </>
    </FullLayout>
  )
}

export default ForgotPasswordPage
