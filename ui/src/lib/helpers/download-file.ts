import { File } from '@/client/api/file'
import { getAccessTokenOrRedirect } from '@/infra/token'

export default async function downloadFile(file: File) {
  if (!file.snapshot?.original || file.type !== 'file') {
    return
  }
  const a: HTMLAnchorElement = document.createElement('a')
  a.href = `/proxy/api/v2/files/${file.id}/original${
    file.snapshot?.original.extension
  }?${new URLSearchParams({
    access_token: getAccessTokenOrRedirect(),
    download: 'true',
  })}`
  a.download = file.name
  a.style.display = 'none'
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
}
