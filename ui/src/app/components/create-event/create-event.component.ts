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
  selectedCreator: number = 0;
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
            console.log(this.availableCreators)
          }
        )
      }
    )

  }

  onSubmit() {
    this.newEvent.Date = this.date.value
    this.newEvent.Description = this.description
    const request = {
      UserId: this.me.Id,
      SocietyId: this.availableCreators[this.selectedCreator].Id,
      AsSociety: this.availableCreators[this.selectedCreator].AsSociety,
      Description: this.description,
      Date: this.date.value,
      Trash: [],
    }
    console.log('new event je')
    console.log(request)
  }



}
