import useSWR, { SWRConfiguration } from 'swr'
import { apiFetcher } from '@/client/fetcher'

export type NotificationType = 'new_invitation'

export type Notification = {
  type: NotificationType
  body: unknown
}

export default class NotificationAPI {
  static useGetAll(swrOptions?: SWRConfiguration) {
    const url = `/notifications`
    return useSWR<Notification[]>(
      url,
      () => apiFetcher({ url, method: 'GET' }) as Promise<Notification[]>,
      swrOptions,
    )
  }
}
