import {Component, OnInit} from '@angular/core';
import {SocietyWithPagingAnsw} from "../../../models/society.model";
import {SocietyService} from "../../../services/society/society.service";
import {PagingModel} from "../../../models/shared.models";
import {PageEvent} from '@angular/material/paginator';
import { MatTableModule } from '@angular/material/table';
import {SocietiesTableElementModel} from "./societiesTable.model";

@Component({
  selector: 'app-societies',
  templateUrl: './societies.component.html',
  styleUrls: ['./societies.component.css']
})
export class SocietiesComponent implements OnInit {
  actualPaging: PagingModel;
  pageEvent: PageEvent;
  displayedColumns: string[] = ['position', 'name', 'members', 'createdAt'];
  dataSource: SocietiesTableElementModel[];

  constructor(
    private societyService: SocietyService,
  ) {
    this.dataSource = [];
    this.actualPaging = {
        From: 0,
        To: 10,
        TotalCount: 10,
      }
  }


  ngOnInit(): void {
    this.societyService.getSocieties(this.actualPaging)
      .subscribe(resp => {
      this.actualPaging = resp.Paging
        this.dataSource = [];
        resp.Societies.map( (soc, i) => this.dataSource.push(
          {
            Society: soc,
            Number: this.actualPaging.From + i + 1
          }
        ))
    })
  }

  public fetchNewSocieties(event?: PageEvent) {
    this.actualPaging.From = event.pageIndex*event.pageSize
    this.actualPaging.To = (event.pageIndex*event.pageSize) + event.pageSize
    this.societyService.getSocieties(this.actualPaging)
      .subscribe(resp => {
        this.actualPaging = resp.Paging
        this.dataSource = [];
        resp.Societies.map( (soc, i) => this.dataSource.push(
          {
            Society: soc,
            Number: this.actualPaging.From + i + 1
          }
        ))
        console.log('Data source ',this.dataSource)
      })
    return event;
  }

  //createSociety()


}
