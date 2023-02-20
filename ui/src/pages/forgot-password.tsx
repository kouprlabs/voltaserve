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
import Logo from '@/components/common/logo'
import FullLayout from '@/components/layout/full'
import variables from '@/theme/variables'

type FormValues = {
  email: string
}

const ForgotPasswordPage = () => {
  const toast = useToast()
  const [isFormDisaled, setIsFormDisabled] = useState(false)
  const formSchema = Yup.object().shape({
    email: Yup.string()
      .email('Email is not valid')
      .required('Email is required'),
  })

  const handleSubmit = useCallback(
    async (
      { email }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>
    ) => {
      try {
        await AccountAPI.sendResetPasswordEmail({
          email,
        })
        toast({
          title:
            'If your email belongs to an account, you will receive the recovery instructions in your inbox shortly',
          status: 'success',
          isClosable: true,
        })
        setIsFormDisabled(true)
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
          <title>Forgot Password</title>
        </Helmet>
        <VStack spacing="25px" w="100%">
          <Logo className="w-16" isGlossy={true} />
          <h1 className="font-display text-2xl font-medium">Forgot Password</h1>
          <Text align="center">
            Please provide your account Email where we can send you the password
            recovery instructions.
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
                <VStack spacing={variables.spacing}>
                  <Field name="email">
                    {({ field }: FieldAttributes<FieldProps>) => (
                      <FormControl
                        isInvalid={errors.email && touched.email ? true : false}
                      >
                        <Input
                          {...field}
                          id="email"
                          placeholder="Email"
                          disabled={isSubmitting || isFormDisaled}
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
                    disabled={isFormDisaled}
                    isLoading={isSubmitting}
                  >
                    Send pasword recovery instructions
                  </Button>
                </VStack>
              </Form>
            )}
          </Formik>
          <HStack spacing={variables.spacingXs}>
            <Text>Password recovered?</Text>
            <ChakraLink as={Link} to="/sign-in">
              Sign in
            </ChakraLink>
          </HStack>
        </VStack>
      </>
    </FullLayout>
  )
}

export default ForgotPasswordPage
