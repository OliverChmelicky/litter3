import {TrashModel} from "./trash.model";
import {PagingModel} from "./shared.models";

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

export interface EventUserModel {
  UserId: string,
  EventId: string,
  Permission:string,
}

export interface EventSocietyModel {
  SocietyId: string,
  EventId: string,
  Permission:string,
}

export interface EventModel {
  Id?: string,
  Date: Date,
  Description: string,
  CreatedAt?: Date,
  Trash?: TrashModel[],
  UsersIds?: EventUserModel[],
  SocietiesIds?: EventSocietyModel[],
}

export interface ListEventsModel {
  Id: string,
  Date: Date,
  NumOfAttendants: number,
}

export interface EventRequestModel {
  Id?: string,
  UserId: string,
  SocietyId: string,
  AsSociety: boolean,
  Description: string,
  Date: Date,
  Trash?: string[],
}

export interface EventPickerModel {
  AsSociety: boolean,
  Id: string,
  VisibleName: string
}

export interface AttendanceRequestModel {
  PickerId: string,
  EventId: string
  AsSociety: boolean,
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

export interface EventWithPagingAnsw {
  Events: EventModel[];
  Paging: PagingModel;
}
