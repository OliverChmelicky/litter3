import {Component, OnInit, ViewChild} from '@angular/core';
import {LocationService} from "../../services/location/location.service";
import {MapLocationModel} from "../../models/GPSlocation.model";
import {GoogleMap, LatLng, LatLngBounds, LatLngLiteral} from "@agm/core/services/google-maps-types";
import {MarkerModel} from "src/app/components/google-map/Marker.model";
import {MouseEvent} from '@agm/core';
import {Router} from "@angular/router";
import {TrashService} from "../../services/trash/trash.service";
import { AgmMap } from '@agm/core';

export const czechPosition: MapLocationModel = {
  lat: 49.81500022397678,
  lng: 20.0,
  zoom: 7,
  minZoom: 3,
};

@Component({
  selector: 'app-google-map',
  templateUrl: './google-map.component.html',
  styleUrls: ['./google-map.component.css']
})
export class GoogleMapComponent implements OnInit {
  @ViewChild('mymap') agmMap : AgmMap

  funcActivatedByParentComponentOnAccordionOpen() {
    this.agmMap.triggerResize();
  }

  location: MapLocationModel;
  defaultLocation = czechPosition;
  map: GoogleMap;
  markers: MarkerModel[];
  selectedMarker: MarkerModel;
  currentLocation: LatLng;

  visibleTop: number;
  visibleBottom: number;
  visibleLeft: number;
  visibleRight: number;

  constructor(
    private readonly locationService: LocationService,
    private trashService: TrashService,
    private router: Router,
  ) {
  }

  ngOnInit() {
    this.location = this.defaultLocation;
    this.locationService.getPosition().then(data => {
      this.location = data;
    }).catch(err => window.alert('Error getting location: ' + err));
    this.markers = [];
  }

  async onMapReady(map: GoogleMap) {
    this.map = map;
    setTimeout(()=>{ this.agmMap.triggerResize(); },1000)
    this.loadNewMarkers()
  }

  addMarker(lat: number, lng: number) {
    this.selectedMarker = {
      lat: lat,
      lng: lng,
      new: true,
    };
    this.markers.push(this.selectedMarker)
  }

  createTrash(marker: MarkerModel) {
    this.router.navigate(['report', marker.lat, marker.lng])
  }

  selectMarker(i: number, event) {
    console.log(event)
    console.log(event.infoWindow.close())
    this.selectedMarker = {
      lat: event.latitude,
      lng: event.longitude,
      new: this.markers[i].new
    }
  }

  dragging(i: number, $event: MouseEvent) {
    this.markers[i].lat = $event.coords.lat;
    this.markers[i].lng = $event.coords.lng;
    this.selectedMarker.lat = $event.coords.lat;
    this.selectedMarker.lng = $event.coords.lng;
  }

  onBoundsChange() {
    const currentCenter = this.map.getCenter()

    if (currentCenter.lng() > this.visibleRight || currentCenter.lng() < this.visibleLeft) {
      console.log('reached visible sides')
      this.loadNewMarkers()
    } else if (currentCenter.lat() < this.visibleBottom || currentCenter.lat() > this.visibleTop) {
      console.log('reached visible tops and bottoms')
      this.loadNewMarkers()
    }

  }

  loadNewMarkers() {
    const p1 = this.map.getBounds().getNorthEast()
    const p2 = this.map.getBounds().getSouthWest()

    this.visibleTop = p1.lat()
    this.visibleRight = p1.lng()

    this.visibleBottom = p2.lat()
    this.visibleLeft = p2.lng()

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
        this.markers = [];
        for (let i = 0; i < trash.length; i++) {
          this.markers.push({
            lat: trash[i].Location[0],
            lng: trash[i].Location[1],
            new: false,
          })
        }
      }
    )
  }

  hideMarkerWindows() {
  }
}
