import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from "@angular/router";
import {TrashService} from "../../services/trash/trash.service";
import {CommentModel, CommentViewModel, TrashModel, TrashTypeBooleanValues} from "../../models/trash.model";
import {GoogleMap} from "@agm/core/services/google-maps-types";
import {AuthService} from "../../services/auth/auth.service";
import {CollectionTableDisplayedColumns} from "./collectionTableModel";
import {UserModel} from "../../models/user.model";
import {UserService} from "../../services/user/user.service";

@Component({
  selector: 'app-trash-details',
  templateUrl: './trash-details.component.html',
  styleUrls: ['./trash-details.component.css']
})
export class TrashDetailsComponent implements OnInit {
  isLoggedIn: boolean
  map: GoogleMap;
  trashId: string;
  trash: TrashModel;
  trashTypeBool: TrashTypeBooleanValues;
  tableColumnsTrashCollections = CollectionTableDisplayedColumns;
  finder: UserModel = null;
  comments: CommentViewModel[] =[];
  message: string = '';
  me: UserModel;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private trashService: TrashService,
    private authService: AuthService,
    private userService: UserService,
  ) {
  }

  ngOnInit(): void {
    this.authService.isLoggedIn.subscribe( isLogged => {
      this.isLoggedIn = isLogged
      this.userService.getMe().subscribe(me => this.me = me)
    })

    this.route.paramMap.subscribe(params => {
      this.trashId = params.get('id');
      this.trashService.getTrashById(this.trashId).subscribe(
        trash => {
          console.log(trash)
          console.log('Comments: ',trash.Comments)

          if (trash.FinderId) {
            this.userService.getUser(trash.FinderId).subscribe( u => {
              this.finder = u
              console.log('Finder je: ', this.finder)
            })
          }

          if (!trash.Collections) {
            trash.Collections = []
          }
          if (!trash.Images) {
            trash.Images = []
          }
          if (!trash.Comments){
            trash.Comments = []
          }

          this.trash = trash
          this.trashTypeBool = this.trashService.convertTrashTypeNumToBools(this.trash.TrashType);

          if (trash.Comments.length > 0){
            const usersCommented = trash.Comments.map( c => c.UserId);
            this.userService.getUsersDetails(usersCommented).subscribe(
              users => this.addUsersToComments(users)
            )
          }
        })
    });
  }

  onMapReady(map: GoogleMap) {
    this.map = map;
  }

  onEdit() {
    this.router.navigateByUrl('trash/edit/'+this.trash.Id)
  }

  showCollectionDetails(Id: string) {
    this.router.navigateByUrl('collection/details/'+this.trash.Id)
  }

  onCreateEvent() {
    this.router.navigateByUrl('events/create')
  }

  private addUsersToComments(users: UserModel[]) {
    let unsortedArray: CommentViewModel[] = []

    this.trash.Comments.map( c => {
      users.map( u => {
        if (u.Id === c.UserId){
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

    this.comments = unsortedArray.sort( (a,b) => a.CreatedAt.getTime() - b.CreatedAt.getTime() )

  }

  commentOnTrash() {
    console.log('msg: ', this.message)
    if (this.message.length > 0) {
      console.log('msg: ', this.message)
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
}
