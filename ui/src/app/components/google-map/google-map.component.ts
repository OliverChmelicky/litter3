import {Component, OnInit, ViewChild} from '@angular/core';
import {LocationService} from "../../services/location/location.service";
import {MapLocationModel} from "../../models/GPSlocation.model";
import {GoogleMap, LatLng, LatLngBounds, LatLngLiteral} from "@agm/core/services/google-maps-types";
import {MarkerModel} from "src/app/components/google-map/Marker.model";
import {MouseEvent} from '@agm/core';
import {Router} from "@angular/router";
import {TrashService} from "../../services/trash/trash.service";
import { AgmMap } from '@agm/core';
import {MatCheckboxChange} from "@angular/material/checkbox";

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
  @ViewChild('agmMap') agmMap: AgmMap;
  exampleBinUrl: 'https://www.google.com/url?sa=i&url=https%3A%2F%2Fvisualpharm.com%2Ffree-icons%2Fblank%2Fblank%2520trash&psig=AOvVaw07-6SZ8RD7AhPn2ddRQm6W&ust=1589094336519000&source=images&cd=vfe&ved=0CAIQjRxqFwoTCKjtme2bpukCFQAAAAAdAAAAABAF';

  location: MapLocationModel;
  defaultLocation = czechPosition;
  map: GoogleMap;
  allMarkers: MarkerModel[];
  filteredMarkers: MarkerModel[];
  selectedMarker: MarkerModel;

  hideCleaned: boolean = true;
  hideNotCleaned: boolean = true;

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
    this.allMarkers = [];
  }

  sleep(ms): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  async onMapReady(map: GoogleMap) {
    this.map = map;
    //In an issue it was written that this helps but don`t
    // await this.sleep(2000)
    // await setTimeout(()=>{ this.agmMap.triggerResize(); },500)
    // await this.sleep(2000)
    this.loadNewMarkers()
  }

  addMarker(lat: number, lng: number) {
    this.selectedMarker = {
      lat: lat,
      lng: lng,
      new: true,
      id: '',
    };
    this.allMarkers.push(this.selectedMarker)
  }

  createTrash(marker: MarkerModel) {
    this.router.navigate(['report', marker.lat, marker.lng])
  }

  selectMarker(i: number, event) {
    // console.log(event)
    // console.log(event.infoWindow.close())
    this.selectedMarker = this.allMarkers[i]
  }

  dragging(i: number, $event: MouseEvent) {
    this.allMarkers[i].lat = $event.coords.lat;
    this.allMarkers[i].lng = $event.coords.lng;
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
        this.allMarkers = [];
        for (let i = 0; i < trash.length; i++) {
          this.allMarkers.push({
            lat: trash[i].Location[0],
            lng: trash[i].Location[1],
            new: false,
            id: trash[i].Id,
            images: trash[i].Images ? trash[i].Images : [this.exampleBinUrl],
            numOfCollections: trash[i].Collections ? trash[i].Collections.length : 0
          })
        }

        let futureVisibleMarkers = this.allMarkers

        if (this.hideNotCleaned) {
          futureVisibleMarkers = this.filterNotCleaned(futureVisibleMarkers)
        }
        if (this.hideCleaned) {
          futureVisibleMarkers = this.filterCleaned(futureVisibleMarkers)
        }
        this.filteredMarkers = futureVisibleMarkers

      }
    )
  }

  navigateToTrash(id: string) {
    this.router.navigate(['trash/details', id])
  }

  onRightClick() {
    //https://github.com/SebastianM/angular-google-maps/issues/797
  }

  onCleanedOption(event: MatCheckboxChange) {
    this.hideCleaned = event.checked;
    this.processFilterChange()
  }

  onNotCleanedOption(event: MatCheckboxChange) {
    this.hideNotCleaned = event.checked;
    this.processFilterChange()
  }

  processFilterChange() {
    let futureVisibleMarkers = this.allMarkers;

    if (this.hideCleaned) {
      this.filterCleaned(futureVisibleMarkers)
    }
    if (this.hideNotCleaned) {
      this.filterNotCleaned(futureVisibleMarkers)
    }

    this.filteredMarkers = futureVisibleMarkers
  }


  filterCleaned(markers: MarkerModel[]): MarkerModel[]{
    return markers.filter( marker => marker.cleaned === true )
  }

  filterNotCleaned(markers: MarkerModel[]): MarkerModel[]{
    return markers.filter( marker => marker.cleaned === false )
  }

}
