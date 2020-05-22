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
  CreatedAt?: Date,
  Anonymously?: boolean,
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
  CreatedAt: null,
  Anonymously: false,
}

export interface CollectionModel {
  Id?: string,
  Weight: number,
  CleanedTrash: boolean,
  TrashId: string,
  EventId: string,
  Images: string[],
  CreatedAt?: Date,
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
  collectionImages: string[],

  isInList: boolean,
}

export const TrashTypeHousehold     = 0b00000000001
export const TrashTypeAutomotive    = 0b00000000010
export const TrashTypeConstruction  = 0b00000000100
export const TrashTypePlastics      = 0b00000001000
export const TrashTypeElectronic    = 0b00000010000
export const TrashTypeGlass         = 0b00000100000
export const TrashTypeMetal         = 0b00001000000
export const TrashTypeDangerous     = 0b00010000000
export const TrashTypeCarcass       = 0b00100000000
export const TrashTypeOrganic       = 0b01000000000
export const TrashTypeOther         = 0b10000000000

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
