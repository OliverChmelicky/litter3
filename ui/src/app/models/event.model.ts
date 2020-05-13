export const permissionEventCreator = 'creator'
export const permissionEventEditor = 'editor'
export const permissionEventViewer = 'viewer'

export interface ChangePermisssionRequest {
  ChangingRightsTo: string,
  EventId: string,
  Permission: string,
  AsSociety: boolean
  SocietyId: string,
}

export interface EventModel {
  Id: string,
  Date: Date,
  Description: string,
  Publc: boolean,
  CreatedAt: Date,
  TrashIds?: string[],
  UsersIds?: string[],
  SocietiesIds?: string[],
}
