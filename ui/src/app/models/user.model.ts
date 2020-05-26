import {MemberModel, SocietyModel} from "./society.model";
import {CollectionModel} from "./trash.model";
import {EventModel} from "./event.model";

export interface UserModel {
  Id: string;
  FirstName: string;
  LastName: string;
  Email: string;
  Uid: string;
  Avatar: string;
  Societies?: SocietyModel[];
  Collections?: CollectionModel[];
  Events?: EventModel[];
  CreatedAt: Date;
}

export interface UserGroupModel {
  UserId: string,
  SocietyId: string,
}

export interface ShowUsersInMatSelect {
  email: string;
  id: string;
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
  User1Id: string;
  User2Id: string;
  CreatedAt: Date;
}

export interface UserInSocietyModel {
  user: UserModel,
  role: string,
  showRemove: boolean,
}
