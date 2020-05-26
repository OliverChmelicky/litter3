import { Component, OnInit } from '@angular/core';
import {defaultCollectionModel,CollectionModel} from "../../models/trash.model";
import {TrashService} from "../../services/trash/trash.service";
import {ActivatedRoute, Router} from "@angular/router";
import {UserService} from "../../services/user/user.service";
import {UserModel} from "../../models/user.model";

@Component({
  selector: 'app-collection-details',
  templateUrl: './collection-details.component.html',
  styleUrls: ['./collection-details.component.css']
})
export class CollectionDetailsComponent implements OnInit {
  collection: CollectionModel = defaultCollectionModel;
  me: UserModel;
  isInCollection: boolean = false;

  constructor(
    private trashService: TrashService,
    private userService: UserService,
    private route: ActivatedRoute,
    private router: Router,
  ) { }

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      const collectionId = params.get('id');
      this.trashService.getCollectionById(collectionId).subscribe(
        collection => {
          this.collection = collection;
          this.userService.getMe().subscribe(
            me => {
              this.me = me
              this.collection.Users.map( u => {
                if (u.Id === me.Id) {
                  this.isInCollection = true;
                }
              })
            }
          )

        })
    });
  }

  onGoToTrash() {
    this.router.navigate(['trash/details', this.collection.TrashId]);
  }

  onGoToEvent() {
    this.router.navigate(['events/details', this.collection.EventId]);
  }

  onEdit() {
    this.router.navigate(['collections/edit', this.collection.Id]);
  }
}
