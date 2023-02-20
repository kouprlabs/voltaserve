export default function mapFileList(fileList: FileList | null): File[] {
  if (!fileList || fileList.length === 0) {
    return []
  }
  const files = []
  for (let i = 0; i < fileList.length; i++) {
    files.push(fileList[i])
  }
  return files
}
