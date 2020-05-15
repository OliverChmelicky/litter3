import {Component, OnInit} from '@angular/core';
import {FormBuilder, FormControl} from "@angular/forms";
import {MatDatepickerInputEvent} from "@angular/material/datepicker";
import {EventService} from "../../services/event/event.service";
import {EventCreatorModel, EventModel} from "../../models/event.model";
import {UserModel} from "../../models/user.model";
import {SocietyModel} from "../../models/society.model";
import {UserService} from "../../services/user/user.service";
import {SocietyService} from "../../services/society/society.service";

@Component({
  selector: 'app-create-event',
  templateUrl: './create-event.component.html',
  styleUrls: ['./create-event.component.css']
})
export class CreateEventComponent implements OnInit {
  me: UserModel;
  availableCreators: EventCreatorModel[] = [];
  selectedCreator: number;
  newEvent: EventModel = {
    Date: new Date(),
    Description: '',
  };
  description: string;
  date = new FormControl(new Date());

  constructor(
    private formBuilder: FormBuilder,
    private eventService: EventService,
    private userService: UserService,
  ) {
  }

  ngOnInit(): void {
    this.userService.getMe().subscribe(
      me => {
        this.me = me
        this.availableCreators.push({
          VisibleName: me.Email,
          Id: me.Id,
          AsSociety: false
        })
        this.userService.getMyEditableSocieties().subscribe(
          editable => {
            console.log(editable)
            if (editable) {
              editable.map(soc => this.availableCreators.push({
                VisibleName: soc.Name,
                Id: soc.Id,
                AsSociety: true
              }))
            }
          }
        )
      }
    )

  }

  onSubmit() {
    this.newEvent.Date = this.date.value
    this.newEvent.Description = this.description
    console.log(this.newEvent)
  }

}
