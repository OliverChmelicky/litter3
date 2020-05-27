import {Component, OnInit} from '@angular/core';
import {GoogleMap} from "@agm/core/services/google-maps-types";
import {ActivatedRoute, Router} from "@angular/router";
import {TrashService} from "../../services/trash/trash.service";
import {
  CollectionModel,
  CreateCollectionModel,
  defaultTrashImage,
  MarkerCollectionModel,
  TrashModel
} from "../../models/trash.model";
import {czechPosition} from "../event-details/event-details.component";
import {MatTableDataSource} from "@angular/material/table";
import {EventService} from "../../services/event/event.service";
import {EventPickerModel} from "../../models/event.model";
import {FileuploadService} from "../../services/fileupload/fileupload.service";

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
  notSelectedMarkers: MarkerCollectionModel[] = [];
  initLat: number = czechPosition.lat;
  initLng: number = czechPosition.lng;

  selectedMarkers: MarkerCollectionModel[] = [];
  tableColumns = createCollectionTrashColumns;
  organizerId: EventPickerModel;
  errorMessage: string;


  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private trashService: TrashService,
    private eventService: EventService,
    private fileuploadService: FileuploadService,
  ) {
  }

  ngOnInit(): void {
    this.route.queryParamMap.subscribe(params => {
      this.eventId = params.get('eventId')
      this.trashIds = params.getAll('trashIds')
      this.trashService.getTrashByIds(this.trashIds).subscribe(trash => {
        this.trash = trash
        this.initLat = trash[0].Location[0]
        this.initLng = trash[0].Location[1]
        this.assignMarkers()
      })
    });
    this.organizerId = this.eventService.getEventEditor()
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
        t.Images = [defaultTrashImage];
      }
      this.notSelectedMarkers.push({
        trashId: t.Id,
        lat: t.Location[0],
        lng: t.Location[1],
        cleaned: t.Cleaned,
        image: t.Images[0],
        numOfCollections: collLength,
        collectionWeight: 0,
        collectionCleanedTrash: false,
        collectionEventId: this.eventId,
        collectionImages: new FormData(),
        isInList: false,
      })
    })
  }

  addToList(marker: MarkerCollectionModel) {
    marker.isInList = true
    this.selectedMarkers.push(marker)

    const index = this.notSelectedMarkers.findIndex(t => t.trashId === marker.trashId)
    this.notSelectedMarkers.splice(index, 1)
  }

  removeFromList(trashId: string) {
    const index = this.selectedMarkers.findIndex(t => t.trashId === trashId)
    this.selectedMarkers.splice(index, 1)

    let marker = this.selectedMarkers[index]
    marker.isInList = false
    this.notSelectedMarkers.push(marker)

    //rerender table
    const newData = new MatTableDataSource<MarkerCollectionModel>(this.selectedMarkers);
    this.selectedMarkers = []
    for (let i = 0; i < newData.data.length; i++) {
      this.selectedMarkers.push(newData.data[i])
    }

  }

  onFileSelected(event, i: number) {
    this.selectedMarkers[i].collectionImages.delete('files')
    for (let i = 0; i < event.target.files.length; i++) {
      this.selectedMarkers[i].collectionImages.append("files", event.target.files[i], event.target.files[i].name);
    }
  }

  onCreate() {
    const colections = this.mapFromColMarkersToColections(this.selectedMarkers)
    let invalidInput = false
    this.selectedMarkers.map( m => {
      if (m.collectionWeight <= 0) {
        this.errorMessage = 'You must specify a weght greater than 0!'
        invalidInput = true
      }
    })

    if (invalidInput) {
      return
    }

    let collectionsToCreate: CreateCollectionModel = {
      Collections: colections,
      AsSociety: this.organizerId.AsSociety,
      OrganizerId: this.organizerId.Id,
      EventId: this.eventId,
    }

    console.log('Collections to create: ', this.selectedMarkers)
    console.log('images: ' ,this.selectedMarkers[0].collectionImages.has('files'))

    this.eventService.createCollectionsOrganized(collectionsToCreate).subscribe(res => {
      res.map((c, i) => {
        console.log('sending fies: ',i)
        this.fileuploadService.uploadCollectionImages(this.selectedMarkers[i].collectionImages, c.Id).subscribe(() => {},);
      })
    })
    this.router.navigate(['events/details', this.eventId]);
  }

  mapFromColMarkersToColections(selected: MarkerCollectionModel[]): CollectionModel[] {
    const toCreate = selected.map(s => <CollectionModel>{
      Weight: s.collectionWeight,
      CleanedTrash: s.collectionCleanedTrash,
      TrashId: s.trashId,
      EventId: this.eventId,
    })
    return toCreate
  }

  getImages() {

  }

}
