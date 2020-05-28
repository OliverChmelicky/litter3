import {Component, OnInit, ViewChild} from '@angular/core';
import {FormBuilder, FormControl} from "@angular/forms";
import {EventService} from "../../services/event/event.service";
import {EventModel, EventPickerModel, EventRequestModel} from "../../models/event.model";
import {UserModel} from "../../models/user.model";
import {UserService} from "../../services/user/user.service";
import {LocationService} from "../../services/location/location.service";
import {ActivatedRoute, Router} from "@angular/router";
import {GoogleMap} from "@agm/core/services/google-maps-types";
import {MarkerModel} from "../google-map/Marker.model";
import {ApisModel} from "../../api/api-urls";
import {TrashService} from "../../services/trash/trash.service";
import {createTrashkColumnsDefinition} from "./table-definitions";
import {AgmMap} from "@agm/core";
import {MatTable, MatTableDataSource} from "@angular/material/table";
import {defaultTrashImage} from "../../models/trash.model";
import {initialDistance} from "../../models/shared.models";

@Component({
  selector: 'app-create-event',
  templateUrl: './create-event.component.html',
  styleUrls: ['./create-event.component.css']
})
export class CreateEventComponent implements OnInit {
  @ViewChild('agmMap') agmMap: AgmMap;
  @ViewChild('table') table: MatTable<any>;

  allMarkers: MarkerModel[] = [];
  selectedTrash: MarkerModel[] = [];
  tableColumns = createTrashkColumnsDefinition
  borderTop: number;
  borderBottom: number;
  borderLeft: number;
  borderRight: number;

  map: GoogleMap;
  initMapLat: number;
  initMapLng: number;
  me: UserModel;
  availableCreators: EventPickerModel[] = [];
  selectedCreator: number = 0;
  newEvent: EventModel = {
    Date: new Date(),
    Description: '',
  };
  description: string;
  date = new FormControl(new Date());
  private initialDistance: number = initialDistance;

  constructor(
    private formBuilder: FormBuilder,
    private eventService: EventService,
    private userService: UserService,
    private locationService: LocationService,
    private trashService: TrashService,
    private route: ActivatedRoute,
    private router: Router,
  ) {
  }

  ngOnInit(): void {
    //get location
    this.route.paramMap.subscribe(params => {
      this.initMapLat = +params.get('lat');
      this.initMapLng = +params.get('lng');

      if (this.initMapLat === 0 && this.initMapLng === 0) {
        this.locationService.getPosition().then(data => {
          this.initMapLat = data.lat;
          this.initMapLng = data.lng;
        }).catch(
          () => {
            this.initMapLat = 49;
            this.initMapLng = 16;
          }
        );
      }
    });

    //get my creator information
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

  ngAfterViewInit(){
    this.agmMap.mapReady.subscribe(map => {
      this.map = map

      let c = this.map.getCenter()
      this.borderTop = c.lat() + 3.4
      this.borderBottom = c.lat() - 3.4

      this.borderRight = c.lng() + 8.82
      this.borderLeft = c.lng() - 8.82

      this.trashService.getTrashInRange(this.map.getCenter().lat(), this.map.getCenter().lng(), this.initialDistance).subscribe(
        trash => {
          for (let i = 0; i < trash.length; i++) {
            if (!trash[i].Images) {
              trash[i].Images = [];
            }
            this.allMarkers.push({
              lat: trash[i].Location[0],
              lng: trash[i].Location[1],
              new: false,
              id: trash[i].Id,
              cleaned: trash[i].Cleaned,
              images: trash[i].Images,
              numOfCollections: trash[i].Collections ? trash[i].Collections.length : 0
            })

            this.allMarkers = this.trashService.filterCleanedAndSelected(this.allMarkers, this.selectedTrash)
          }
        })

    });
  }

  onSubmit() {
    this.newEvent.Date = this.date.value
    this.newEvent.Description = this.description
    const trashIds = this.selectedTrash.map(t => t.id)
    const request: EventRequestModel = {
      UserId: this.me.Id,
      SocietyId: this.availableCreators[this.selectedCreator].Id,
      AsSociety: this.availableCreators[this.selectedCreator].AsSociety,
      Description: this.description,
      Date: this.date.value,
      Trash: trashIds,
    }
    this.eventService.createEvent(request).subscribe(
      e => {
        this.router.navigateByUrl('events')
      }
    )
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
            images: trash[i].Images ? trash[i].Images : [defaultTrashImage],
            numOfCollections: trash[i].Collections ? trash[i].Collections.length : 0
          })

          this.allMarkers = this.trashService.filterCleanedAndSelected(this.allMarkers, this.selectedTrash)
        }

        const viewCenter = this.map.getCenter()
        let r = 2 * Math.abs(p1.lat() - viewCenter.lat())

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

  addToList(marker: MarkerModel) {
    this.selectedTrash.push(marker)

    const index = this.allMarkers.findIndex(t => t.id === marker.id)
    this.allMarkers.splice(index, 1)

    //rerender table
    const newData = new MatTableDataSource<MarkerModel>(this.selectedTrash);
    this.selectedTrash = []
    for (let i = 0; i < newData.data.length; i++) {
      this.selectedTrash.push(newData.data[i])
    }
  }

  removeFromList(trashId: string) {
    const index = this.selectedTrash.findIndex(t => t.id === trashId)
    this.allMarkers.push(this.selectedTrash[index])
    this.selectedTrash.splice(index, 1)

    //rerender table
    const newData = new MatTableDataSource<MarkerModel>(this.selectedTrash);
    this.selectedTrash = []
    for (let i = 0; i < newData.data.length; i++) {
      this.selectedTrash.push(newData.data[i])
    }

  }

}
