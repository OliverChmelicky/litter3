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
  Images?: string[],
  CreatedAt?: Date,
  Anonymously?: boolean,
}

export interface CollectionModel {
  Id: string,
  Weight: number,
  CleanedTrash: boolean,
  TrashId: string,
  EventId: string,
  Images: string[],
  CreatedAt?: Date,
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
