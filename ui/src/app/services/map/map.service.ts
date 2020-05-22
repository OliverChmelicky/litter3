import {Injectable} from '@angular/core';
import {TrashService} from "../trash/trash.service";
import {GoogleMap} from "@agm/core/services/google-maps-types";
import {MarkerModel} from "../../components/google-map/Marker.model";
import {MapLoadAnswer} from "../../models/shared.models";
import {BehaviorSubject, Observable} from "rxjs";
import {defaultTrashImage} from "../../models/trash.model";

@Injectable({
  providedIn: 'root'
})
export class MapService {

  //TODO use this service for loading new markers in map and createEvent component

  allMarkers: MarkerModel[];
  exampleBinUrl: string = ''
  borderTop: number;
  borderBottom: number;
  borderLeft: number;
  borderRight: number;
  private result = new BehaviorSubject<MapLoadAnswer>({
    LoadedMarkers: [],
    BorderBottom: 0,
    BorderLeft: 0,
    BorderRight: 0,
    BorderTop: 0,
  });

  constructor(
    private trashService: TrashService,
  ) {
  }


  loadNewMarkers(map: GoogleMap) {
    const p1 = map.getBounds().getNorthEast()
    const p2 = map.getBounds().getSouthWest()

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
    this.trashService.getTrashInRange(map.getCenter().lat(), map.getCenter().lng(), d * 2).subscribe(
      trash => {
        this.allMarkers = [];
        for (let i = 0; i < trash.length; i++) {
          this.allMarkers.push({
            lat: trash[i].Location[0],
            lng: trash[i].Location[1],
            new: false,
            id: trash[i].Id,
            cleaned: trash[i].Cleaned,
            images: trash[i].Images ? trash[i].Images : [defaultTrashImage],
            numOfCollections: trash[i].Collections ? trash[i].Collections.length : 0
          })
        }

        const viewCenter = map.getCenter()
        let r = 2 * Math.abs(p1.lat() - viewCenter.lat())
        console.log('R: ', r)

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

      },
      err => console.log(err),
      () => this.result.next({
        LoadedMarkers: this.allMarkers,
        BorderTop: this.borderTop,
        BorderBottom: this.borderBottom,
        BorderLeft: this.borderLeft,
        BorderRight: this.borderRight,
      })
    )
  }


}
