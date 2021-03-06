import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from "@angular/router";
import {SocietyService} from "../../services/society/society.service";
import {DefaultSociety, MemberModel, SocietyModel} from "../../models/society.model";
import {loggoutUser, UserModel} from "../../models/user.model";
import {UserService} from "../../services/user/user.service";
import {EventModel, EventModelTable} from "../../models/event.model";
import {EventService} from "../../services/event/event.service";
import {requestsReceivedColumnsDefinition} from "../my-profile/table-definitions";
import {EventsColumnsDefinition, membersColumnsDefinition,} from "./table-definitions";
import {MatTableDataSource} from "@angular/material/table";
import {MarkerModel} from "../google-map/Marker.model";
import {AuthService} from "../../services/auth/auth.service";
import {MarkerCollectionModel} from "../../models/trash.model";


@Component({
  selector: 'app-society-details',
  templateUrl: './society-details.component.html',
  styleUrls: ['./society-details.component.css'],
})
export class SocietyDetailsComponent implements OnInit {
  society: SocietyModel;
  adminIds: string[] = [];
  me: UserModel = loggoutUser;

  events: EventModel[] = [];
  futureEvents: EventModelTable[] = [];
  participatedEvents: EventModelTable[] = [];
  membersIds: MemberModel[] = [];
  members: UserModel[] = [];
  applicants: UserModel[]

  requestsReceivedColumns = requestsReceivedColumnsDefinition;
  membersColumns = membersColumnsDefinition;
  participatedEventsColumns = EventsColumnsDefinition;
  upcomingEventsColumns = EventsColumnsDefinition;

  isAdmin: boolean = false;
  askedForMembership: boolean = false;
  isMember: boolean = false
  isLoggedIn: boolean = false;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private societyService: SocietyService,
    private userService: UserService,
    private eventService: EventService,
    private authService: AuthService,
  ) {
  }

  ngOnInit(): void {
    this.authService.isLoggedIn.subscribe(res => this.isLoggedIn = res)
    this.route.paramMap.subscribe(params => {
      this.societyService.getSociety(params.get('societyId')).subscribe(
        society => {
          this.society = society
          this.members = this.society.Users
          this.applicants = this.society.Applicants

          if (this.society.MemberRights) {
            this.adminIds = society.MemberRights.filter(m => m.Permission === 'admin').map(m => m.UserId)
          }

          this.userService.getMe().subscribe(
            me => {
              this.me = me
              this.adminIds.map(adminId => {
                if (adminId === this.me.Id) {
                  this.isAdmin = true
                }
              })
            },
            () => {
            },
            () => {
              if (this.society.Applicants) {
                this.society.Applicants.map(r => {
                  if (r.Id === this.me.Id) {
                    this.askedForMembership = true
                  }
                })
              }

              if (this.society.Users) {
                this.society.Users.map(mem => {
                  if (mem.Id === this.me.Id) {
                    this.isMember = true
                  }
                })
              }
            }
          )

          this.eventService.getSocietyEvents(this.society.Id).subscribe(e => this.filterEvents(e))
        }
      )
    });
  }

  filterEvents(events: EventModel[]) {
    const now = new Date().getTime()
    if (events) {
      events.map(e => {
        const eventTime = new Date(e.Date).getTime()
          let peopleAttend = 0
          if (e.UsersIds) {
            peopleAttend += e.UsersIds.length
          }
          if (e.SocietiesIds) {
            peopleAttend += e.SocietiesIds.length
          }

          if (eventTime > now) {
            this.futureEvents.push({
              id: e.Id,
              date: new Date(),
              attendingPeople: peopleAttend,
            })
          } else {
            this.participatedEvents.push({
              id: e.Id,
              date: new Date(),
              attendingPeople: peopleAttend,
            })
          }
        }
      )
      //rerender table
      let newData2 = new MatTableDataSource<EventModelTable>(this.participatedEvents);
      this.participatedEvents = []
      for (let i = 0; i < newData2.data.length; i++) {
        this.participatedEvents.push(newData2.data[i])
      }

      //rerender table
      let newData3 = new MatTableDataSource<EventModelTable>(this.futureEvents);
      this.futureEvents = []
      for (let i = 0; i < newData3.data.length; i++) {
        this.futureEvents.push(newData3.data[i])
      }
    }
  }

  onEdit() {
    this.router.navigateByUrl('societies/edit/' + this.society.Id)
  }

  onAskForMembership() {
    this.societyService.askForMembership(this.society.Id).subscribe(
      () => this.askedForMembership = true
    )
  }

  onRemoveApplication() {
    this.societyService.removeApplication(this.society.Id).subscribe(
      () => this.askedForMembership = false
    )
  }

  onLeave() {
    this.societyService.leaveSociety(this.society.Id, this.me.Id).subscribe(
      () => this.router.navigateByUrl('map')
    )
  }

  onAccept(userId: string) {
    this.societyService.acceptApplicant(this.society.Id, userId).subscribe( a => {
      this.addFromApplicantsToMembersTable(userId)
    })
  }

  onDeny(userId: string) {
    this.societyService.dismissApplicant(this.society.Id, userId).subscribe(
      () => this.removeFromRequestsTable(userId)
    )
  }

  onSeeDetails(eventId: string) {
    this.router.navigate(['events/details', eventId]);
  }

  addFromApplicantsToMembersTable(userId: string) {
    const index = this.applicants.findIndex(u => u.Id === userId)
    this.members.push(this.applicants[index])
    this.applicants.splice(index, 1)

    //rerender table
    let newData = new MatTableDataSource<UserModel>(this.applicants);
    this.applicants = []
    for (let i = 0; i < newData.data.length; i++) {
      this.applicants.push(newData.data[i])
    }

    //rerender table
    newData = new MatTableDataSource<UserModel>(this.members);
    this.members = []
    for (let i = 0; i < newData.data.length; i++) {
      this.members.push(newData.data[i])
    }
  }

  removeFromRequestsTable(userId: string) {
    const index = this.applicants.findIndex(u => u.Id === userId)
    this.applicants.splice(index, 1)

    //rerender table
    const newData = new MatTableDataSource<UserModel>(this.applicants);
    this.applicants = []
    for (let i = 0; i < newData.data.length; i++) {
      this.applicants.push(newData.data[i])
    }
  }
}
