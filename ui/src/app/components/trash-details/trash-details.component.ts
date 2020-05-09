import { Component, OnInit } from '@angular/core';
import {ActivatedRoute} from "@angular/router";
import {TrashService} from "../../services/trash/trash.service";
import {TrashModel} from "../../models/trash.model";

@Component({
  selector: 'app-trash-details',
  templateUrl: './trash-details.component.html',
  styleUrls: ['./trash-details.component.css']
})
export class TrashDetailsComponent implements OnInit {
  trashId: string;
  trash: TrashModel;

  constructor(
    private route: ActivatedRoute,
    private trashService: TrashService,
  ) { }

  ngOnInit(): void {
      this.route.paramMap.subscribe(params => {
        this.trashId = params.get('id');
        this.trashService.getTrashById(this.trashId).subscribe(
          trash => this.trash = trash
        )
        });
  }

}
