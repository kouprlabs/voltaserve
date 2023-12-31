import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { Breadcrumb, BreadcrumbItem, BreadcrumbLink } from '@chakra-ui/react'
import FileAPI, { File } from '@/client/api/file'
import WorkspaceAPI from '@/client/api/workspace'

const Path = () => {
  const { id, fileId } = useParams()
  const { data: workspace } = WorkspaceAPI.useGetById(id)
  const [path, setPath] = useState<File[]>([])

  useEffect(() => {
    ;(async () => {
      setPath(await FileAPI.getPath(fileId!))
    })()
  }, [fileId])

  return (
    <Breadcrumb>
      <BreadcrumbItem>
        {workspace && (
          <BreadcrumbLink
            as={Link}
            to={`/workspace/${id}/file/${workspace.rootId}`}
          >
            Home
          </BreadcrumbLink>
        )}
      </BreadcrumbItem>
      {path.slice(1).map((f) => (
        <BreadcrumbItem key={f.id}>
          <BreadcrumbLink
            as={Link}
            to={`/workspace/${id}/file/${f.id}`}
            isCurrentPage={fileId === f.id}
          >
            {f.name}
          </BreadcrumbLink>
        </BreadcrumbItem>
      ))}
    </Breadcrumb>
  )
}

export default Path
