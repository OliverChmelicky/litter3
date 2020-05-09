export interface MarkerModel {
  lat: number;
  lng: number;
  new?: boolean;
  id: string,
  cleaned?: boolean,
  images?: string[],
  numOfCollections?: number,
}
