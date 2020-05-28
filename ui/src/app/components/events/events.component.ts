import {Component, OnInit} from '@angular/core';
import {PagingModel} from "../../models/shared.models";
import {PageEvent} from '@angular/material/paginator';
import {Router} from "@angular/router";
import {EventService} from "../../services/event/event.service";
import {EventModel, ListEventsModel} from "../../models/event.model";
import {AuthService} from "../../services/auth/auth.service";

@Component({
  selector: 'app-events',
  templateUrl: './events.component.html',
  styleUrls: ['./events.component.css'],
})
export class EventsComponent implements OnInit {
  actualPaging: PagingModel;
  isLoggedIn: boolean;
  pageEvent: PageEvent;
  displayedColumns: string[] = ['date', 'num-of-attendants','remove'];
  dataSource: ListEventsModel[];

  constructor(
    private eventService: EventService,
    private router: Router,
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
    this.authService.isLoggedIn.subscribe( res => this.isLoggedIn = res)

    this.eventService.getEvents(this.actualPaging)
      .subscribe(resp => {
        this.actualPaging = resp.Paging
        this.dataSource = this.mapResponsToDataSource(resp.Events);
      })
  }

  public fetchNewEvents(event?: PageEvent) {
    this.actualPaging.From = event.pageIndex*event.pageSize
    this.actualPaging.To = (event.pageIndex*event.pageSize) + event.pageSize
    this.eventService.getEvents(this.actualPaging)
      .subscribe(resp => {
        this.actualPaging = resp.Paging
        this.dataSource = this.mapResponsToDataSource(resp.Events);
      })
    return event;
  }

  showEventDetails(eventId: string) {
    this.router.navigate(['events/details', eventId])
  }

  onCreateEvent() {
    this.router.navigateByUrl('events/create')
  }

  private mapResponsToDataSource(events: EventModel[]): ListEventsModel[] {
    return events.map(e => {
      let attNum = 0
      if (e.SocietiesIds) {
        attNum = e.SocietiesIds.length
      }
      if (e.UsersIds) {
        attNum = attNum + e.UsersIds.length
      }
      return {
        Id: e.Id,
        Date: e.Date,
        NumOfAttendants: attNum,
      }
    })
  }

  private convertToLocalTime(events: EventModel[]): EventModel[] {
     return events.map( e => {
      let dateStr = e.Date.toString()
      dateStr += ' UTC'
      e.Date = new Date(dateStr)
      return e
    })
  }
}

