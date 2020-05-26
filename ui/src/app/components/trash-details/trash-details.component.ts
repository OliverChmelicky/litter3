import {Component, Inject, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from "@angular/router";
import {TrashService} from "../../services/trash/trash.service";
import {
  CollectionModel,
  CommentViewModel,
  defaultCollectionModel, defaultTrashModel, defaultTrashTypeBooleanValues,
  TrashModel,
  TrashTypeBooleanValues
} from "../../models/trash.model";
import {GoogleMap} from "@agm/core/services/google-maps-types";
import {AuthService} from "../../services/auth/auth.service";
import {CollectionTableDisplayedColumns} from "./collectionTableModel";
import {UserModel} from "../../models/user.model";
import {UserService} from "../../services/user/user.service";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material/dialog";
import {FormControl} from "@angular/forms";
import {FileuploadService} from "../../services/fileupload/fileupload.service";


export interface DialogDataCreateCollection {
  collection: CollectionModel,
  collectionImages: FormData,
  friends: UserModel[],
  collectedWithFriends: UserModel[],
}

export interface DialogDataEditCollection {
  collection: CollectionModel,
  leaveCollection: boolean
  deleteImages: string[],
  uploadImages: FormData,
  friends: UserModel[],
  newFriends: UserModel[],
}

export interface DialogDataShowCollection {
  collection: CollectionModel,
}


export interface ShownCollectonsModel {
  collection: CollectionModel,
  canEdit: boolean,
}

export interface ShowUsersInModal {
  email: string;
  id: string;
}

@Component({
  selector: 'app-trash-details',
  templateUrl: './trash-details.component.html',
  styleUrls: ['./trash-details.component.css']
})
export class TrashDetailsComponent implements OnInit {
  isLoggedIn: boolean
  map: GoogleMap;
  trashId: string;
  trash: TrashModel = defaultTrashModel;
  trashTypeBool: TrashTypeBooleanValues = defaultTrashTypeBooleanValues;
  tableColumnsTrashCollections = CollectionTableDisplayedColumns;
  finder: UserModel = null;
  comments: CommentViewModel[] = [];
  message: string = '';

  me: UserModel;
  friends: UserModel[] = [];

  shownCollections: ShownCollectonsModel[] = [];

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private trashService: TrashService,
    private authService: AuthService,
    private userService: UserService,
    private fileuploadService: FileuploadService,
    private createCollectionRandomDialog: MatDialog,
    private editCollectionRandomDialog: MatDialog,
    private showCollectionRandomDialog: MatDialog,
  ) {
  }

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      this.trashId = params.get('id');
      this.trashService.getTrashById(this.trashId).subscribe(
        trash => {
          trash.Collections.map(c => {
            this.shownCollections.push({
              collection: c,
              canEdit: false,
            })
          })

          if (trash.FinderId) {
            this.userService.getUser(trash.FinderId).subscribe(u => {
              this.finder = u
            })
          }

          if (!trash.Collections) {
            trash.Collections = []
          }
          if (!trash.Images) {
            trash.Images = []
          }
          if (!trash.Comments) {
            trash.Comments = []
          }

          this.trash = trash
          this.trashTypeBool = this.trashService.convertTrashTypeNumToBools(this.trash.TrashType);

          if (trash.Comments.length > 0) {
            const usersCommented = trash.Comments.map(c => c.UserId);
            this.userService.getUsersDetails(usersCommented).subscribe(
              users => this.addUsersToComments(users)
            )
          }
          this.authService.isLoggedIn.subscribe(isLogged => {
            this.isLoggedIn = isLogged
            if (!isLogged) {
              return
            }
            this.userService.getMe().subscribe(me => {
              this.me = me
              this.getMyFriends();

              this.trashService.getIdsOfTrashOfUsers().subscribe(ids => {
                console.log('moje collections: ', ids)
                console.log('existing cols: ',this.trash.Collections )
                ids.map(u => {
                  this.shownCollections.map(c => {
                    if (u.CollectionId === c.collection.Id) {
                      c.canEdit = true
                    }
                  })
                })
                console.log('na konci mapovania:')
                const vals = this.shownCollections.map(a => a.canEdit)
                console.log(vals)
              })
            })
          })
        })
    });
  }

  onMapReady(map: GoogleMap) {
    this.map = map;
  }

  onEdit() {
    this.router.navigateByUrl('trash/edit/' + this.trash.Id)
  }

  // showCollectionDetails(Id: string) {
  //   this.router.navigateByUrl('collection/details/' + this.trash.Id)
  // }

  onCreateEvent() {
    this.router.navigateByUrl('events/create')
  }

  private addUsersToComments(users: UserModel[]) {
    let unsortedArray: CommentViewModel[] = []

    this.trash.Comments.map(c => {
      users.map(u => {
        if (u.Id === c.UserId) {
          unsortedArray.push({
            Id: c.Id,
            TrashId: c.TrashId,
            UserName: u.FirstName,
            Message: c.Message,
            CreatedAt: new Date(c.CreatedAt),
          })
        }
      })
    })

    this.comments = unsortedArray.sort((a, b) => a.CreatedAt.getTime() - b.CreatedAt.getTime())

  }

  commentOnTrash() {
    if (this.message.length > 0) {
      this.trashService.commentTrash(this.message, this.trash.Id).subscribe(
        rec => {
          this.comments.push({
            Id: rec.Id,
            TrashId: rec.TrashId,
            UserName: this.me.FirstName,
            Message: rec.Message,
            CreatedAt: new Date(rec.CreatedAt),
          })
        }
      )
    }
  }

  onCreateCollection() {
    this.router.navigate(['collection/create',this.trashId ])
    //badly shows window
    // const dialogRef = this.createCollectionRandomDialog.open(CreateCollectionRandomFromTrashComponent, {
    //   width: '800px',
    //   data: {
    //     collection: defaultCollectionModel,
    //     friends: this,
    //     collectedWithFriends: [],
    //   }
    // });
    //
    // dialogRef.afterClosed().subscribe(() => {
    // });
  }

  onShowCollection(collectionId: string) {
      let collection: CollectionModel = defaultCollectionModel
      this.trash.Collections.map(c => {
        if (c.Id === collectionId) {
          collection = c
        }
      })

      if (!collection.Images) {
        collection.Images = [];
      }

      const dialogRef = this.showCollectionRandomDialog.open(ShowCollectionFromTrashComponent, {
        width: '800px',
        data: {
          collection: collection,
        }
      });

      dialogRef.afterClosed().subscribe(() => {});
  }

  onEditCollection(collectionId: string) {
    let collection: CollectionModel = defaultCollectionModel
    this.trash.Collections.map(c => {
      if (c.Id === collectionId) {
        collection = c
      }
    })

    if (!collection.Images) {
      collection.Images = [];
    }

    const oldCollection: CollectionModel = {
      Weight: collection.Weight,
      CleanedTrash: collection.CleanedTrash,
      Id: collection.Id,
      EventId: collection.TrashId,
      TrashId: collection.TrashId,
      Users: collection.Users,
      Images: collection.Images,
      CreatedAt: collection.CreatedAt
    }

    const dialogRef = this.editCollectionRandomDialog.open(EditCollectionRandomFromTrashComponent, {
      width: '800px',
      data: {
        collection: collection,
        leaveCollection: false,
        deleteImages: [],
        uploadImages: new FormData(),
        friends: this.friends,
        newFriends: [],
      }
    });

    dialogRef.afterClosed().subscribe(res => {
      console.log('upravil som to tak, ze:', res)
      if (res.leaveCollection) {

      } else {
        console.log(oldCollection.Weight != res.collection.Weight ||
          oldCollection.CleanedTrash != res.collection.CleanedTrash)
        if (oldCollection.Weight != res.collection.Weight ||
          oldCollection.CleanedTrash != res.collection.CleanedTrash) {
          console.log('updatujem')
          this.trashService.updateCollection(res.collection).subscribe()
        }
        //TODO new friends
        if (res.uploadImages.has('files')) {
          this.fileuploadService.uploadCollectionImages(res.uploadImages, collectionId).subscribe()
        }
      }
    });
  }


  private getMyFriends() {
    this.userService.getMyFriendsIds().subscribe(relationship => {
        if (relationship != null) {
          const userIds = relationship.map(friend => {
            if (friend.User1Id !== this.me.Id)
              return friend.User1Id;
            if (friend.User2Id !== this.me.Id)
              return friend.User2Id;
          });
          if (userIds.length !== 0) {
            this.userService.getUsersDetails(userIds).subscribe(
              users => {
                this.friends = users
              },
              err => console.log('Error fetching user details ', err)
            );
          }
        } else {
          this.friends = []
        }
      },
      error => console.log('Error GetMyFriends ', error)
    )
  }
}

