import {UserModel} from "./user.model";
import {PagingModel} from "./shared.models";

export interface SocietyModel {
  Id?: string;
  Name: string;
  Avatar?: string;
  Description?: string;
  Users?: UserModel[];
  Applicants?: UserModel[],
  MemberRights?: MemberModel[],
  ApplicantsIds?: ApplicantModel,
  CreatedAt?: Date;
}

export const DefaultSociety: SocietyModel = {
  Name: '',
  Description: '',
  Avatar: ''
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

export interface UserSocietyRequestModel {
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
