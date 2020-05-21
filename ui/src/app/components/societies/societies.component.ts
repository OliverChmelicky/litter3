import {Component, Inject, OnInit} from '@angular/core';
import {SocietyService} from "../../services/society/society.service";
import {PagingModel} from "../../models/shared.models";
import {PageEvent} from '@angular/material/paginator';
import {SocietiesTableElementModel} from "./societiesTable.model";
import {animate, state, style, transition, trigger} from '@angular/animations';
import {Router} from "@angular/router";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material/dialog";
import {ApisModel} from "../../api/api-urls";
import {AuthService} from "../../services/auth/auth.service";

export interface DialogData {
  name: string;
  description: string;
}

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
  displayedColumns: string[] = ['position', 'avatar','name', 'members', 'createdAt', 'showMore'];
  dataSource: SocietiesTableElementModel[];
  expandedElement: SocietiesTableElementModel | null;
  private isLoggedIn: boolean = false;

  constructor(
    private societyService: SocietyService,
    private router: Router,
    public createSocietyDialog: MatDialog,
    private authService: AuthService,
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
        resp.Societies.map( (soc, i) => {
          this.dataSource.push(
            {
              Society: soc,
              Number: this.actualPaging.From + i + 1
            }
          )
        })
    })
    this.authService.isLoggedIn.subscribe( res => this.isLoggedIn = res)
  }

  openDialog(): void {
    const dialogRef = this.createSocietyDialog.open(CreateSocietyComponent, {
      width: '800px',
      data: {
        name: '',
        description: '',
      }
    });

    dialogRef.afterClosed().subscribe(result => {
      console.log(result)
      if (result) {
        if (result.name != '') {
          this.societyService.createSociety({
            Name: result.name,
            Description: result.description,
          }).subscribe(
            newSoc => {
              this.router.navigate(['societies', newSoc.Id])
            }
          )
        }
      }
    });
  }

  public fetchNewSocieties(event?: PageEvent) {
    this.actualPaging.From = event.pageIndex*event.pageSize
    this.actualPaging.To = (event.pageIndex*event.pageSize) + event.pageSize
    this.societyService.getSocieties(this.actualPaging)
      .subscribe(resp => {
        this.actualPaging = resp.Paging
        this.dataSource = [];
        resp.Societies.map( (soc, i) => {
          soc.Avatar = ApisModel.pictureBucketPrefix + soc.Avatar
          console.log(soc.Avatar)
          this.dataSource.push(
          {
            Society: soc,
            Number: this.actualPaging.From + i + 1
          }
        )
        })
        console.log(resp.Societies)
      })
    return event;
  }

  showSocietyDetails(Id: string) {
    this.router.navigateByUrl('societies/'+Id)
  }

  createSociety() {

  }
}


  @Component({
    selector: 'app-edit-profile',
    templateUrl: './dialog/create-society.component.html',
    //styleUrls: ['./dialog/edit-create-society.component.css]
  })
  export class CreateSocietyComponent {

  constructor( public dialogRef: MatDialogRef<CreateSocietyComponent>,
               @Inject(MAT_DIALOG_DATA) public data: DialogData) {
  }

  onNoClick(): void {
    this.dialogRef.close();
  }

}
