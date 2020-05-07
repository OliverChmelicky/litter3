export interface TrashModel {
  Id: string,
  Cleaned: boolean,
  Size: string,
  Accessibility: string,
  TrashType: number,
  Location: number[],
  Description: string,
  FinderId: string,
  Collections?: Collection[],
  Images?: string[],
  CreatedAt?: Date,
}

export interface Collection {
  Id: string,
  Weight: number,
  CleanedTrash: boolean,
  TrashId: string,
  EventId: string,
  Images: string[],
  CreatedAt?: Date,
}
