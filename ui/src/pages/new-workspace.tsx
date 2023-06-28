import { useCallback, useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import {
  Button,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Heading,
  HStack,
  Input,
  Select,
  Stack,
} from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { useSWRConfig } from 'swr'
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
import OrganizationAPI, { Organization } from '@/client/api/organization'
import { geEditorPermission } from '@/client/api/permission'
import WorkspaceAPI from '@/client/api/workspace'
import OrganizationSelector from '@/components/common/organization-selector'
import StorageInput from '@/components/common/storage-input'
import { gigabyteToByte } from '@/helpers/convert-storage'

type FormValues = {
  name: string
  organizationId: string
  storageCapacity: number
}

const NewWorkspacePage = () => {
  const navigate = useNavigate()
  const { mutate } = useSWRConfig()
  const [orgs, setOrgs] = useState<Organization[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const formSchema = Yup.object().shape({
    name: Yup.string().required('Name is required').max(255),
    organizationId: Yup.string().required('Organization is required'),
    storageCapacity: Yup.number()
      .required('Storage capacity is required')
      .positive()
      .integer()
      .min(1, 'Invalid storage usage value'),
  })

  useEffect(() => {
    ;(async () => {
      const list = await OrganizationAPI.list()
      setOrgs(list.data.filter((o) => geEditorPermission(o.permission)))
    })()
  }, [])

  const handleSubmit = useCallback(
    async (
      { name, organizationId, storageCapacity }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>
    ) => {
      setSubmitting(true)
      setIsLoading(true)
      try {
        const result = await WorkspaceAPI.create({
          name,
          organizationId,
          storageCapacity,
        })
        mutate(`/workspaces/${result.id}`, result)
        mutate(`/workspaces`)
        setSubmitting(false)
        navigate(`/workspace/${result.id}/file/${result.rootId}`)
      } catch (e) {
        setIsLoading(false)
      } finally {
        setSubmitting(false)
      }
    },
    [navigate, mutate]
  )

  return (
    <>
      <Helmet>
        <title>New Workspace</title>
      </Helmet>
      <Stack spacing={variables.spacing2Xl}>
        <Heading size="lg">New Workspace</Heading>
        <Formik
          enableReinitialize={true}
          initialValues={{
            name: '',
            organizationId: '',
            storageCapacity: gigabyteToByte(100),
          }}
          validationSchema={formSchema}
          validateOnBlur={false}
          onSubmit={handleSubmit}
        >
          {({ errors, touched, isSubmitting, setFieldValue }) => (
            <Form>
              <Stack spacing={variables.spacing2Xl}>
                <Stack spacing={variables.spacing}>
                  <Field name="name">
                    {({ field }: FieldAttributes<FieldProps>) => (
                      <FormControl
                        maxW="400px"
                        isInvalid={errors.name && touched.name ? true : false}
                      >
                        <FormLabel>Name</FormLabel>
                        <Input {...field} disabled={isSubmitting} autoFocus />
                        <FormErrorMessage>{errors.name}</FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                  <Field name="organizationId">
                    {({ field }: FieldAttributes<FieldProps>) => (
                      <FormControl
                        maxW="400px"
                        isInvalid={
                          errors.organizationId && touched.organizationId
                            ? true
                            : false
                        }
                      >
                        <FormLabel>Organization</FormLabel>
                        <OrganizationSelector
                          onConfirm={(value) =>
                            setFieldValue(field.name, value.id)
                          }
                        />
                        <FormErrorMessage>
                          {errors.organizationId}
                        </FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                  <Field name="storageCapacity">
                    {(props: FieldAttributes<FieldProps>) => (
                      <FormControl
                        maxW="400px"
                        isInvalid={
                          errors.storageCapacity && touched.storageCapacity
                            ? true
                            : false
                        }
                      >
                        <FormLabel>Storage capacity</FormLabel>
                        <StorageInput {...props} />
                        <FormErrorMessage>
                          {errors.storageCapacity}
                        </FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                </Stack>
                <HStack>
                  <Button
                    type="submit"
                    variant="solid"
                    colorScheme="blue"
                    isDisabled={isSubmitting || isLoading}
                    isLoading={isSubmitting || isLoading}
                  >
                    Create
                  </Button>
                  <Button as={Link} to="/workspace" variant="solid">
                    Cancel
                  </Button>
                </HStack>
              </Stack>
            </Form>
          )}
        </Formik>
      </Stack>
    </>
  )
}

export default NewWorkspacePage
