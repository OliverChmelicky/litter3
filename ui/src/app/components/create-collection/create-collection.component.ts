import { Component, OnInit } from '@angular/core';
import {GoogleMap} from "@agm/core/services/google-maps-types";
import {ActivatedRoute} from "@angular/router";
import {TrashService} from "../../services/trash/trash.service";
import {TrashModel, MarkerCollectionModel, defaultTrashImage} from "../../models/trash.model";
import {MarkerModel} from "../google-map/Marker.model";
import {czechPosition} from "../event-details/event-details.component";
import {MatTableDataSource} from "@angular/material/table";
import {AuthService} from "../../services/auth/auth.service";

export const createCollectionTrashColumns: string[] = [
  'trash-image',
  'weight',
  'cleaned-trash',
  'images-btn',
]


@Component({
  selector: 'app-create-collection',
  templateUrl: './create-collection.component.html',
  styleUrls: ['./create-collection.component.css']
})
export class CreateCollectionComponent implements OnInit {
  trashIds: string[] = [];
  trash: TrashModel[] = [];
  eventId: string = '';

  map: GoogleMap;
  allMarkers: MarkerCollectionModel[] = [];
  initLat: number = czechPosition.lat;
  initLng: number = czechPosition.lng;

  selectedMarkers: MarkerCollectionModel[] = [];
  tableColumns = createCollectionTrashColumns;


  constructor(
    private route: ActivatedRoute,
    private trashService: TrashService,
    ) {
  }

  ngOnInit(): void {
    this.route.queryParamMap.subscribe(params => {
      this.eventId = params.get('eventId')
      this.trashIds = params.getAll('trashIds')
      console.log('param event: ', this.eventId)
      console.log('param trash: ', this.trashIds)
      this.trashService.getTrashByIds(this.trashIds).subscribe( trash => {
        this.trash = trash
        this.initLat = trash[0].Location[0]
        this.initLng = trash[0].Location[1]
        this.assignMarkers()
      })
    });
  }

  onMapReady(map: GoogleMap) {
    this.map = map
  }

  private assignMarkers() {
    this.trash.map(t => {
      let collLength = 0
      if (t.Collections) {
        collLength = t.Collections.length
      }
      if (!t.Images) {
        t.Images = [];
      }
      this.allMarkers.push({
        trashId: t.Id,
        lat: t.Location[0],
        lng: t.Location[1],
        cleaned: t.Cleaned,
        image: t.Images ? t.Images[0] : defaultTrashImage,
        numOfCollections: collLength,
        collectionWeight: 0,
        collectionCleanedTrash: false,
        collectionEventId: this.eventId,
        collectionImages: [],
        isInList: false,
      })
    })
  }

  addToList(marker: MarkerCollectionModel) {
    marker.isInList = true
    this.selectedMarkers.push(marker)
  }

  removeFromList(trashId: string) {
    const index = this.selectedMarkers.findIndex(t => t.trashId === trashId)
    this.selectedMarkers.splice(index, 1)

    let marker = this.selectedMarkers[index]
    marker.isInList = false
    this.allMarkers.push(marker)

    //rerender table
    const newData = new MatTableDataSource<MarkerCollectionModel>(this.selectedMarkers);
    this.selectedMarkers = []
    for (let i = 0; i < newData.data.length; i++) {
      this.selectedMarkers.push(newData.data[i])
    }

  }
}
