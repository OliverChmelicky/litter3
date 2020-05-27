import { Component, OnInit } from '@angular/core';
import {ActivatedRoute, Router} from "@angular/router";
import {TrashService} from "../../services/trash/trash.service";
import {UserService} from "../../services/user/user.service";
import {ShowUsersInMatSelect, UserModel} from "../../models/user.model";
import {FormBuilder, FormControl, Validators} from "@angular/forms";
import {CreateCollectionRandomRequest} from "../../models/trash.model";
import {FileuploadService} from "../../services/fileupload/fileupload.service";


@Component({
  selector: 'app-create-collection-random',
  templateUrl: './create-collection-random.component.html',
  styleUrls: ['./create-collection-random.component.css']
})
export class CreateCollectionRandomComponent implements OnInit {
  me: UserModel;
  trashId: string;
  friends: UserModel[] = [];
  weight: number;

  showUsers: ShowUsersInMatSelect[];
  uploadImages = new FormData();
  errorMessage: string = '';

  selectedFriends: FormControl = new FormControl();
  cleanedtrash: boolean = false;

  collectionForm = this.formBuilder.group({
    selectedFriends: this.formBuilder.control(['']),
    cleanedTrash: this.formBuilder.control(false),
    weight: this.formBuilder.control(0, [Validators.required, Validators.min(0)])
  });

  constructor(
    private route: ActivatedRoute,
    private trashService: TrashService,
    private userService: UserService,
    private formBuilder: FormBuilder,
    private fileUpload: FileuploadService,
    private router: Router,
  ) { }

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      this.trashId = params.get('trashId');
      this.userService.getMe().subscribe(me => {
        this.me = me;
        this.getMyFriends();
      })
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
                this.fillUsersForMatSelect(users)
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

  private fillUsersForMatSelect(users: UserModel[]) {
    this.showUsers = users.map( u => {
      return {
        email: u.Email,
        id: u.Id,
      }
    })
  }

  onUpload(event) {
    this.uploadImages.delete('files')
    for (let i = 0; i < event.target.files.length; i++) {
      this.uploadImages.append("files", event.target.files[i], event.target.files[i].name);
    }
  }

  onCreate() {
    if (this.weight <= 0) {
      this.errorMessage = 'Collection needs to have more than 0 kilograms'
      return
    }
    const friendsIds = this.selectedFriends.value
    const collectionRequest: CreateCollectionRandomRequest = {
      TrashId: this.trashId,
      CleanedTrash: this.cleanedtrash,
      Weight: this.weight,
      Friends: friendsIds,
    }
    console.log('idem vytvorit: ',collectionRequest)
    this.trashService.createCollection(collectionRequest).subscribe(
      res => {
        console.log('vytvorene')
        if (this.uploadImages.has('files')) {
          this.fileUpload.uploadCollectionImages(this.uploadImages, res.Id).subscribe(
            () => this.router.navigate(['trash/details/', this.trashId])
          )
        } else {
          this.router.navigate(['trash/details/', this.trashId])
        }
      }
    )
  }
}

