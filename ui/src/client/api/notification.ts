/* eslint-disable react-hooks/rules-of-hooks */
import useSWR from 'swr'
import { apiFetcher } from '@/client/fetch'
import { Invitation } from './invitation'

export type NotificationType = 'new_invitation'

export type Notification = {
  type: NotificationType
  body: Invitation | any
}

export default class NotificationAPI {
  static useGetAll(swrOptions?: any) {
    return useSWR<Notification[]>(`/notifications`, apiFetcher, swrOptions)
  }
}