//DIALOG INFO

@Component({
  selector: 'app-show-collection',
  templateUrl: './dialog-collection-detail/detail-dialog.component.html',
  styleUrls: ['./dialog-collection-detail/detail-dialog.component.css']
})
export class ShowCollectionFromTrashComponent {

  constructor( public dialogRef: MatDialogRef<ShowCollectionFromTrashComponent>,
               @Inject(MAT_DIALOG_DATA) public data: DialogDataShowCollection) {
    console.log(data)
  }

  onNoClick(): void {
    this.dialogRef.close();
  }

}


@Component({
  selector: 'app-edit-collection',
  templateUrl: './dialog-edit-collection/edit-collection-dialog.component.html',
  styleUrls: ['./dialog-edit-collection/edit-collection-dialog.component.css'],
})
export class EditCollectionRandomFromTrashComponent {
  images = []
  show: boolean = true;

  constructor(public dialogRef: MatDialogRef<EditCollectionRandomFromTrashComponent>,
              @Inject(MAT_DIALOG_DATA) public data: DialogDataEditCollection) {
    console.log('images in: ', this.data.collection.Images)
    data.collection.Images.map(i => this.images.push({
      Url: i.Url,
      CollectionId: i.CollectionId,
    }))
    console.log('images in dialog: ', this.images)
  }

