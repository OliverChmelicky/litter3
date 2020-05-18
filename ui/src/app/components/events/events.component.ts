import {Component, Inject, OnInit} from '@angular/core';
import {PagingModel} from "../../models/shared.models";
import {PageEvent} from '@angular/material/paginator';
import {Router} from "@angular/router";
import {EventService} from "../../services/event/event.service";
import {EventModel} from "../../models/event.model";

@Component({
  selector: 'app-events',
  templateUrl: './events.component.html',
  styleUrls: ['./events.component.css'],
})
export class EventsComponent implements OnInit {
  actualPaging: PagingModel;
  pageEvent: PageEvent;
  displayedColumns: string[] = ['date', 'num-of-attendants','remove'];
  dataSource: EventModel[];

  constructor(
    private eventService: EventService,
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
    this.eventService.getEvents(this.actualPaging)
      .subscribe(resp => {
        this.actualPaging = resp.Paging
        this.dataSource = resp.Events;
      })
  }

  public fetchNewEvents(event?: PageEvent) {
    this.actualPaging.From = event.pageIndex*event.pageSize
    this.actualPaging.To = (event.pageIndex*event.pageSize) + event.pageSize
    this.eventService.getEvents(this.actualPaging)
      .subscribe(resp => {
        this.actualPaging = resp.Paging
        this.dataSource = resp.Events;
      })
    return event;
  }

  showEventDetails(eventId: string) {
    this.router.navigate(['events/details', eventId])
  }

  onCreateEvent() {
    this.router.navigateByUrl('events/create')
  }
}

