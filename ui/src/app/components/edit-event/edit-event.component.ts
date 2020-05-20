import { Component, OnInit } from '@angular/core';
import {FormBuilder} from "@angular/forms";
import {ActivatedRoute, Router} from "@angular/router";
import {UserService} from "../../services/user/user.service";
import {EventModel, EventPickerModel} from "../../models/event.model";
import {EventService} from "../../services/event/event.service";
import {attendantsTableColumns} from "../event-details/table-definitions";
import {MatSelectChange} from "@angular/material/select";

@Component({
  selector: 'app-edit-event',
  templateUrl: './edit-event.component.html',
  styleUrls: ['./edit-event.component.css']
})
export class EditEventComponent implements OnInit {
  event: EventModel;
  eventEditor: EventPickerModel;
  eventForm = this.formBuilder.group({
    date: new Date(),
    description: '',
    trashIds: [''], //modified users will be served separately
  });
  tableColumns = attendantsTableColumns;

  constructor(
    private formBuilder: FormBuilder,
    private route: ActivatedRoute,
    private router: Router,
    private eventService: EventService,
    private userService: UserService,
  ) { }

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      this.eventService.getEvent(params.get('eventId')).subscribe( e => {
        this.event = e;
      } )
    })

    this.eventEditor = this.eventService.getEventEditor()


  }

  memberPermissionChange(event: MatSelectChange, i: any) {
    // if (event.value === this.origMembers[i].role) {
    //   const index = this.changeMemberPermission.findIndex(u => u.UserId === this.members[i].user.Id)
    //   this.changeMemberPermission.splice(index, 1)
    // } else {
    //   const exists = this.changeMemberPermission.filter( mem => mem.UserId === this.members[i].user.Id)
    //   if (exists.length !== 0) {
    //     const index = this.changeMemberPermission.findIndex(u => u.UserId === this.members[i].user.Id)
    //     this.changeMemberPermission.splice(index, 1)
    //   }
    //
    //   this.changeMemberPermission.push({
    //     UserId: this.members[i].user.Id,
    //     SocietyId: this.society.Id,
    //     Permission: event.value.toString(),
    //     CreatedAt: new Date(),  //server does not use this property
    //   })
    // }
  }

  onAttendantsPermissionAcceptChanges() {

  }
}
