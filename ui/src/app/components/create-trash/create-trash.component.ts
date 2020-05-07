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
    trashTypes: [''],
    accessibility: [''],
    description: [''],
    anonymously: [''],
  });

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

  map: GoogleMap
  initMapLat: number;
  initMapLng: number;
  markerLat: number;
  markerLng: number;


  constructor(
    private route: ActivatedRoute,
    private trashService: TrashService,
    private formBuilder: FormBuilder,
    private readonly locationService: LocationService,
    private http: HttpClient,
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
    this.trashForm.value["lat"] = this.markerLat
    this.trashForm.value["lng"] = this.markerLng
    console.log(this.trashForm.value)
  }


  onFileSelected(event) {
    console.log(event)

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
        if (event.target.files.length !== 0) {
          const fd = new FormData();
          for(let i =0; i < event.target.files.length; i++){
            fd.append("files", event.target.files[i], event.target.files[i].name);
          }
          //fd.append('file', event.target.files[0], event.target.files[0].name);
          const url = `${ApisModel.apiUrl}/${ApisModel.fileupload}/${ApisModel.trash}/${trash.Id}`;
          return this.http.post(url, fd).subscribe(
            res => {
              console.log('Huraaa uploadlo')
              console.log(res)
            },
            error => console.log(error)
          )
        } else {
          console.log('Less then 0 pictures')
        }
      })

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
