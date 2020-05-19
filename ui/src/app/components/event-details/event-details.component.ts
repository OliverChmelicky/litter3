import { Component, OnInit } from '@angular/core';
import {UserModel} from "../../models/user.model";
import {EventModel, EventPickerModel, EventSocietyModel, EventUserModel} from "../../models/event.model";
import {attendantsTableColumns} from "./table-definitions";
import {AttendantsModel} from "../../models/shared.models";
import {ActivatedRoute, Router} from "@angular/router";
import {UserService} from "../../services/user/user.service";
import {EventService} from "../../services/event/event.service";
import {SocietyService} from "../../services/society/society.service";
import {SocietyModel} from "../../models/society.model";
import {MatTableDataSource} from "@angular/material/table";

@Component({
  selector: 'app-event-details',
  templateUrl: './event-details.component.html',
  styleUrls: ['./event-details.component.css']
})
export class EventDetailsComponent implements OnInit {
  isLoggedIn: boolean = false;
  statusAttend: boolean = false;
  me: UserModel;
  event: EventModel = {
    Date: new Date,
    Description: '',
  };
  attendants: AttendantsModel[] = [];
  tableColumns = attendantsTableColumns;
  availableDecisionsAs: EventPickerModel[] = [];
  selectedCreator: number = 0;

  isAdmin: boolean = false;
  editableSocieties: SocietyModel[] = [];
  editableSocietiesIds: string[] = [];



  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private userService: UserService,
    private societyService: SocietyService,
    private eventService: EventService,
  ) { }

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      this.userService.getMe().subscribe(
        me => {
          this.me = me;
          this.availableDecisionsAs.push({
            VisibleName: me.Email,
            Id: me.Id,
            AsSociety: false
          })
          this.isLoggedIn = true;
          this.userService.getMyEditableSocieties().subscribe(
            societies => {
              if (societies) {
                this.editableSocieties = societies
                this.editableSocietiesIds = societies.map( soc => {
                  this.availableDecisionsAs.push({
                    VisibleName: soc.Name,
                    Id: soc.Id,
                    AsSociety: true
                  })
                  return soc.Id
                })
              }
            }
          )
        },
        () => this.isLoggedIn = false,
        () => {
          this.eventService.getEvent(params.get('eventId')).subscribe( event => {
            this.convertToLocalTime()
            this.event = event
            if (event.UsersIds) {
              event.UsersIds.map( attendant => {
                if (attendant.Permission === 'creator' && attendant.UserId === this.me.Id) {
                  this.isAdmin = true
                }
              })
            }

            //continiue if not user admin and
            //is society admin and I have society rights editor and more so I am admin
            if (!this.isAdmin && event.SocietiesIds && this.editableSocieties.length > 0) {
              event.SocietiesIds.map( attendant => {
                if (attendant.Permission === 'creator') {
                  if (this.editableSocietiesIds.includes(attendant.SocietyId)){
                    this.isAdmin = true
                  }
                }
              })
            }

            if (event.UsersIds) {
              const userIds = event.UsersIds.map(u => u.UserId)
              this.fetchUsersWhoAttends(userIds, event.UsersIds)
            }

            if (event.SocietiesIds){
              const societyIds = event.SocietiesIds.map(s => s.SocietyId)
              this.fetchSocietiesWhichAttend(societyIds, event.SocietiesIds)
            }

            this.getEventAttendanceOnUser()
          })
        }
      )

    })
  }

  onCreateCollections() {
    this.router.navigateByUrl('/collections/event')
  }

  onDesideAsChange() {
    if (this.selectedCreator === 0){  //User
      this.getEventAttendanceOnUser()
    } else {
      this.getSocietyEventAttendance()
    }
  }

  private getEventAttendanceOnUser() {
    this.statusAttend = false
    this.isAdmin = false
    if (this.event.UsersIds) {
      this.event.UsersIds.map(
        user => {
          if (user.UserId == this.me.Id) {
            this.statusAttend = true
            if (user.Permission === 'admin') {
              this.isAdmin = true
            }
          }
        }
      )
    }
  }

  private getSocietyEventAttendance(){
    this.statusAttend = false
    this.isAdmin = false

    if (this.event.SocietiesIds) {
      this.event.SocietiesIds.map(
        society => {
          if (society.SocietyId == this.availableDecisionsAs[this.selectedCreator].Id) {
            this.statusAttend = true
            if (society.Permission === 'admin') {
              this.isAdmin = true
            }
          }
        }
      )
    }
  }

  onAttend() {
    this.eventService.attendEvent(this.event.Id, this.availableDecisionsAs[this.selectedCreator]).subscribe(
      () => {}
    )
  }

  onNotAttend() {
    this.eventService.notAttendEvent(this.event.Id, this.availableDecisionsAs[this.selectedCreator])
  }

  onEdit() {
    this.router.navigate(['event/details', this.event.Id])
  }

  private fetchUsersWhoAttends(userIds: string[], userEventDetails: EventUserModel[]) {
    this.userService.getUsersDetails(userIds).subscribe(
      users => {
        users.map( u => {
          userEventDetails.map( detail => {
            if (u.Id === detail.UserId) {
              this.attendants.push({
                name: u.Email,
                avatar: u.Avatar? u.Avatar : '',
                role: detail.Permission,
              })
            }
            const newData = new MatTableDataSource<AttendantsModel>(this.attendants);
            this.attendants = []
            for (let i = 0; i < newData.data.length; i++) {
              this.attendants.push(newData.data[i])
            }
          });
        })
      })
  }

  private fetchSocietiesWhichAttend(societiesIds: string[], societyEventDetails: EventSocietyModel[]) {
    this.societyService.getSocietiesByIds(societiesIds).subscribe(
      societies => {
        societies.map( s => {
          societyEventDetails.map( detail => {
            if (s.Id === detail.SocietyId) {
              this.attendants.push({
                name: s.Name,
                avatar: s.Avatar? s.Avatar : '',
                role: detail.Permission,
              })
            }
          });
        })
        const newData = new MatTableDataSource<AttendantsModel>(this.attendants);
        this.attendants = []
        for (let i = 0; i < newData.data.length; i++) {
          this.attendants.push(newData.data[i])
        }
      })

  }

  private convertToLocalTime() {
    let dateStr = this.event.Date.toString()
    dateStr += ' UTC'
    console.log('old date: ', dateStr)
    this.event.Date = new Date(dateStr)
    console.log('new date: ', this.event.Date)
  }
}
