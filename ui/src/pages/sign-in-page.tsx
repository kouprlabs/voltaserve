import { useCallback } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import {
  Button,
  FormControl,
  FormErrorMessage,
  Input,
  Link as ChakraLink,
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
import GroupAPI from '@/client/api/group'
import OrganizationAPI from '@/client/api/organization'
import WorkspaceAPI from '@/client/api/workspace'
import TokenAPI from '@/client/idp/token'
import Logo from '@/components/common/logo'
import LayoutFull from '@/components/layout/layout-full'
import { saveToken } from '@/infra/token'
import { gigabyteToByte } from '@/lib/helpers/convert-storage'

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
      { setSubmitting }: FormikHelpers<FormValues>,
    ) => {
      try {
        const token = await TokenAPI.exchange({
          username,
          password,
          grant_type: 'password',
        })
        saveToken(token)
        const orgList = await OrganizationAPI.list()
        if (orgList.totalElements === 0) {
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
          const workspaceList = await WorkspaceAPI.list()
          if (workspaceList.totalElements === 1) {
            navigate(
              `/workspace/${workspaceList.data[0].id}/file/${workspaceList.data[0].rootId}`,
            )
          } else {
            navigate('/workspace')
          }
        }
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
          <title>Sign In to Voltaserve</title>
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
          <Heading className={cx('text-heading')}>
            Sign In to Voltaserve
          </Heading>
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
                Sign Up
              </ChakraLink>
            </div>
            <div className={cx('flex', 'flex-row', 'items-center', 'gap-0.5')}>
              <span>Cannot sign in?</span>
              <ChakraLink as={Link} to="/forgot-password">
                Reset Password
              </ChakraLink>
            </div>
          </div>
        </div>
      </>
    </LayoutFull>
  )
}

export default SignInPage
