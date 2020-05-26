import {Component, Inject, OnInit} from '@angular/core';
import {UserModel, FriendRequestModel, FriendsModel} from "../../models/user.model";
import {UserService} from "../../services/user/user.service";
import {UserViewModel} from "./friendRequestView";
import {
  friendsColumnsDefinition,
  requestsSendColumnsDefinition,
  societiesColumnsDefinition,
  requestsReceivedColumnsDefinition,
  myCollectionsColumns, myEventsColumns
} from "./table-definitions";
import {SocietyModel} from "../../models/society.model";
import {SocietyService} from "../../services/society/society.service";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material/dialog";
import {CollectionModel} from "../../models/trash.model";
import {EventModel} from "../../models/event.model";
import {Router} from "@angular/router";
import {TrashService} from "../../services/trash/trash.service";
import {FileuploadService} from "../../services/fileupload/fileupload.service";

export interface ProfileDialogData {
  viewName?: string;
  firstName: string;
  lastName: string;
  email: string;
  deleteAccount: boolean;
}

export interface CollectionDialogData {
  collection: CollectionModel;
  update: boolean;
  deleteFromUser: boolean;
  deletedImages: string[],
  newImages: FormData,
}

@Component({
  selector: 'app-my-profile',
  templateUrl: './my-profile.component.html',
  styleUrls: ['./my-profile.component.css']
})
export class MyProfileComponent implements OnInit {
  me: UserModel;
  IsendFriendRequests: UserViewModel[] = [];
  IreceivedFriendRequests: UserViewModel[] = [];
  myFriendsIds: FriendsModel[] = [];
  myFriendsView: UserViewModel[] = [];
  mySocietiesView: SocietyModel[];
  newFriendEmail: string;
  myCollections: CollectionModel[] = [];
  myEvents: EventModel[] = [];

  friendsColumns = friendsColumnsDefinition;
  societiesColumns = societiesColumnsDefinition;
  requestsSendColumns = requestsSendColumnsDefinition;
  requestsReceivedColumns = requestsReceivedColumnsDefinition;
  myCollectionsColumns = myCollectionsColumns;
  myEventsColumns = myEventsColumns;

  constructor(
    private router: Router,
    private userService: UserService,
    private societyService: SocietyService,
    private trashService: TrashService,
    private fileuploadService: FileuploadService,
    public editProfileDialog: MatDialog,
    public showCollectionDialog: MatDialog,
  ) {
  }

  ngOnInit() {
    this.userService.getMe().subscribe(
      user => {
        this.me = user;
        this.mySocietiesView = user.Societies
        this.myCollections = user.Collections
        this.myEvents = user.Events
        this.userService.getMyFriendRequests().subscribe(
          requests => {
            if (requests != null) {
              this.fetchUserDetailsForFriendRequests(requests)
            }
          },
          err => {
            console.log('Error fetching my requests ', err)
          },
        )
      })

    this.userService.getMyFriendsIds().subscribe(
      relationship => {
        if (relationship != null) {
          this.myFriendsIds = relationship;
          this.fetchUserDetailsForFriends(relationship)
        }
      },
      error => console.log('Error GetMyFriends ', error)
    )

  }

  openDialog(): void {
    const dialogRef = this.editProfileDialog.open(EditProfileComponent, {
      width: '800px',
      data: {
        firstName: this.me.FirstName,
        lastName: this.me.LastName,
        email: this.me.Email,
        deleteAccount: false,
      }
    });

    dialogRef.afterClosed().subscribe(result => {
      console.log(result.firstName)
      console.log(result.data)
      if (result) {
        console.log(result.deleteAccount)
        if (result.deleteAccount) {
          this.userService.deleteAccount().subscribe(res => console.log(res))
          return
        }
        if (result.firstName) {
          this.me.FirstName = result.FirstName;
        }
        if (result.LastName) {
          this.me.LastName = result.LastName;
        }
        if (result.firstName != this.me.FirstName || result.lastName != this.me.LastName) {
          this.userService.updateUser(this.me).subscribe()
        }
      }
    });
  }

  removeFriend(userId: string) {
    this.userService.removeFriend(userId).subscribe(
      () => {
        const index = this.myFriendsView.findIndex(u => u.UserId === userId)
        this.myFriendsView.splice(index, 1)
      }
    )

  }

  sendFriendRequest() {
    if (this.newFriendEmail != '') {
      this.userService.requestFriend(this.newFriendEmail).subscribe(
        newRequest => {
          this.userService.getUserByEmail(this.newFriendEmail).subscribe(
            user => {
              this.pushUserToFriendRequests([newRequest], [user])
              this.newFriendEmail = ''
            },
            error => console.log('Err getUser by email'))
        },)
    }
  }

  acceptFriendRequest(userId: string) {
    this.userService.acceptFriend(userId).subscribe(
      () => {
        const index = this.IreceivedFriendRequests.findIndex(u => u.UserId === userId)
        this.myFriendsView.push(this.IreceivedFriendRequests[index])
        this.IreceivedFriendRequests.splice(index, 1)
      }
    )
  }

  denyFriendRequest(userId: string) {
    this.userService.denyFriend(userId).subscribe(
      () => {
        const index = this.IreceivedFriendRequests.findIndex(u => u.UserId === userId)
        this.IreceivedFriendRequests.splice(index, 1)
      }
    )
  }

