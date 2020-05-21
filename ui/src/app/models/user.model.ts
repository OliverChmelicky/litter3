import {MemberModel, SocietyModel} from "./society.model";

export interface UserModel {
  Id: string;
  FirstName: string;
  LastName: string;
  Email: string;
  Uid: string;
  Avatar: string;
  Admins?:    MemberModel[];
  Societies?: SocietyModel[];
  CreatedAt: Date;
}

export interface UserGroupModel {
  UserId: string,
  SocietyId: string,
}

export const loggoutUser = {
  Id: '',
  FirstName: '',
  LastName: '',
  Email: '',
  Uid: '',
  Avatar: '',
  CreatedAt: null,
}

export interface FriendsModel {
  User1Id: string;
  User2Id: string;
  CreatedAt: Date;
}

export interface FriendRequestModel {
  User1Id:   string;
  User2Id:   string;
  CreatedAt: Date;
}

export interface UserInSocietyModel {
  user: UserModel,
  role: string,
  showRemove: boolean,
}
