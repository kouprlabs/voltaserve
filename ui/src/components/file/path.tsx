import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { Breadcrumb, BreadcrumbItem, BreadcrumbLink } from '@chakra-ui/react'
import FileAPI, { File } from '@/api/file'
import WorkspaceAPI from '@/api/workspace'

const FilePath = () => {
  const params = useParams()
  const workspaceId = params.id as string
  const fileId = params.fileId as string
  const { data: workspace } = WorkspaceAPI.useGetById(workspaceId)
  const [path, setPath] = useState<File[]>([])

  useEffect(() => {
    ;(async () => {
      setPath(await FileAPI.getPath(fileId))
    })()
  }, [fileId])

  return (
    <Breadcrumb>
      <BreadcrumbItem>
        {workspace && (
          <BreadcrumbLink
            as={Link}
            to={`/workspace/${workspaceId}/file/${workspace.rootId}`}
          >
            Home
          </BreadcrumbLink>
        )}
      </BreadcrumbItem>
      {path.slice(1).map((f) => (
        <BreadcrumbItem key={f.id}>
          <BreadcrumbLink
            as={Link}
            to={`/workspace/${workspaceId}/file/${f.id}`}
            isCurrentPage={fileId === f.id}
          >
            {f.name}
          </BreadcrumbLink>
        </BreadcrumbItem>
      ))}
    </Breadcrumb>
  )
}

export default FilePath
