import {MarkerModel} from "../components/google-map/Marker.model";

export interface IdMessageModel {
  Id: string;
}

export interface IdsMessageModel {
  Ids: string[];
}

export interface EmailMessageModel {
  Email: string;
}

export interface PagingModel {
  TotalCount: number;
  From: number;
  To: number;
}

export interface MapLoadAnswer {
  LoadedMarkers: MarkerModel[];
  BorderTop: number;
  BorderBottom: number;
  BorderLeft: number;
  BorderRight: number;
}

export interface AttendantsModel {
  name: string
  avatar: string
  role: string
}
