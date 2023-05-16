import { useCallback } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import {
  Button,
  Container,
  FormControl,
  FormErrorMessage,
  HStack,
  Input,
  Link as ChakraLink,
  Text,
  VStack,
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
import { Helmet } from 'react-helmet-async'
import GroupAPI from '@/api/group'
import OrganizationAPI from '@/api/organization'
import TokenAPI from '@/api/token'
import WorkspaceAPI from '@/api/workspace'
import Logo from '@/components/common/logo'
import FullLayout from '@/components/layout/full'
import { saveToken } from '@/infra/token'
import { gigabyteToByte } from '@/helpers/convert-storage'

type FormValues = {
  email: string
  password: string
}

const SignInPage = () => {
  const navigate = useNavigate()
  const formSchema = Yup.object().shape({
    email: Yup.string()
      .email('Email is not valid')
      .required('Email is required'),
    password: Yup.string().required('Password is required'),
  })

  const handleSignIn = useCallback(
    async (
      { email: username, password }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>
    ) => {
      try {
        const token = await TokenAPI.exchange({
          username,
          password,
          grant_type: 'password',
        })
        saveToken(token)
        const { length: organizationCount } = await OrganizationAPI.getAll()
        if (organizationCount === 0) {
          const { id: organizationId } = await OrganizationAPI.create({
            name: 'My Organization',
          })
          await GroupAPI.create({
            name: 'My Group',
            organizationId,
          })
          const { id: workspaceId, rootId } = await WorkspaceAPI.create({
            name: 'My Workspace',
            organizationId,
            storageCapacity: gigabyteToByte(100),
          })
          navigate(`/workspace/${workspaceId}/file/${rootId}`)
        } else {
          const result = await WorkspaceAPI.getAll()
          if (result.length === 1) {
            navigate(`/workspace/${result[0].id}/file/${result[0].rootId}`)
          } else {
            navigate('/workspace')
          }
        }
      } finally {
        setSubmitting(false)
      }
    },
    [navigate]
  )

  return (
    <FullLayout>
      <>
        <Helmet>
          <title>Sign In to Voltaserve</title>
        </Helmet>
        <VStack spacing="25px" w="100%">
          <Logo className="w-16" isGlossy={true} />
          <h1 className="font-display text-2xl font-medium">
            Sign In to Voltaserve
          </h1>
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
                  <Button
                    variant="solid"
                    colorScheme="blue"
                    w="100%"
                    type="submit"
                    isLoading={isSubmitting}
                  >
                    Sign in
                  </Button>
                </VStack>
              </Form>
            )}
          </Formik>
          <Container centerContent>
            <HStack spacing={variables.spacingXs}>
              <Text>{"Don't have an account yet?"}</Text>
              <ChakraLink as={Link} to="/sign-up">
                Sign up
              </ChakraLink>
            </HStack>
            <HStack spacing={variables.spacingXs}>
              <Text>Cannot sign in?</Text>
              <ChakraLink as={Link} to="/forgot-password">
                Reset password
              </ChakraLink>
            </HStack>
          </Container>
        </VStack>
      </>
    </FullLayout>
  )
}

export default SignInPage
