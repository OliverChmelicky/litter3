import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from "@angular/router";
import {TrashService} from "../../services/trash/trash.service";
import {GoogleMap} from "@agm/core/services/google-maps-types";
import {MouseEvent} from '@agm/core';
import {FormArray, FormBuilder, FormControl} from "@angular/forms";
import {LocationService} from "../../services/location/location.service";

@Component({
  selector: 'app-create-trash',
  templateUrl: './create-trash.component.html',
  styleUrls: ['./create-trash.component.css']
})
export class CreateTrashComponent implements OnInit {
  trashForm = this.formBuilder.group({
    lat: [''],
    lng: [''],
    size: [1],
    trashTypes: [''],
    description: [''],
    anonymously: [''],
  });
  map: GoogleMap

  trashTypesChoices: string[] = [
    '---other---',
    'household',
    'automotive',
    'construction',
    'plastics',
    'electronic',
    'glass',
    'metal',
    'dangerous',
    'carcass',
    'organic',
  ]

  initMapLat: number;
  initMapLng: number;
  markerLat: number;
  markerLng: number;


  constructor(
    private route: ActivatedRoute,
    private trashService: TrashService,
    private formBuilder: FormBuilder,
    private readonly locationService: LocationService,
  ) {
  }

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      this.initMapLat = +params.get('lat');
      this.initMapLng = +params.get('lng');

      if (this.initMapLat === 0 && this.initMapLng === 0) {
        this.locationService.getPosition().then(data => {
          this.initMapLat = data.lat;
          this.initMapLng = data.lng;
          this.markerLat = this.initMapLat;
          this.markerLng = this.initMapLng;
        }).catch(
          () => {
            this.initMapLat = 49;
            this.initMapLng = 16;
            this.markerLat = this.initMapLat;
            this.markerLng = this.initMapLng;
          }
        );
      }
    });
  }

  onMapReady(map: GoogleMap) {
    this.map = map;
  }

  onDragging($event: MouseEvent) {
    this.markerLat = $event.coords.lat;
    this.markerLng = $event.coords.lng;
  }

  printSize() {
    if (this.trashForm.value["size"] == 0) {
      return 'undefined';
    }
    if (this.trashForm.value["size"] == 1) {
      return 'small';
    }
    if (this.trashForm.value["size"] == 2) {
      return 'medium';
    }
    if (this.trashForm.value["size"] == 3) {
      return 'big';
    }
  }


  onSubmit() {
    this.trashForm.value["lat"] = this.markerLat
    this.trashForm.value["lng"] = this.markerLng
    console.log(this.trashForm.value)
  }


}
