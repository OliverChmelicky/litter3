import { Component, OnInit } from '@angular/core';
import {UserModel} from "../../models/user.model";
import {EventModel, EventPickerModel} from "../../models/event.model";
import {attendantsTableColumns} from "./table-definitions";
import {AttendantsModel} from "../../models/shared.models";
import {ActivatedRoute, Router} from "@angular/router";
import {UserService} from "../../services/user/user.service";
import {EventService} from "../../services/event/event.service";
import {SocietyService} from "../../services/society/society.service";
import {SocietyModel} from "../../models/society.model";

@Component({
  selector: 'app-event-details',
  templateUrl: './event-details.component.html',
  styleUrls: ['./event-details.component.css']
})
export class EventDetailsComponent implements OnInit {
  isLoggedIn: boolean = false;
  statusAttend: boolean = false;
  me: UserModel;
  event: EventModel;
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
          )
        },
        () => this.isLoggedIn = false,
        () => {
          this.eventService.getEvent(params.get('eventId')).subscribe( event => {
            this.event = event
            if (event.UsersIds) {
              event.UsersIds.map( attendant => {
                if (attendant.Permission === 'creator' && attendant.UserId === this.me.Id) {
                  this.isAdmin = true
                }
              })
            }

            if (!this.isAdmin && event.SocietiesIds && this.editableSocieties.length > 0) { //is society admin and I have society rights editor and more so I am admin
              event.SocietiesIds.map( attendant => {
                if (attendant.Permission === 'creator') {
                  if (this.editableSocietiesIds.includes(attendant.SocietyId)){
                    this.isAdmin = true
                  }
                }
              })
            }

            this.getMyEventAttendance()
          })
        }
      )

    })
  }

  onCreateCollections() {
    this.router.navigateByUrl('/collections/event')
  }

  desideAsChange() {
    if (this.selectedCreator === 0){  //User
      this.getMyEventAttendance()
    } else {
      this.getSocietyEventAttendance()
    }
  }

  private getMyEventAttendance() {
    this.statusAttend = false
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

  private getSocietyEventAttendance(){
    this.statusAttend = false
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

  onAttend() {
    this.eventService.attendEvent(this.event.Id, this.availableDecisionsAs[this.selectedCreator])
  }

  onNotAttend() {
    this.eventService.notAttendEvent(this.event.Id, this.availableDecisionsAs[this.selectedCreator])
  }

  onEdit() {
    this.router.navigate(['event/details', this.event.Id])
  }
}
