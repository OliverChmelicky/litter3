import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from "@angular/router";
import {TrashService} from "../../services/trash/trash.service";
import {GoogleMap} from "@agm/core/services/google-maps-types";
import {MouseEvent} from '@agm/core';
import {FormBuilder} from "@angular/forms";
import {LocationService} from "../../services/location/location.service";
import {TrashModel,} from "../../models/trash.model";
import {FileuploadService} from "../../services/fileupload/fileupload.service";
import {accessibilityChoces} from "../../models/accessibilityChocies";
import {AuthService} from "../../services/auth/auth.service";

@Component({
  selector: 'app-create-trash',
  templateUrl: './create-trash.component.html',
  styleUrls: ['./create-trash.component.css']
})
export class CreateTrashComponent implements OnInit {
  trash: TrashModel;
  trashForm = this.formBuilder.group({
    lat: '',
    lng: '',
    size: 0,

    trashTypeHousehold: false,
    trashTypeAutomotive: false,
    trashTypeConstruction: false,
    trashTypePlastics: false,
    trashTypeElectronic: false,
    trashTypeGlass: false,
    trashTypeMetal: false,
    trashTypeDangerous: false,
    trashTypeCarcass: false,
    trashTypeOrganic: false,
    trashTypeOther: false,

    accessibility: [''],
    description: '',
    anonymously: false,
  });

  isLoggedIn: boolean = false

  accessibilityChoices = accessibilityChoces;

  map: GoogleMap
  initMapLat: number;
  initMapLng: number;
  markerLat: number;
  markerLng: number;
  fd: FormData = new FormData();


  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private trashService: TrashService,
    private fileuploadService: FileuploadService,
    private formBuilder: FormBuilder,
    private readonly locationService: LocationService,
    private authService: AuthService,
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
    this.authService.isLoggedIn.subscribe( res => this.isLoggedIn = res)
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
      Accessibility: this.trashForm.value['accessibility'],
      TrashType: this.changeTrashTypeToInt(),
      Location: [this.markerLat, this.markerLng],
      Description: this.trashForm.value['description'],
      FinderId: '',
      Anonymously: this.trashForm.value['anonymously'],
    }

    this.trashService.createTrash(this.trash).subscribe(
      trash => {
        if (this.fd.getAll('files').length !== 0) {
          this.fileuploadService.uploadTrashImages(this.fd, trash.Id).subscribe(
            () => {
              this.fd.delete('files')
              this.router.navigateByUrl('map')
            })
        } else {
          this.router.navigateByUrl('map')
        }
      })
  }


  onFileSelected(event) {
    this.fd.delete('files')
    for (let i = 0; i < event.target.files.length; i++) {
      this.fd.append("files", event.target.files[i], event.target.files[i].name);
    }
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

  private changeTrashTypeToInt(): number {
    return this.trashService.convertTrashTypeBoolsToNums(
      {
        TrashTypeHousehold: !!this.trashForm.value.trashTypeHousehold,
        TrashTypeAutomotive: !!this.trashForm.value.trashTypeAutomotive,
        TrashTypeConstruction: !!this.trashForm.value.trashTypeConstruction,
        TrashTypePlastics: !!this.trashForm.value.trashTypePlastics,
        TrashTypeElectronic: !!this.trashForm.value.trashTypeElectronic,
        TrashTypeGlass: !!this.trashForm.value.trashTypeGlass,
        TrashTypeMetal: !!this.trashForm.value.trashTypeMetal,
        TrashTypeDangerous: !!this.trashForm.value.trashTypeDangerous,
        TrashTypeCarcass: !!this.trashForm.value.trashTypeCarcass,
        TrashTypeOrganic: !!this.trashForm.value.trashTypeOrganic,
        TrashTypeOther: !!this.trashForm.value.trashTypeOther,
      }
    );
  }

}