  onNoClick(): void {
    this.dialogRef.close();
  }

  onDelete(url: string) {
    this.data.deleteImages.push(url)
    const index = this.images.findIndex(i => i.Url === url)
    this.images[index] = true
    this.reload()
  }

  reload() {
    this.show = false;
    setTimeout(() => this.show = true);
  }

  onUpload(event) {
    this.data.uploadImages.delete('files')
    for (let i = 0; i < event.target.files.length; i++) {
      this.data.uploadImages.append("files", event.target.files[i], event.target.files[i].name);
    }
  }

  onSetDelete() {
    this.data.leaveCollection = true
    this.data.uploadImages = new FormData()
    this.dialogRef.close({
      leaveCollection: true,
      collection: this.data.collection
    });
  }
}


//
// @Component({
//   selector: 'app-create-collection-random',
//   templateUrl: './dialog-create-collection-random/create-collection-dialog.component.html',
//   styleUrls: ['./dialog-create-collection-random/create-collection-dialog.component.css'],
// })
// export class CreateCollectionRandomFromTrashComponent {
//   showUsers: ShowUsersInModal[];
//   selectedFriends: FormControl;
//
//   constructor(public dialogRef: MatDialogRef<CreateCollectionRandomFromTrashComponent>,
//               @Inject(MAT_DIALOG_DATA) public data: DialogDataCreateCollection) {
//     if (data.friends) {
//       this.showUsers = this.data.friends.map(f => {
//         return {
//           email: f.Email,
//           id: f.Id,
//         }
//       })
//     }
//     this.selectedFriends = new FormControl()
//   }
//
//   onNoClick(): void {
//     this.dialogRef.close();
//   }
//
//
//   onUpload(event) {
//     this.data.collectionImages.delete('files')
//     for (let i = 0; i < event.target.files.length; i++) {
//       this.data.collectionImages.append("files", event.target.files[i], event.target.files[i].name);
//     }
//   }
//
//   onCreate() {
//     console.log('From form controll: ', this.selectedFriends)
//     this.dialogRef.close({
//       collection: this.data.collection,
//       collectionImages: this.data.collectionImages,
//       collectionWithFriends:  this.selectedFriends
//     });
//   }
// }






