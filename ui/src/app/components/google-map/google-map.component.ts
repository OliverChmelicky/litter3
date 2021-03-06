import {Component, OnInit, ViewChild} from '@angular/core';
import {LocationService} from "../../services/location/location.service";
import {MapLocationModel} from "../../models/GPSlocation.model";
import {GoogleMap, LatLng} from "@agm/core/services/google-maps-types";
import {MarkerModel} from "src/app/components/google-map/Marker.model";
import {MouseEvent} from '@agm/core';
import {Router} from "@angular/router";
import {TrashService} from "../../services/trash/trash.service";
import {AgmMap} from '@agm/core';
import {MatCheckboxChange} from "@angular/material/checkbox";
import {initialDistance} from "../../models/shared.models";

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

  location: MapLocationModel;
  defaultLocation = czechPosition;
  map: GoogleMap;
  allMarkers: MarkerModel[];
  filteredMarkers: MarkerModel[];
  selectedMarker: MarkerModel;

  showCleaned: boolean = true;
  showNotCleaned: boolean = true;

  borderTop: number;
  borderBottom: number;
  borderLeft: number;
  borderRight: number;

  initialDistance:number = initialDistance

  constructor(
    private readonly locationService: LocationService,
    private trashService: TrashService,
    private router: Router,
  ) {
  }

  ngOnInit() {
    //https://github.com/SebastianM/angular-google-maps/issues/1238 => grey map
    this.location = this.defaultLocation;
    this.locationService.getPosition().then(data => {
      this.location = data;
    }).catch(() => {});
    this.allMarkers = [];
  }

  sleep(ms): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  async onMapReady(map: GoogleMap) {
    this.map = map;
    //In an issue it was written that this helps but don`t
    await setTimeout(() => {
      this.agmMap.triggerResize();
    }, 1000)
    this.initialMapFetching();
  }

  addMarker(lat: number, lng: number) {
    this.selectedMarker = {
      lat: lat,
      lng: lng,
      new: true,
      cleaned: false,
      id: '',
    };
    this.allMarkers.push(this.selectedMarker)
    this.applyMarkerFilters()
  }

  createTrash(marker: MarkerModel) {
    this.router.navigate(['report', marker.lat, marker.lng])
  }

  selectMarker(i: number, event) {
    this.selectedMarker = this.allMarkers[i]
  }

  dragging(i: number, $event: MouseEvent) {
    this.allMarkers[i].lat = $event.coords.lat;
    this.allMarkers[i].lng = $event.coords.lng;
    this.selectedMarker.lat = $event.coords.lat;
    this.selectedMarker.lng = $event.coords.lng;
  }

  onBoundsChange() {
    let c = this.map.getCenter()

    let visibleTop = c.lat() + 3.4
    let visibleBottom = c.lat() - 3.4
    let visibleRight = c.lng() + 8.82
    let visibleLeft = c.lng() - 8.82

    if (!this.map.getBounds()) {
      this.initialMapFetching()
    } else {
      const p1 = this.map.getBounds().getNorthEast()
      const p2 = this.map.getBounds().getSouthWest()

      visibleTop = p1.lat()
      visibleRight = p1.lng()
      visibleBottom = p2.lat()
      visibleLeft = p2.lng()
    }


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
        this.allMarkers = this.getOnlyNewMarkers();
        for (let i = 0; i < trash.length; i++) {
          this.allMarkers.push({
            lat: trash[i].Location[0],
            lng: trash[i].Location[1],
            new: false,
            id: trash[i].Id,
            cleaned: trash[i].Cleaned,
            images: trash[i].Images ? trash[i].Images : [],
            numOfCollections: trash[i].Collections ? trash[i].Collections.length : 0
          })

          this.applyMarkerFilters()
        }

        const viewCenter = this.map.getCenter()
        let r = 2 * Math.abs(p1.lat() - viewCenter.lat())

        if (p1.lat() < 0) {
          this.borderTop = p1.lat() + r
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

  initialMapFetching() {
    let c = this.map.getCenter()
    this.borderTop = c.lat() + 3.4
    this.borderBottom = c.lat() - 3.4

    this.borderRight = c.lng() + 8.82
    this.borderLeft = c.lng() - 8.82

    this.trashService.getTrashInRange(this.map.getCenter().lat(), this.map.getCenter().lng(), this.initialDistance).subscribe(
      trash => {
        this.allMarkers = this.getOnlyNewMarkers();
        for (let i = 0; i < trash.length; i++) {
          this.allMarkers.push({
            lat: trash[i].Location[0],
            lng: trash[i].Location[1],
            new: false,
            id: trash[i].Id,
            cleaned: trash[i].Cleaned,
            images: trash[i].Images ? trash[i].Images : [],
            numOfCollections: trash[i].Collections ? trash[i].Collections.length : 0
          })

          this.applyMarkerFilters()
        }
      })
  }

  navigateToTrash(id: string) {
    this.router.navigate(['trash/details', id])
  }

  onRightClick() {
    //https://github.com/SebastianM/angular-google-maps/issues/797
    // console.log(event)
    // console.log(event.infoWindow.close())
  }

  onCleanedOption(event: MatCheckboxChange) {
    this.showCleaned = event.checked;
    this.applyMarkerFilters()
  }

  onNotCleanedOption(event: MatCheckboxChange) {
    this.showNotCleaned = event.checked;
    this.applyMarkerFilters()
  }

  applyMarkerFilters() {
    let futureVisibleMarkers = this.allMarkers;

    if (!this.showCleaned) {
      futureVisibleMarkers = this.filterCleaned(futureVisibleMarkers)
    }
    if (!this.showNotCleaned) {
      futureVisibleMarkers = this.filterNotCleaned(futureVisibleMarkers)
    }

    this.filteredMarkers = futureVisibleMarkers
  }


  //I want not cleaned
  filterCleaned(markers: MarkerModel[]): MarkerModel[] {
    return markers.filter(marker => {
      if (marker.cleaned === false || marker.new === true) {
        return marker
      }
    })
  }

  //I want cleaned
  filterNotCleaned(markers: MarkerModel[]): MarkerModel[] {
    return markers.filter(marker => {
      if (marker.cleaned === true || marker.new === true) {
        return marker
      }
    })
  }

  private getOnlyNewMarkers() {
    return this.allMarkers.filter(marker => {
      if (marker.new === true) {
        return marker
      }
    })
  }
}
