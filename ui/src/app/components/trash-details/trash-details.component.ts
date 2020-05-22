import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from "@angular/router";
import {TrashService} from "../../services/trash/trash.service";
import {TrashModel, TrashTypeBooleanValues} from "../../models/trash.model";
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

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private trashService: TrashService,
    private authService: AuthService,
    private userService: UserService,
  ) {
  }

  ngOnInit(): void {
    this.authService.isLoggedIn.subscribe( isLogged => this.isLoggedIn = isLogged)

    this.route.paramMap.subscribe(params => {
      this.trashId = params.get('id');
      this.trashService.getTrashById(this.trashId).subscribe(
        trash => {
          console.log(trash)

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
          this.trash = trash
          this.trashTypeBool = this.trashService.convertTrashTypeNumToBools(this.trash.TrashType);
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
    this.router.navigateByUrl('collection/'+this.trash.Id)
  }

  onCreateEvent() {
    this.router.navigateByUrl('events/create')
  }
}
