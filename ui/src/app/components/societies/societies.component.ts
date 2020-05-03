import {Component, OnInit} from '@angular/core';
import {SocietyService} from "../../services/society/society.service";
import {PagingModel} from "../../models/shared.models";
import {PageEvent} from '@angular/material/paginator';
import {SocietiesTableElementModel} from "./societiesTable.model";
import {animate, state, style, transition, trigger} from '@angular/animations';
import {Router} from "@angular/router";

@Component({
  selector: 'app-societies',
  templateUrl: './societies.component.html',
  styleUrls: ['./societies.component.css'],
  animations: [
    trigger('detailExpand', [
      state('collapsed', style({height: '0px', minHeight: '0'})),
      state('expanded', style({height: '*'})),
      transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
    ]),
  ],
})
export class SocietiesComponent implements OnInit {
  actualPaging: PagingModel;
  pageEvent: PageEvent;
  displayedColumns: string[] = ['position', 'name', 'members', 'createdAt'];
  dataSource: SocietiesTableElementModel[];
  expandedElement: SocietiesTableElementModel | null;

  constructor(
    private societyService: SocietyService,
    private router: Router,
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

  showSocietyDetails(Id: string) {
    this.router.navigateByUrl('societies/'+Id)
  }

  createSociety() {

  }
}
