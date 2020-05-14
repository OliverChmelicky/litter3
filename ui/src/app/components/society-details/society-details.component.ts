import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from "@angular/router";
import {SocietyService} from "../../services/society/society.service";
import {ApplicantModel, MemberModel, SocietyModel} from "../../models/society.model";
import {UserModel} from "../../models/user.model";
import {UserService} from "../../services/user/user.service";
import {EventModel} from "../../models/event.model";
import {EventService} from "../../services/event/event.service";
import {
  friendsColumnsDefinition, requestsReceivedColumnsDefinition,
  requestsSendColumnsDefinition,
  societiesColumnsDefinition
} from "../my-profile/table-definitions";
import {
  membersColumnsDefinition,
  participatedEventsColumnsDefinition,
  upcommingEventsColumnsDefinition
} from "./table-definitions";
import {ApisModel} from "../../api/api-urls";

@Component({
  selector: 'app-society-details',
  templateUrl: './society-details.component.html',
  styleUrls: ['./society-details.component.css']
})
export class SocietyDetailsComponent implements OnInit {
  society: SocietyModel;
  adminIds: string[];
  me: UserModel;
  isAdmin: boolean = false;
  events: EventModel[];
  futureEvents: EventModel[];
  participatedEvents: EventModel[];
  membersIds: MemberModel[];
  members: UserModel[];
  societyRequests: ApplicantModel[];

  upcommingEventsColumns = upcommingEventsColumnsDefinition;
  requestsReceivedColumns = requestsReceivedColumnsDefinition;
  membersColumns = membersColumnsDefinition;
  participatedEventsColumns = participatedEventsColumnsDefinition;

  requesters: UserModel[]

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private societyService: SocietyService,
    private userService: UserService,
    private eventService: EventService,
  ) {
  }

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      this.societyService.getSociety(params.get('societyId')).subscribe(
        society => {
          if (society.Avatar) {
            society.Avatar = ApisModel.pictureBucketPrefix + society.Avatar
          }
          this.society = society
          this.societyService.getSocietyAdmins(society.Id).subscribe(
            res => {
              this.adminIds = res
              this.userService.getMe().subscribe(
                res => {
                  this.me = res
                  this.adminIds.map(adminId => {
                    if (adminId === this.me.Id) {
                      this.isAdmin = true
                    }
                  })
                }
              )

              //getSocietyRequests
              this.societyService.getSocietyRequests(this.society.Id).subscribe(requests => {
                this.societyRequests = requests;
                this.getRequesters(requests)
              })
              //getSocietyMembersDetails
              this.societyService.getSocietyMembers(this.society.Id).subscribe(m => {
                this.membersIds = m
                this.getMembers(m)
              })
              //get society events
              this.eventService.getSocietyEvents(this.society.Id).subscribe(e => this.filterEvents(e))


            }
          )
        }
      )
    });
  }

  filterEvents(events: EventModel[]) {
    const now = new Date()
    events.map(e => {
        if (e.Date > now) {
          this.futureEvents.push(e)
        } else {
          this.participatedEvents.push(e)
        }
      }
    )
  }

  onEdit() {
    this.router.navigateByUrl('societies/edit/' + this.society.Id)
  }

  private getRequesters(requests: ApplicantModel[]) {
    const requestersIds = requests.map( r => r.UserId)
    console.log(requestersIds)
    this.userService.getUsersDetails(requestersIds).subscribe( r => this.requesters = r)
  }

  private getMembers(membersIds: MemberModel[]) {
    const membIds = membersIds.map( m => m.UserId)
    this.userService.getUsersDetails(membIds).subscribe( m => this.members = m)
  }
}
