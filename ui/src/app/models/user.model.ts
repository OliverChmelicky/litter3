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

export interface SocietyModel {
  Id: string;
  Name: string;
  Avatar: string;
  Users?: UserModel[];
  CreatedAt: Date;
}

export interface IdMessageModel {
  Id: string;
}

export interface IdsMessageModel {
  Ids: string[];
}

export interface EmailMessageModel {
  Email: string;
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
