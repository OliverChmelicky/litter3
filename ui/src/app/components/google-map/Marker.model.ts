import {TrashImageModel} from "../../models/trash.model";

export interface MarkerModel {
  lat: number;
  lng: number;
  new?: boolean;
  id: string,
  cleaned?: boolean,
  images?: TrashImageModel[],
  numOfCollections?: number,
}
