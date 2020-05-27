import {Component, OnInit} from '@angular/core';
import {ShowUsersInMatSelect, UserModel} from "../../models/user.model";
import {CollectionImageModel, CollectionModel, defaultCollectionModel} from "../../models/trash.model";
import {ActivatedRoute} from "@angular/router";
import {TrashService} from "../../services/trash/trash.service";
import {UserService} from "../../services/user/user.service";
import {AuthService} from "../../services/auth/auth.service";
import {Location} from "@angular/common";
import {FormControl, FormGroup} from "@angular/forms";
import {FileuploadService} from "../../services/fileupload/fileupload.service";

@Component({
  selector: 'app-edit-collection-random',
  templateUrl: './edit-collection-random.component.html',
  styleUrls: ['./edit-collection-random.component.css']
})
export class EditCollectionRandomComponent implements OnInit {
  collection: CollectionModel = defaultCollectionModel;
  newImages: FormData = new FormData()
  deleteImages: string[] = [];
  me: UserModel

  friendsNotInCollection: UserModel[] = [];
  peopleInCollection: UserModel[] = [];
  allFriends: UserModel[] = [];

  showUsers: ShowUsersInMatSelect[] = []
  addFriends: FormControl = new FormControl(['']);
  collectionImages: CollectionImageModel[] = [];

  collectionForm: FormGroup;
  newWeight: number;
  newCleanedTrash: boolean;

  constructor(
    private trashService: TrashService,
    private userService: UserService,
    private authService: AuthService,
    private location: Location,
    private route: ActivatedRoute,
    private fileuploadService: FileuploadService,
  ) {
  }

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      const collectionId = params.get('collectionId');

      this.authService.isLoggedIn.subscribe(loggedId => {
        if (loggedId) {
          this.userService.getMe().subscribe(me => {
            this.me = me;
            this.getMyFriends();

            this.trashService.getCollectionById(collectionId).subscribe(c => {
              this.collection = c
              if (c.Images) {
                c.Images.map( col => {
                  this.collectionImages.push({
                    Url: col.Url,
                    CollectionId: col.CollectionId,
                    inDeleteList: false,
                  })
                })

              }
              this.newWeight = c.Weight
              this.newCleanedTrash = c.CleanedTrash
              if (c.Users) {
                this.peopleInCollection = c.Users
                this.filterFirendsAlreadyInCollection()
              }
            })
          })
        }
      })
    });
  }

  onUpdate() {
    const newCollection: CollectionModel = {
      Weight: this.newWeight,
      CleanedTrash: this.newCleanedTrash,
      Id: this.collection.Id,
      EventId: this.collection.EventId,
      TrashId: this.collection.TrashId,
      Users: this.collection.Users,
      Images: this.collection.Images,
      CreatedAt: this.collection.CreatedAt
    }

    if (this.collection.Weight != this.newWeight ||
      this.collection.CleanedTrash != this.newCleanedTrash) {
      this.trashService.updateCollection(newCollection).subscribe( )
    }
    console.log('mam tychto friendsL ', this.addFriends)
    // if (this.addFriends) {
    //   this.trashService.addFriendsToCollection().subscribe()
    // }
    if (this.newImages.has('files')) {
      this.fileuploadService.uploadCollectionImages(this.newImages, this.collection.Id).subscribe()
    }
    this.location.back();
  }

  onLeave() {
    this.trashService.deleteCollectionFromUser(this.collection.Id).subscribe(
      () => this.location.back()
    )
  }

  onUpload(event) {
    this.newImages.delete('files')
    for (let i = 0; i < event.target.files.length; i++) {
      this.newImages.append("files", event.target.files[i], event.target.files[i].name);
    }
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
                this.allFriends = users
                this.fillUsersForMatSelect(users)
              },
              err => console.log('Error fetching user details ', err)
            );
          }
        } else {
          this.allFriends = []
        }
      },
      error => console.log('Error GetMyFriends ', error)
    )
  }

  private fillUsersForMatSelect(users: UserModel[]) {
    this.showUsers = users.map(u => {
      return {
        email: u.Email,
        id: u.Id,
      }
    })
  }

  private filterFirendsAlreadyInCollection() {
    for (let i = 0; i < this.allFriends.length; i++) {
      let found = false
      for (let j = 0; j < this.collection.Users.length; j++) {
        if (this.allFriends[i].Id === this.collection.Users[j].Id) {
          found = true;
        }
      }

      if (!found) {
        this.friendsNotInCollection.push(this.allFriends[i])
      }
    }

  }

  onDeleteImage(url: string) {
    this.collectionImages.map(i => {
      if (i.Url === url) {
        i.inDeleteList = true
      }
    })
    this.deleteImages.push(url)
  }

  onRemoveFromDeleteList(url: any) {
    this.collectionImages.map(i => i.inDeleteList = false)
    const index = this.deleteImages.findIndex(delUrl => delUrl === url)
    this.deleteImages.splice(index, 1)
  }

}

