import {Component, OnInit} from '@angular/core';
import {LocationService} from "../services/location/location.service";
import {GPSlocationModel} from "../models/GPSlocation.model";
import {GoogleMap} from "@agm/core/services/google-maps-types";
import {MarkerModel} from "src/app/google-map/Marker.model";

export const czechPosition: GPSlocationModel = {
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
  location: GPSlocationModel;
  defaultLocation = czechPosition;
  map: GoogleMap;
  markers: MarkerModel[];

  constructor(private readonly locationService: LocationService) {
  }

  ngOnInit() {
    this.location = this.defaultLocation;
    this.locationService.getPosition().then(data => {
      this.location = data;
    });
    this.markers = [];
  }

  onMapReady(map: GoogleMap) {
    this.map = map;
  }

  addMarker(lat: number, lng: number) {
    this.markers.push({
      lat,
      lng
    });
  }

}