  cancelFriendRequest(userId: string) {
    this.userService.denyFriend(userId).subscribe(
      () => {
        const index = this.IsendFriendRequests.findIndex(u => u.UserId === userId)
        this.IsendFriendRequests.splice(index, 1)
      }
    )
  }


  private fetchUserDetailsForFriendRequests(requests: FriendRequestModel[]) {
    const userIds = requests.map(r => {
      if (r.User1Id !== this.me.Id)
        return r.User1Id;
      if (r.User2Id !== this.me.Id)
        return r.User2Id;
    });
    if (userIds.length !== 0) {
      this.userService.getUsersDetails(userIds).subscribe(
        users => {
          this.pushUserToFriendRequests(requests, users)
        },
        err => console.log('Error fetching user details ', err)
      );
    }
  }

  private fetchUserDetailsForFriends(friends: FriendsModel[]) {
    const userIds = friends.map(friend => {
      if (friend.User1Id !== this.me.Id)
        return friend.User1Id;
      if (friend.User2Id !== this.me.Id)
        return friend.User2Id;
    });
    if (userIds.length !== 0) {
      this.userService.getUsersDetails(userIds).subscribe(
        users => {
          this.pushUserToMyFriends(users)
        },
        err => console.log('Error fetching user details ', err)
      );
    }
  }

  private pushUserToMyFriends(users: UserModel[]) {
    this.myFriendsIds.map(friendship => {
      users.map(user => {
          if (user.Id === friendship.User1Id || user.Id === friendship.User2Id) {
            this.myFriendsView.push(
              {
                UserId: user.Id,
                FirstName: user.FirstName,
                LastName: user.LastName,
                Email: user.Email,
                Avatar: user.Avatar,
                CreatedAt: friendship.CreatedAt,
              }
            )
          }
        }
      )
    })
  }

  private pushUserToFriendRequests(requests: FriendRequestModel[], users: UserModel[]) {
    this.IreceivedFriendRequests = []
    this.IsendFriendRequests = []

    requests.map(request => {
      users.map(user => {
        if (request.User1Id == user.Id) {
          this.IreceivedFriendRequests.push(
            {
              UserId: user.Id,
              FirstName: user.FirstName,
              LastName: user.LastName,
              Email: user.Email,
              Avatar: user.Avatar,
              CreatedAt: request.CreatedAt,
            }
          )
        } else if (request.User2Id == user.Id) {
          this.IsendFriendRequests.push(
            {
              UserId: user.Id,
              FirstName: user.FirstName,
              LastName: user.LastName,
              Email: user.Email,
              Avatar: user.Avatar,
              CreatedAt: request.CreatedAt,
            }
          )
        }
      })
    })
  }

  onSocietyDetails(socId: string) {
    this.router.navigate(['societies/', socId])
  }

  onGoToEvent(eventId: string) {
    this.router.navigate(['events/details/', eventId])
  }

  showCollectionDetails(collectionId: string) {
    let index = this.myCollections.findIndex(c => c.Id === collectionId)
    let data: CollectionDialogData = {
      collection: this.myCollections[index],
      update: false,
      deleteFromUser: false,
      deletedImages: [],
      newImages: new FormData()
    }

    const dialogRef = this.showCollectionDialog.open(ShowCollectionRandomDetails, {
      width: '800px',
      data: data,
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result.deleteFromUser) {
        this.trashService.deleteCollectionFromUser(collectionId).subscribe()
        //TODO update table
      } else {
        if (result.update) {
          this.trashService.updateCollection(result.data.collection).subscribe()
        }
        if (result.deletedImages.length > 0) {
          //TODO delete images of col rand
        }
        if (result.newImages.has('files')) {
          this.fileuploadService.uploadCollectionImages(result.uploadImages, collectionId).subscribe()
        }
      }
    });
  }

}

@Component({
  selector: 'app-edit-profile',
  templateUrl: './dialog/edit-profile.component.html',
  //styleUrls: ['./dialog/edit-profile.component.css']
})
export class EditProfileComponent {

  constructor(public dialogRef: MatDialogRef<EditProfileComponent>,
              @Inject(MAT_DIALOG_DATA) public data: ProfileDialogData) {
    this.data.viewName = this.data.firstName + ' ' + this.data.lastName
  }

  onNoClick(): void {
    this.dialogRef.close();
  }

  onDeleteAccount() {
    this.data.deleteAccount = true;
    this.dialogRef.close({data: this.data});
  }

}

@Component({
  selector: 'app-collection-detail',
  templateUrl: './dialog/collection-details.component.html',
  //styleUrls: ['./dialog/collection-details.component.css']
})
export class ShowCollectionRandomDetails {

  //in this wnidow user can switch to edit mode and edit collection
  //Inputs will be enabeled and update button will be shown
  //images can be deleted and new files loaded
  editMode: boolean = false

  constructor(public dialogRef: MatDialogRef<ShowCollectionRandomDetails>,
              @Inject(MAT_DIALOG_DATA) public data: CollectionDialogData) {
  }

  onNoClick(): void {
    this.dialogRef.close();
  }

}
