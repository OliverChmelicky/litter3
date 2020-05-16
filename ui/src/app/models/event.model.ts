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
  Id?: string,
  Date: Date,
  Description: string,
  CreatedAt?: Date,
  TrashIds?: string[],
  UsersIds?: string[],
  SocietiesIds?: string[],
}

export interface EventCreatorModel {
  AsSociety: boolean,
  Id: string,
  VisibleName: string
}

export interface EventRequest {
  Id?:          string
  UserId?:      string
  SocietyId?:   string
  AsSociety?:   boolean
  Description?: string
  Date?:        Date
  Trash?:       string[]
}
