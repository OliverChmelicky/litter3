import { Injectable } from '@angular/core';
import {GoogleMap} from "@agm/core/services/google-maps-types";
import {TrashService} from "../trash/trash.service";
import {MarkerModel} from "../../components/google-map/Marker.model";
import {MarkersAfterInitModel} from "../../models/shared.models";

@Injectable({
  providedIn: 'root'
})
export class LocationService {

  initialDistance: number = 30000;
  allMarkers: MarkerModel[] = [];

  constructor(
    private trashService: TrashService,
  ) { }

  getPosition(): Promise<any>
  {
    return new Promise((resolve, reject) => {
      navigator.geolocation.getCurrentPosition(resp => {
          resolve({lng: resp.coords.longitude, lat: resp.coords.latitude});
        },
        err => {
          reject(err);
        },
        {timeout: 10000}
        );
    });

  }

  //I can just return a promise
  // getMarkersAfterMapInit(map: GoogleMap): MarkersAfterInitModel {
  //   this.allMarkers = [];
  //
  //   let c = map.getCenter()
  //   const borderTop = c.lat() + 3.4
  //   const borderBottom = c.lat() - 3.4
  //
  //   const borderRight = c.lng() + 8.82
  //   const borderLeft = c.lng() - 8.82
  //
  //   this.trashService.getTrashInRange(c.lat(), c.lng(), this.initialDistance).subscribe(
  //     trash => {
  //       for (let i = 0; i < trash.length; i++) {
  //         this.allMarkers.push({
  //           lat: trash[i].Location[0],
  //           lng: trash[i].Location[1],
  //           new: false,
  //           id: trash[i].Id,
  //           cleaned: trash[i].Cleaned,
  //           images: trash[i].Images ? trash[i].Images : [],
  //           numOfCollections: trash[i].Collections ? trash[i].Collections.length : 0
  //         })
  //       }
  //     },
  //     () => {},
  //     () => {
  //       return {
  //         markers: this.allMarkers,
  //         borderTop: borderTop,
  //         borderBottom: borderBottom,
  //         borderLeft: borderLeft,
  //         borderRight: borderRight,
  //       }
  //     }
  //     )
  // }
}
