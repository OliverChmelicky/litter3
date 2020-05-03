import {Component, OnInit} from '@angular/core';
import {LocationService} from "../../services/location/location.service";
import {MapLocationModel} from "../../models/GPSlocation.model";
import {google, GoogleMap} from "@agm/core/services/google-maps-types";
import {MarkerModel} from "src/app/components/google-map/Marker.model";
import {NgZone} from '@angular/core';
import {MouseEvent} from '@agm/core';

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
    private ngZone: NgZone,
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
    this.markers.push({
      lat: lat,
      lng: lng,
      new: true,
    });
  }

  createTrash(marker: MarkerModel) {
    window.alert(marker.lng + ' ' + marker.lat)
  }

  selectMarker(event) {
    this.selectedMarker = {
      lat: event.latitude,
      lng: event.longitude
    }
  }

  markerDragEnd(m: MarkerModel, i: number, $event: MouseEvent) {
    this.markers[i].lat = $event.coords.lat;
    this.markers[i].lng = $event.coords.lng;
    window.alert(this.markers[i].lat + ' ' + this.markers[i].lng)
  }
}
