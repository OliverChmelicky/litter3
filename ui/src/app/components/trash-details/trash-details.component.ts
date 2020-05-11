import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from "@angular/router";
import {TrashService} from "../../services/trash/trash.service";
import {TrashModel} from "../../models/trash.model";
import {GoogleMap} from "@agm/core/services/google-maps-types";
import {accessibilityChoces} from "../../models/accessibilityChocies";
import {MatSliderChange} from "@angular/material/slider";

@Component({
  selector: 'app-trash-details',
  templateUrl: './trash-details.component.html',
  styleUrls: ['./trash-details.component.css']
})
export class TrashDetailsComponent implements OnInit {
  map: GoogleMap;
  trashId: string;
  trash: TrashModel;
  sizeView: string;

  accessibilityChoices = accessibilityChoces;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private trashService: TrashService,
  ) {
  }

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      this.trashId = params.get('id');
      this.trashService.getTrashById(this.trashId).subscribe(
        trash => this.trash = trash
      )
    });
  }

  onMapReady(map: GoogleMap) {
    this.map = map;
  }

  printSize(event: number) {
    if (event === 0) {
      return 'unknown';
    }
    if (event === 1) {
      return 'bag';
    }
    if (event === 2) {
      return 'wheelbarrow';
    }
    if (event === 3) {
      return 'car';
    }
  }

  onEdit() {
    this.router.navigateByUrl('trash/edit/'+this.trash.Id)
  }
}
