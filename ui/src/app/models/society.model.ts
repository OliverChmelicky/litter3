import {UserModel} from "./user.model";
import {PagingModel} from "./shared.models";

export interface SocietyModel {
  Id: string;
  Name: string;
  Avatar: string;
  Users?: UserModel[];
  Description?: string;
  CreatedAt: Date;
}

export interface SocietyAnswSimpleModel {
  Id: string;
  Name: string;
  Avatar: string;
  UsersNumb: number;
  Description?: string;
  CreatedAt: Date;
}

export interface SocietyWithPagingAnsw {
  Societies: SocietyAnswSimpleModel[];
  Paging: PagingModel;
}

export interface UserGroupRequestModel {
  UserId: string;
  SocietyId: string;
}

export interface MemberModel {
  UserId: string;
  SocietyId: string;
  Permission: string;
  CreatedAt: Date;
}

export interface ApplicantModel {
  UserId: string;
  SocietyId: string;
  CreatedAt: Date;
}
