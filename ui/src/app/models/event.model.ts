import {CollectionModel, TrashModel} from "./trash.model";
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
  ChangingToSociety: boolean,
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

export const defaultEventModel: EventModel = {
  Id: '',
  Date: new Date,
  Description: '',
  CreatedAt: new Date,
  Trash: [],
  UsersIds: [],
  SocietiesIds: [],
}

export interface EventWithCollectionsModel {
  Id?: string,
  Date: Date,
  Description: string,
  CreatedAt?: Date,
  Trash?: TrashModel[],
  UsersIds?: EventUserModel[],
  SocietiesIds?: EventSocietyModel[],
  Collections?: CollectionModel[],
}

export interface EventModelTable {
  id: string,
  date: Date,
  attendingPeople:number;
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

export interface EventWithPagingAnsw {
  Events: EventModel[];
  Paging: PagingModel;
}

interface rolesInterface {
  key:string,
  value: string,
}

export const roles: rolesInterface[] = [
  {
    key:'editor',
    value: 'Editor',
  },
  {
    key:'viewer',
    value: 'Viewer',
  },
]
