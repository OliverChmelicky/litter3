import {Component, OnInit, ViewChild} from '@angular/core';
import {FormBuilder} from "@angular/forms";
import {ActivatedRoute, Router} from "@angular/router";
import {UserService} from "../../services/user/user.service";
import {EventModel, EventPickerModel, EventUserModel, roles} from "../../models/event.model";
import {EventService} from "../../services/event/event.service";
import {attendantsTableColumns, editAttendantsTableColumns} from "../event-details/table-definitions";
import {MatSelectChange} from "@angular/material/select";
import {AttendantsModel} from "../../models/shared.models";
import {MarkerModel} from "../google-map/Marker.model";
import {TrashService} from "../../services/trash/trash.service";
import {MatTable, MatTableDataSource} from "@angular/material/table";
import {GoogleMap} from "@agm/core/services/google-maps-types";
import {AgmMap} from "@agm/core";
import {createTrashkColumnsDefinition} from "../create-event/table-definitions";

@Component({
  selector: 'app-edit-event',
  templateUrl: './edit-event.component.html',
  styleUrls: ['./edit-event.component.css']
})
export class EditEventComponent implements OnInit {
  @ViewChild('agmMap') agmMap: AgmMap;

  event: EventModel;
  eventEditor: EventPickerModel;
  eventForm = this.formBuilder.group({
    date: new Date(),
    description: '',
    trash: [''], //modified users will be served separately
  });
  usersInEvent: EventUserModel[] = [];

  editAttendantsTableColumns = editAttendantsTableColumns;
  trashListTableColumns = createTrashkColumnsDefinition
  attendants: AttendantsModel[] = [];
  attendantsToUpdate: AttendantsModel[] = [];
  roles = roles;

  map: GoogleMap;
  initMapLat: number = 49;
  initMapLng: number = 19;
  allMarkers: MarkerModel[] = [];
  selectedMarkers: MarkerModel[] = [];
  borderTop: number;
  borderBottom: number;
  borderLeft: number;
  borderRight: number;
  initialDistance: number = 3000000

