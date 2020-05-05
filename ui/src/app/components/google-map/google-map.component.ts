import {Component, OnInit} from '@angular/core';
import {LocationService} from "../../services/location/location.service";
import {MapLocationModel} from "../../models/GPSlocation.model";
import {GoogleMap} from "@agm/core/services/google-maps-types";
import {MarkerModel} from "src/app/components/google-map/Marker.model";
import {MouseEvent} from '@agm/core';
import {Router} from "@angular/router";

export const czechPosition: MapLocationModel = {
  lat: 49.81500022397678,
  lng: 20.0,
  zoom: 7,
};

@Component({
  selector: 'app-google-map',
  templateUrl: './google-map.component.html',
  styleUrls: ['./google-map.component.css']
})
export class GoogleMapComponent implements OnInit {
  location: MapLocationModel;
  defaultLocation = czechPosition;
  map: GoogleMap;
  markers: MarkerModel[];
  selectedMarker: MarkerModel;

  constructor(
    private readonly locationService: LocationService,
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

  onMapReady(map: GoogleMap) {
    this.map = map;
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
}
