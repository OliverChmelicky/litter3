import {Component, Inject, OnInit} from '@angular/core';
import {UserModel, FriendRequestModel, FriendsModel} from "../../models/user.model";
import {UserService} from "../../services/user/user.service";
import {UserViewModel} from "./friendRequestView";
import {
  friendsColumnsDefinition,
  requestsSendColumnsDefinition,
  societiesColumnsDefinition,
  requestsReceivedColumnsDefinition,
  myCollectionsColumns,
  myEventsColumns
} from "./table-definitions";
import {SocietyModel} from "../../models/society.model";
import {SocietyService} from "../../services/society/society.service";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material/dialog";
import {CollectionModel} from "../../models/trash.model";
import {EventModel} from "../../models/event.model";
import {Router} from "@angular/router";
import {TrashService} from "../../services/trash/trash.service";
import {FileuploadService} from "../../services/fileupload/fileupload.service";
import {MatTableDataSource} from "@angular/material/table";
import {ApisModel} from "../../api/api-urls";
import {AuthService} from "../../services/auth/auth.service";

export interface ProfileDialogData {
  viewName?: string;
  firstName: string;
  lastName: string;
  newPicture: FormData;
  deletePicture: boolean;
  updateUser: boolean;
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
  imageUrlPrefix: string = ApisModel.apiUrl + '/' + ApisModel.fileupload + '/' + ApisModel.user + '/load/';

  constructor(
    private router: Router,
    private userService: UserService,
    private societyService: SocietyService,
    private trashService: TrashService,
    private authService: AuthService,
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
        if (user.Collections) {
          this.fillCollections(user.Collections)
        }
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
        this.userService.getMyFriendsIds().subscribe(
          relationship => {
            if (relationship != null) {
              this.myFriendsIds = relationship;
              this.fetchUserDetailsForFriends(relationship)
            }
          },
          error => console.log('Error GetMyFriends ', error)
        )
      })

  }

  openDialog(): void {
    const dialogRef = this.editProfileDialog.open(EditProfileComponent, {
      width: '800px',
      data: {
        viewName: this.me.FirstName + ' ' + this.me.LastName,
        firstName: this.me.FirstName,
        lastName: this.me.LastName,
        newPicture: new FormData(),
        deletePicture: false,
        updateUser: false,
        deleteAccount: false,
      }
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        if (result.deleteAccount) {
          this.userService.deleteAccount().subscribe(res => {
            this.router.navigateByUrl('register')
            this.authService.deleteUser()
          })
          return
        }
        if (result.updateUser) {
          if (result.firstName != this.me.FirstName || result.lastName != this.me.LastName) {
            this.me.FirstName = result.firstName;
            this.me.LastName = result.lastName;
            this.userService.updateUser(this.me).subscribe()
          }
          if (result.newPicture.has('file')) {
            this.fileuploadService.uploadUserImage(result.newPicture).subscribe( () => window.location.reload())
          } else if (result.deletePicture) {
            this.fileuploadService.deleteUserImage().subscribe(
              () => window.location.reload()
            )
          }
        }
      }
    });
  }

  removeFriend(userId: string) {
    this.userService.removeFriend(userId).subscribe(
      () => {
        const index = this.myFriendsView.findIndex(u => u.UserId === userId)
        this.myFriendsView.splice(index, 1)
        this.reinitMyFriendsTable()
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
        this.reinitReceivedRequestsTable()
        this.reinitMyFriendsTable()
      }
    )
  }

  denyFriendRequest(userId: string) {
    this.userService.denyFriend(userId).subscribe(
      () => {
        const index = this.IreceivedFriendRequests.findIndex(u => u.UserId === userId)
        this.IreceivedFriendRequests.splice(index, 1)
        this.reinitReceivedRequestsTable()
      }
    )
  }

  cancelFriendRequest(userId: string) {
    this.userService.denyFriend(userId).subscribe(
      () => {
        const index = this.IsendFriendRequests.findIndex(u => u.UserId === userId)
        this.IsendFriendRequests.splice(index, 1)
        this.reinitSendRequestTable()
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
          console.log(user.Id === friendship.User1Id || user.Id === friendship.User2Id)
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
    this.reinitMyFriendsTable()
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


  reinitSendRequestTable() {
    const newData = new MatTableDataSource<UserViewModel>(this.IsendFriendRequests);
    this.IsendFriendRequests = []
    for (let i = 0; i < newData.data.length; i++) {
      this.IsendFriendRequests.push(newData.data[i])
    }
  }

  reinitReceivedRequestsTable() {
    const newData = new MatTableDataSource<UserViewModel>(this.IreceivedFriendRequests);
    this.IreceivedFriendRequests = []
    for (let i = 0; i < newData.data.length; i++) {
      this.IreceivedFriendRequests.push(newData.data[i])
    }
  }

  reinitMyFriendsTable() {
    const newData = new MatTableDataSource<UserViewModel>(this.myFriendsView);
    this.myFriendsView = []
    for (let i = 0; i < newData.data.length; i++) {
      this.myFriendsView.push(newData.data[i])
    }
  }

  reinitMyCollectionsTable() {
    const newData = new MatTableDataSource<CollectionModel>(this.myCollections);
    this.myCollections = []
    for (let i = 0; i < newData.data.length; i++) {
      this.myCollections.push(newData.data[i])
    }
  }


  private fillCollections(collections: CollectionModel[]) {
    collections.map(c => {
      if (!c.Images) {
        c.Images = [{
          Url: '',
          CollectionId: '',
        }]
      }
    })
    this.myCollections = collections
  }
}

@Component({
  selector: 'app-edit-profile',
  templateUrl: './dialog/edit-profile/edit-profile.component.html',
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
    this.dialogRef.close(this.data);
  }

  onUpload(event) {
    this.data.newPicture.delete('file')
    this.data.newPicture.append("file", event.target.files[0], event.target.files[0].name)
    this.data.updateUser = true
  }

  onSaveChanges() {
    this.data.updateUser = true
    this.dialogRef.close(this.data);
  }
}

@Component({
  selector: 'app-collection-detail',
  templateUrl: './dialog/collection-details/collection-details.component.html',
  styleUrls: ['./dialog/collection-details/collection-details.component.css']
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
