import {UserModel} from "./user.model";

export interface CommentModel {
  Id: string,
  UserId: string,
  TrashId: string,
  Message: string,
  CreatedAt?: Date,
}

export interface CommentViewModel {
  Id: string,
  UserName: string,
  TrashId: string,
  Message: string,
  CreatedAt: Date,
}

export interface TrashModel {
  Id: string,
  Cleaned: boolean,
  Size: string,
  Accessibility: string,
  TrashType: number,
  Location: number[],
  Description: string,
  FinderId: string,
  Collections?: CollectionModel[],
  Images?: TrashImageModel[],
  Comments?: CommentModel[],
  CreatedAt?: Date,
  Anonymously?: boolean,
}

export interface MarkerCollectionModel {
  lat: number;
  lng: number;
  trashId: string,
  cleaned?: boolean,
  image: TrashImageModel,
  numOfCollections?: number,

  collectionWeight: number,
  collectionCleanedTrash: boolean,
  collectionEventId: string,
  collectionImages: FormData,

  isInList: boolean,
}

export interface TrashImageModel {
  Url: string,
  TrashId: string,
}

export const defaultTrashImage: TrashImageModel = {
  Url: '',
  TrashId: '',
}

export const defaultTrashModel: TrashModel = {
  Id: '',
  Cleaned: false,
  Size: '',
  Accessibility: '',
  TrashType: 0,
  Location: [],
  Description: '',
  FinderId: '',
  Collections: [],
  Images: [],
  Comments: [],
  CreatedAt: null,
  Anonymously: false,
}

export interface CollectionModel {
  Id?: string,
  Weight: number,
  CleanedTrash: boolean,
  TrashId: string,
  EventId: string,
  Users: UserModel[],
  Images: CollectionImageModel[],
  CreatedAt?: Date,
}

export interface CollectionUserModel {
  UserId: string,
  CollectionId: string,
}

export interface CreateCollectionModel {
  Collections: CollectionModel[],
  AsSociety: boolean,
  OrganizerId: string,
  EventId: string,
  CreatedAt?: Date,
}

export interface UpdateCollectionModel {
  Collection: CollectionModel,
  AsSociety: boolean,
  OrganizerId: string,
  EventId: string,
}

export const defaultCollectionModel: CollectionModel = {
  Id: '',
  Weight: 0,
  CleanedTrash: false,
  TrashId: '',
  EventId: '',
  Users: null,
  Images: [],
  CreatedAt: null,
}

export interface CollectionImageModel {
  Url: string,
  CollectionId: string,
}

export interface CreateCollectionRandomRequest {
  TrashId:      string,
  CleanedTrash: boolean,
  Weight:       number,
  Friends:      string[],
}

export const TrashTypeHousehold = 0b00000000001
export const TrashTypeAutomotive = 0b00000000010
export const TrashTypeConstruction = 0b00000000100
export const TrashTypePlastics = 0b00000001000
export const TrashTypeElectronic = 0b00000010000
export const TrashTypeGlass = 0b00000100000
export const TrashTypeMetal = 0b00001000000
export const TrashTypeDangerous = 0b00010000000
export const TrashTypeCarcass = 0b00100000000
export const TrashTypeOrganic = 0b01000000000
export const TrashTypeOther = 0b10000000000

export const TrashTypeMask = 0b11111111111

export interface TrashTypeBooleanValues {
  TrashTypeHousehold: boolean,
  TrashTypeAutomotive: boolean,
  TrashTypeConstruction: boolean,
  TrashTypePlastics: boolean,
  TrashTypeElectronic: boolean,
  TrashTypeGlass: boolean,
  TrashTypeMetal: boolean,
  TrashTypeDangerous: boolean,
  TrashTypeCarcass: boolean,
  TrashTypeOrganic: boolean,
  TrashTypeOther: boolean,
}

export const defaultTrashTypeBooleanValues: TrashTypeBooleanValues = {
  TrashTypeHousehold: false,
  TrashTypeAutomotive: false,
  TrashTypeConstruction: false,
  TrashTypePlastics: false,
  TrashTypeElectronic: false,
  TrashTypeGlass: false,
  TrashTypeMetal: false,
  TrashTypeDangerous: false,
  TrashTypeCarcass: false,
  TrashTypeOrganic: false,
  TrashTypeOther: false,
}
