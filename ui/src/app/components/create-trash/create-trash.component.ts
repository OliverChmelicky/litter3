import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from "@angular/router";
import {TrashService} from "../../services/trash/trash.service";
import {GoogleMap} from "@agm/core/services/google-maps-types";
import {MouseEvent} from '@agm/core';
import {FormBuilder} from "@angular/forms";
import {LocationService} from "../../services/location/location.service";
import {TrashModel} from "../../models/trash.model";
import {HttpClient} from "@angular/common/http";
import {ApisModel} from "../../api/api-urls";
import {FileuploadService} from "../../services/fileupload/fileupload.service";
import {accessibilityChoces} from "./accessibilityChocies";

@Component({
  selector: 'app-create-trash',
  templateUrl: './create-trash.component.html',
  styleUrls: ['./create-trash.component.css']
})
export class CreateTrashComponent implements OnInit {
  trash: TrashModel;
  trashForm = this.formBuilder.group({
    lat: [''],
    lng: [''],
    size: [1],

    trashTypeHousehold: [''],
    trashTypeAutomotive: [''],
    trashTypeConstruction: [''],
    trashTypePlastics: [''],
    trashTypeElectronic: [''],
    trashTypeGlass: [''],
    trashTypeMetal: [''],
    trashTypeDangerous: [''],
    trashTypeCarcass: [''],
    trashTypeOrganic: [''],
    trashTypeOther: [''],

    accessibility: [''],
    description: [''],
    anonymously: [''],
  });

  accessibilityChoices = accessibilityChoces;

  map: GoogleMap
  initMapLat: number;
  initMapLng: number;
  markerLat: number;
  markerLng: number;
  fd: FormData = new FormData();


  constructor(
    private route: ActivatedRoute,
    private trashService: TrashService,
    private fileuploadService: FileuploadService,
    private formBuilder: FormBuilder,
    private readonly locationService: LocationService,
  ) {
  }

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      this.initMapLat = +params.get('lat');
      this.initMapLng = +params.get('lng');
      this.markerLat = this.initMapLat;
      this.markerLng = this.initMapLng;

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

  onSubmit() {
    this.trash = {
      Id: '',
      Cleaned: false,
      Size: this.printSize(),
      Accessibility: this.printAccessibility(),
      TrashType: this.changeTrashTypeToInt(),
      Location: [this.markerLat, this.markerLng],
      Description: this.trashForm.value["description"],
      FinderId: '',
    }

    this.trashService.createTrash(this.trash).subscribe(
      trash => {
        if (this.fd.getAll('files').length !== 0) {
          this.fileuploadService.uploadTrashImages(this.fd, trash.Id).subscribe(
            () => {
              this.fd.delete('files')
              //this.trashForm.reset()
            })
        } else {
          console.log('Less then 0 pictures')
          //this.trashForm.reset()
        }
      })
  }


  onFileSelected(event) {
    this.fd.delete('files')
    for (let i = 0; i < event.target.files.length; i++) {
      this.fd.append("files", event.target.files[i], event.target.files[i].name);
    }
    console.log(this.fd.getAll('files').length)
  }

  printSize() {
    if (this.trashForm.value["size"] == 0) {
      return 'unknown';
    }
    if (this.trashForm.value["size"] == 1) {
      return 'bag';
    }
    if (this.trashForm.value["size"] == 2) {
      return 'wheelbarrow';
    }
    if (this.trashForm.value["size"] == 3) {
      return 'car';
    }
  }

  private printAccessibility() {
    if (this.trashForm.value["size"] == 0) {
      return 'unknown';
    }
    if (this.trashForm.value["size"] == 1) {
      return 'easy';
    }
    if (this.trashForm.value["size"] == 2) {
      return 'car';
    }
    if (this.trashForm.value["size"] == 3) {
      return 'cave';
    }
    if (this.trashForm.value["size"] == 4) {
      return 'underWater';
    }
  }

  private changeTrashTypeToInt() {
    return 0;
  }

}
