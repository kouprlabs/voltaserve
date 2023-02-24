import { useCallback, useState } from 'react'
import { Link } from 'react-router-dom'
import {
  Field,
  FieldAttributes,
  FieldProps,
  Form,
  Formik,
  FormikHelpers,
} from 'formik'
import * as Yup from 'yup'
import {
  Button,
  FormControl,
  FormErrorMessage,
  HStack,
  Input,
  Link as ChakraLink,
  Text,
  useToast,
  VStack,
} from '@chakra-ui/react'
import { Helmet } from 'react-helmet-async'
import AccountAPI from '@/api/account'
import { errorToString } from '@/api/error'
import Logo from '@/components/common/logo'
import FullLayout from '@/components/layout/full'
import variables from '@/theme/variables'

type FormValues = {
  fullName: string
  email: string
  password: string
  passwordConfirmation: string
}

const SignUpPage = () => {
  const toast = useToast()
  const [isConfirmationVisible, setIsConfirmationVisible] = useState(false)
  const formSchema = Yup.object().shape({
    fullName: Yup.string().required('Name is required'),
    email: Yup.string()
      .email('Email is not valid')
      .required('Email is required'),
    password: Yup.string().required('Password is required'),
    passwordConfirmation: Yup.string()
      .oneOf([Yup.ref('password'), undefined], 'Passwords must match')
      .required('Confirm your password'),
  })

  const handleSubmit = useCallback(
    async (
      { fullName, email, password }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>
    ) => {
      try {
        await AccountAPI.create({
          fullName,
          email,
          password,
        })
        setIsConfirmationVisible(true)
      } catch (e) {
        toast({
          title: errorToString(e),
          status: 'error',
          isClosable: true,
        })
      } finally {
        setSubmitting(false)
      }
    },
    [toast]
  )

  return (
    <FullLayout>
      <>
        <Helmet>
          <title>Sign Up to Voltaserve</title>
        </Helmet>
        {isConfirmationVisible && (
          <VStack spacing="25px" w="100%">
            <VStack spacing={variables.spacing}>
              <Logo className="w-16" isGlossy={true} />
              <h1 className="font-display text-2xl font-medium text-center">
                Thanks! We just sent you a confirmation email
              </h1>
              <Text align="center">
                Just open your inbox, find the email, and click on the
                confirmation link.
              </Text>
            </VStack>
          </VStack>
        )}
        {!isConfirmationVisible && (
          <VStack spacing="25px" w="100%">
            <Logo className="w-16" isGlossy={true} />
            <h1 className="font-display text-2xl font-medium">
              Sign Up To Voltaserve
            </h1>
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
              {({ errors, touched, isSubmitting }) => (
                <Form className="w-full">
                  <VStack spacing={variables.spacing}>
                    <Field name="fullName">
                      {({ field }: FieldAttributes<FieldProps>) => (
                        <FormControl
                          isInvalid={
                            errors.fullName && touched.fullName ? true : false
                          }
                        >
                          <Input
                            {...field}
                            id="fullName"
                            placeholder="Full name"
                            disabled={isSubmitting}
                          />
                          <FormErrorMessage>{errors.fullName}</FormErrorMessage>
                        </FormControl>
                      )}
                    </Field>
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
                    <Field name="password">
                      {({ field }: FieldAttributes<FieldProps>) => (
                        <FormControl
                          isInvalid={
                            errors.password && touched.password ? true : false
                          }
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
                    <Field name="passwordConfirmation">
                      {({ field }: FieldAttributes<FieldProps>) => (
                        <FormControl
                          isInvalid={
                            errors.passwordConfirmation &&
                            touched.passwordConfirmation
                              ? true
                              : false
                          }
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
                      variant="solid"
                      colorScheme="blue"
                      w="100%"
                      type="submit"
                      isLoading={isSubmitting}
                    >
                      Sign up
                    </Button>
                  </VStack>
                </Form>
              )}
            </Formik>
            <HStack spacing={variables.spacingXs}>
              <Text>Already a member?</Text>
              <ChakraLink as={Link} to="/sign-in">
                Sign in
              </ChakraLink>
            </HStack>
          </VStack>
        )}
      </>
    </FullLayout>
  )
}

export default SignUpPage