  constructor(
    private formBuilder: FormBuilder,
    private route: ActivatedRoute,
    private router: Router,
    private eventService: EventService,
    private userService: UserService,
    private trashService: TrashService,
  ) {
  }

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      this.eventService.getEvent(params.get('eventId')).subscribe(e => {
          this.event = e;
          this.usersInEvent = e.UsersIds
          this.eventForm.value['date'] = e.Date
          this.eventForm.value['description'] = e.Description
          this.eventForm.value['trash'] = e.Trash

        },
        () => {
        },
        () => {
          if (this.usersInEvent) {
            const usersIds = this.usersInEvent.map(u => u.UserId)
            this.userService.getUsersDetails(usersIds).subscribe(
              details => {
                details.map(d => {
                  this.usersInEvent.map(u => {
                    if (d.Id === u.UserId) {
                      this.attendants.push(
                        {
                          id: d.Id,
                          name: d.Email,
                          avatar: d.Avatar,
                          role: u.Permission,
                        });
                    }
                  })
                })
                //reinit attendants table
                const newData = new MatTableDataSource<AttendantsModel>(this.attendants);
                this.attendants = []
                for (let i = 0; i < newData.data.length; i++) {
                  this.attendants.push(newData.data[i])
                }
              }
            )
          }
          if (this.event.Trash) {
            const ids = this.event.Trash.map(t => t.Id)
            this.trashService.getTrashByIds(ids).subscribe(
              trash => {
                if (trash) {
                  trash.map( t => {
                    let numOfCol = 0
                    if (t.Collections) {
                      numOfCol = t.Collections.length
                    }
                    this.selectedMarkers.push({
                      id: t.Id,
                      lat: t.Location[0],
                      lng: t.Location[1],
                      cleaned: t.Cleaned,
                      images: t.Images ? t.Images : [''],
                      numOfCollections: numOfCol,
                    })

                    //reinit trash table
                    const newData = new MatTableDataSource<MarkerModel>(this.selectedMarkers);
                    this.selectedMarkers = []
                    for (let i = 0; i < newData.data.length; i++) {
                      this.selectedMarkers.push(newData.data[i])
                    }
                  })
                }
              }
            )
          }
        }
      );
    })
    this.eventEditor = this.eventService.getEventEditor()
  }

  memberPermissionChange(event: MatSelectChange, i: any) {
    if (event.value === this.attendants[i].role) {
      const index = this.attendantsToUpdate.findIndex(u => u.id === this.attendants[i].id)
      this.attendantsToUpdate.splice(index, 1)
    } else {
      const exists = this.attendantsToUpdate.filter(mem => mem.id === this.attendants[i].id)
      if (exists.length !== 0) {
        const index = this.attendantsToUpdate.findIndex(u => u.id === this.attendants[i].id)
        this.attendantsToUpdate.splice(index, 1)
      }

      this.attendantsToUpdate.push({
        id: this.attendants[i].id,
        role: event.value.toString(),
        avatar: '',
        name: '',
      })
    }
  }

  onAttendantsPermissionAcceptChanges() {
    if (this.attendantsToUpdate) {
      this.attendantsToUpdate.map(a => {
        this.eventService.updateUserPermission(a, this.eventEditor, this.event.Id).subscribe(
          () => {
          }
        )
      })
    }
  }

  onDeleteEvent() {
    this.attendants.map(a => {
      if (a.id === this.eventEditor.Id && a.role === 'admin') {
        this.eventService.deleteEvent(this.eventEditor, this.event.Id).subscribe(
          () => {
          }
        )
      }
    })
  }

  addToList(marker: MarkerModel) {
    this.selectedMarkers.push(marker)

    const index = this.allMarkers.findIndex(t => t.id === marker.id)
    this.selectedMarkers = this.allMarkers.splice(index, 1)
  }

  removeFromList(trashId: string) {
    const index = this.selectedMarkers.findIndex(t => t.id === trashId)
    this.allMarkers.push(this.selectedMarkers[index])
    this.selectedMarkers.splice(index, 1)

    //rerender table
    const newData = new MatTableDataSource<MarkerModel>(this.selectedMarkers);
    this.selectedMarkers = []
    for (let i = 0; i < newData.data.length; i++) {
      this.selectedMarkers.push(newData.data[i])
    }

  }

  navigateToTrash(id: any) {
    this.router.navigate(['trash/details', id])
  }

  async onMapReady(map: GoogleMap) {
    this.map = map;
    //In an issue it was written that this helps but don`t
    await setTimeout(() => {
      this.agmMap.triggerResize();
    }, 1000)
    let c = this.map.getCenter()
    this.borderTop = c.lat() + 3.4
    this.borderBottom = c.lat() - 3.4

    this.borderRight = c.lng() + 8.82
    this.borderLeft = c.lng() - 8.82

    this.trashService.getTrashInRange(this.map.getCenter().lat(), this.map.getCenter().lng(), this.initialDistance).subscribe(
      trash => {
        for (let i = 0; i < trash.length; i++) {
          this.allMarkers.push({
            lat: trash[i].Location[0],
            lng: trash[i].Location[1],
            new: false,
            id: trash[i].Id,
            cleaned: trash[i].Cleaned,
            images: trash[i].Images ? trash[i].Images : [],
            numOfCollections: trash[i].Collections ? trash[i].Collections.length : 0
          })

          this.filterSelected()
        }
      })
  }

  onBoundsChange() {
    const p1 = this.map.getBounds().getNorthEast()
    const p2 = this.map.getBounds().getSouthWest()

    const visibleTop = p1.lat()
    const visibleRight = p1.lng()
    const visibleBottom = p2.lat()
    const visibleLeft = p2.lng()

    if (visibleRight > this.borderRight || visibleLeft < this.borderLeft) {
      this.loadNewMarkers()
    } else if (visibleBottom < this.borderBottom || visibleTop > this.borderTop) {
      this.loadNewMarkers()
    }

  }

  loadNewMarkers() {
    const p1 = this.map.getBounds().getNorthEast()
    const p2 = this.map.getBounds().getSouthWest()

    const R = 6371e3; // metres
    const fi1 = p1.lat() * Math.PI / 180; // φ, λ in radians
    const fi2 = p2.lat() * Math.PI / 180;
    const delta1 = (p2.lat() - p1.lat()) * Math.PI / 180;
    const delta2 = (p2.lng() - p1.lng()) * Math.PI / 180;

    const a = Math.sin(delta1 / 2) * Math.sin(delta1 / 2) +
      Math.cos(fi1) * Math.cos(fi2) *
      Math.sin(delta2 / 2) * Math.sin(delta2 / 2);
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));

    const d = R * c; // in metres

    //get double range for markers
    this.trashService.getTrashInRange(this.map.getCenter().lat(), this.map.getCenter().lng(), d * 2).subscribe(
      trash => {
        this.allMarkers = [];
        for (let i = 0; i < trash.length; i++) {
          this.allMarkers.push({
            lat: trash[i].Location[0],
            lng: trash[i].Location[1],
            new: false,
            id: trash[i].Id,
            cleaned: trash[i].Cleaned,
            images: trash[i].Images ? trash[i].Images : [],
            numOfCollections: trash[i].Collections ? trash[i].Collections.length : 0
          })

          this.filterSelected
        }

        const viewCenter = this.map.getCenter()
        let r = 2 * Math.abs(p1.lat() - viewCenter.lat())
        console.log('R: ', r)

        if (p1.lat() < 0) {
          this.borderTop =  p1.lat() + r
        } else if (p1.lat() >= 0) {
          this.borderTop = p1.lat() + r
        }
        if (p1.lng() < 0) {
          this.borderRight = p1.lng() + r
        } else if (p1.lng() >= 0) {
          this.borderRight = p1.lng() + r
        }

        if (p2.lat() < 0) {
          this.borderBottom = p2.lat() - r
        } else if (p2.lat() >= 0) {
          this.borderBottom = p2.lat() - r
        }
        if (p2.lng() < 0) {
          this.borderLeft = p2.lng() - r
        } else if (p2.lng() >= 0) {
          this.borderLeft = p2.lng() - r
        }

      }
    )
  }

  private filterSelected() {
    this.allMarkers = this.allMarkers.filter( m => {
      this.selectedMarkers.map( selected => {
        if (m.id !== selected.id) {
          return m
        }
      })
    })
  }
}
