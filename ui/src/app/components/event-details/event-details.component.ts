import {Component, Inject, OnInit} from '@angular/core';
import {UserModel} from "../../models/user.model";
import {
  EventPickerModel,
  EventSocietyModel,
  EventUserModel,
  EventWithCollectionsModel
} from "../../models/event.model";
import {attendantsTableColumns} from "./table-definitions";
import {AttendantsModel} from "../../models/shared.models";
import {ActivatedRoute, Router} from "@angular/router";
import {UserService} from "../../services/user/user.service";
import {EventService} from "../../services/event/event.service";
import {SocietyService} from "../../services/society/society.service";
import {SocietyModel} from "../../models/society.model";
import {MatTableDataSource} from "@angular/material/table";
import {GoogleMap} from "@agm/core/services/google-maps-types";
import {MarkerModel} from "../google-map/Marker.model";
import {MapLocationModel} from "../../models/GPSlocation.model";
import {
  CollectionModel,
  defaultCollectionModel,
  defaultTrashImage
} from "../../models/trash.model";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material/dialog";
import {TrashService} from "../../services/trash/trash.service";
import {FileuploadService} from "../../services/fileupload/fileupload.service";

export const czechPosition: MapLocationModel = {
  lat: 49.81500022397678,
  lng: 20.0,
  zoom: 7,
  minZoom: 3,
};

export interface DialogDataEditCollection {
  collection: CollectionModel;
  deleteImages: string[];
  uploadImages: FormData;
  deleteCollection: boolean
}

export interface DialogDataShowCollection {
  collection: CollectionModel;
  deleteImages: string[];
  uploadImages: FormData;
  deleteCollection: boolean
}

@Component({
  selector: 'app-event-details',
  templateUrl: './event-details.component.html',
  styleUrls: ['./event-details.component.css']
})
export class EventDetailsComponent implements OnInit {
  isLoggedIn: boolean = false;
  statusAttend: boolean = false;
  me: UserModel = {
    Id: '',
    FirstName: '',
    LastName: '',
    Email: '',
    Uid: '',
    Avatar: '',
    CreatedAt: new Date(),
  };
  map: GoogleMap;
  event: EventWithCollectionsModel = {
    Date: new Date,
    Description: '',
  };
  attendants: AttendantsModel[] = [];
  tableColumns = attendantsTableColumns;
  availableDecisionsAs: EventPickerModel[] = [];
  selectedCreator: number = 0;

  isAdmin: boolean = false;
  isEditor: boolean = false;
  editableSocieties: SocietyModel[] = [];
  editableSocietiesIds: string[] = [];

  markers: MarkerModel[] = [];
  initLat: number = czechPosition.lat;
  initLng: number = czechPosition.lng;

  societies: SocietyModel[] = [];
  users: UserModel[] = [];
  showCollectionsTable: CollectionModel[] = [];


  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private userService: UserService,
    private societyService: SocietyService,
    private eventService: EventService,
    private trashService: TrashService,
    private fileuploadService: FileuploadService,
    public showCollectionDialog: MatDialog,
    public editCollectionDialog: MatDialog,
  ) {
  }

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      this.eventService.getEvent(params.get('eventId')).subscribe(event => {
          this.convertToLocalTime()
          this.event = event
          if (event.Collections) {
            this.mapCollectionsToTable(event.Collections)
          }
          if (event.Trash) {
            this.initLat = event.Trash[0].Location[0]
            this.initLng = event.Trash[0].Location[1]
          }
          if (this.event.Trash) {
            this.assignMarkers()
          }

          if (event.UsersIds) {
            const userIds = event.UsersIds.map(u => u.UserId)
            this.fetchUsersWhoAttends(userIds, event.UsersIds)
          }

          if (event.SocietiesIds) {
            const societyIds = event.SocietiesIds.map(s => s.SocietyId)
            this.fetchSocietiesWhichAttend(societyIds, event.SocietiesIds)
          }
        },
        () => {
        },
        () => {
          this.userService.getMe().subscribe(
            me => {
              this.me = me;
              this.isLoggedIn = true;
              this.availableDecisionsAs.push({
                VisibleName: me.Email,
                Id: me.Id,
                AsSociety: false
              })
              this.getEventAttendanceOnUser()
              this.userService.getMyEditableSocieties().subscribe(
                societies => {
                  if (societies) {
                    this.editableSocieties = societies
                    this.editableSocietiesIds = societies.map(soc => {
                      this.availableDecisionsAs.push({
                        VisibleName: soc.Name,
                        Id: soc.Id,
                        AsSociety: true
                      })
                      return soc.Id
                    })
                  }
                }
              )
              //find out my status
              if (this.event.UsersIds) {
                this.event.UsersIds.map(attendant => {
                  if (attendant.Permission === 'editor' && attendant.UserId === this.me.Id) {
                    this.isEditor = true
                  }
                  if (attendant.Permission === 'creator' && attendant.UserId === this.me.Id) {
                    this.isAdmin = true
                  }
                })
              }
            },
            () => console.log('You are not registered')
          )
        }
      )

    })
  }

  onMapReady(map: GoogleMap) {
    this.map = map;
  }

  onDesideAsChange() {
    this.statusAttend = false
    this.isAdmin = false
    this.isEditor = false
    if (this.selectedCreator === 0) {  //User
      this.getEventAttendanceOnUser()
    } else {
      this.getSocietyEventAttendance()
    }
  }

  private getEventAttendanceOnUser() {
    if (this.event.UsersIds) {
      this.event.UsersIds.map(
        user => {
          if (user.UserId == this.me.Id) {
            this.statusAttend = true
            if (user.Permission === 'editor') {
              this.isEditor = true
            } else if (user.Permission === 'creator') {
              this.isAdmin = true
              this.isEditor = true
            }
          }
        }
      )
    }
  }

  private getSocietyEventAttendance() {
    this.attendants.map(a => {
      if (a.isSociety && this.availableDecisionsAs[this.selectedCreator].Id === a.id) {
        this.statusAttend = true
        if (a.role === 'editor') {
          this.isEditor = true
        }
        if (a.role === 'creator') {
          this.isAdmin = true
        }
      }
    })
  }

  onAttend() {
    this.eventService.attendEvent(this.event.Id, this.availableDecisionsAs[this.selectedCreator]).subscribe(
      () => {
        this.statusAttend = true;
        this.pushAttendant(this.availableDecisionsAs[this.selectedCreator])

        const newData = new MatTableDataSource<AttendantsModel>(this.attendants);
        this.attendants = []
        for (let i = 0; i < newData.data.length; i++) {
          this.attendants.push(newData.data[i])
        }
      }
    )
  }

  onNotAttend() {
    this.eventService.notAttendEvent(this.event.Id, this.availableDecisionsAs[this.selectedCreator]).subscribe(
      () => {
        this.statusAttend = false
        this.isEditor = false
        this.isAdmin = false

        const index = this.attendants.findIndex(a => a.id === this.availableDecisionsAs[this.selectedCreator].Id)
        this.attendants.splice(index, 1)

        const newData = new MatTableDataSource<AttendantsModel>(this.attendants);
        this.attendants = []
        for (let i = 0; i < newData.data.length; i++) {
          this.attendants.push(newData.data[i])
        }
      }
    )
  }

  onEdit() {
    this.eventService.setEventEditor(this.availableDecisionsAs[this.selectedCreator])
    this.router.navigate(['events/edit', this.event.Id])
  }

  fetchUsersWhoAttends(userIds: string[], userEventDetails: EventUserModel[]) {
    if (userIds.length > 1) {
      this.userService.getUsersDetails(userIds).subscribe(
        users => {
          this.users = users
          users.map(u => {
            userEventDetails.map(detail => {
              if (u.Id === detail.UserId) {
                this.attendants.push({
                  id: u.Id,
                  name: u.Email,
                  avatar: u.Avatar ? u.Avatar : '',
                  role: detail.Permission,
                  isSociety: false,
                })
              }
              const newData = new MatTableDataSource<AttendantsModel>(this.attendants);
              this.attendants = []
              for (let i = 0; i < newData.data.length; i++) {
                this.attendants.push(newData.data[i])
              }
            });
          })
        })
    } else {
      this.userService.getUser(userIds[0]).subscribe(
        user => {
          this.users.push(user)
          userEventDetails.map(detail => {
            if (user.Id === detail.UserId) {
              this.attendants.push({
                id: user.Id,
                name: user.Email,
                avatar: user.Avatar ? user.Avatar : '',
                role: detail.Permission,
                isSociety: false
              })
            }
            const newData = new MatTableDataSource<AttendantsModel>(this.attendants);
            this.attendants = []
            for (let i = 0; i < newData.data.length; i++) {
              this.attendants.push(newData.data[i])
            }
          });

        })
    }
  }

  fetchSocietiesWhichAttend(societiesIds: string[], societyEventDetails: EventSocietyModel[]) {
    if (societiesIds.length > 1) {
      this.societyService.getSocietiesByIds(societiesIds).subscribe(
        societies => {
          this.societies = societies
          societies.map(s => {
            societyEventDetails.map(detail => {
              if (s.Id === detail.SocietyId) {
                this.attendants.push({
                  id: s.Id,
                  name: s.Name,
                  avatar: s.Avatar ? s.Avatar : '',
                  role: detail.Permission,
                  isSociety: true
                })
              }
            })
          })
          const newData = new MatTableDataSource<AttendantsModel>(this.attendants);
          this.attendants = []
          for (let i = 0; i < newData.data.length; i++) {
            this.attendants.push(newData.data[i])
          }
        })
    } else {
      this.societyService.getSociety(societiesIds[0]).subscribe(
        society => {
          this.societies.push(society)
          societyEventDetails.map(detail => {
            if (society.Id === detail.SocietyId) {
              this.attendants.push({
                id: society.Id,
                name: society.Name,
                avatar: society.Avatar ? society.Avatar : '',
                role: detail.Permission,
                isSociety: true,
              })
            }
          })
          const newData = new MatTableDataSource<AttendantsModel>(this.attendants);
          this.attendants = []
          for (let i = 0; i < newData.data.length; i++) {
            this.attendants.push(newData.data[i])
          }
        })
    }
  }

  private convertToLocalTime() {
    let dateStr = this.event.Date.toString()
    dateStr += ' UTC'
    this.event.Date = new Date(dateStr)
  }

  navigateToTrash(trashId: string) {
    this.router.navigate(['trash/details', trashId])
  }

  private assignMarkers() {
    const trashIds = this.event.Trash.map( t => t.Id)
    this.trashService.getTrashByIds(trashIds).subscribe( trash => {
      trash.map(t => {
        let collLength = 0
        if (t.Collections) {
          collLength = t.Collections.length
        }
        if (!t.Images) {
          t.Images = [defaultTrashImage];
        }
        this.markers.push({
          id: t.Id,
          lat: t.Location[0],
          lng: t.Location[1],
          cleaned: t.Cleaned,
          images: t.Images,
          numOfCollections: collLength
        })
      })
    })
  }

  private pushAttendant(picker: EventPickerModel) {
    if (picker.AsSociety) {
      this.editableSocieties.map(s => {
        if (s.Id === picker.Id) {
          this.attendants.push({
            id: s.Id,
            name: picker.VisibleName,
            avatar: s.Avatar ? s.Avatar : '',
            role: 'viewer',
            isSociety: true,
          })
        }
      })
    } else {
      this.attendants.push({
        id: this.me.Id,
        name: this.me.Email,
        avatar: this.me.Avatar ? this.me.Avatar : '',
        role: 'viewer',
        isSociety: false
      })

    }


  }

  onCreateCollections() {
    let trashIds = []
    if (this.event.Trash) {
      trashIds = this.event.Trash.map(t => t.Id)
    }
    this.eventService.setEventEditor(this.availableDecisionsAs[this.selectedCreator])
    this.router.navigate(['collection'], {queryParams: {trashIds: trashIds, 'eventId': this.event.Id}})
  }

  onEditCollection(collectionId: string) {
    let collection: CollectionModel = defaultCollectionModel
    this.event.Collections.map(c => {
      if (c.Id === collectionId) {
        collection = c
      }
    })

    if (!collection.Images) {
      collection.Images = [];
    }

    const weightBefore = collection.Weight
    const cleanedBefore = collection.CleanedTrash

    const dialogRef = this.editCollectionDialog.open(EditCollectionComponent, {
      width: '800px',
      data: {
        collection: collection,
        deleteImages: [],
        uploadImages: new FormData(),
        deleteCollection: false
      }
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        if (result.deleteCollection) {
          this.eventService.deleteCollectionOrganized(this.event.Id, result.collection.Id, this.availableDecisionsAs[this.selectedCreator]).subscribe(
            () => window.location.reload(),
          )
          return
        }
        if (result.deleteImages.length > 0) {
          this.trashService.deleteCollectionEventImages(collectionId, result.deleteImages, this.event.Id, this.availableDecisionsAs[this.selectedCreator]).subscribe(
            () => {
            },
            error => console.log(error)
          )
        }
        if (result.uploadImages) {
          if (result.uploadImages.has('files')) {
            this.fileuploadService.uploadCollectionImages(result.uploadImages, collectionId).subscribe()
          }
        }
        if (result.collection.Weight !== weightBefore || result.collection.CleanedTrash !== cleanedBefore) {
          this.eventService.updateCollectionOrganized(result.collection, this.event.Id ,this.availableDecisionsAs[this.selectedCreator]).subscribe(
            () => window.location.reload(),
            error => {console.log(error)}
            )
        }
      }
    });
  }

  onShowCollection(collectionId: string) {
    let collection: CollectionModel = defaultCollectionModel
    this.event.Collections.map(c => {
      if (c.Id === collectionId) {
        collection = c
      }
    })

    if (!collection.Images) {
      collection.Images = [];
    }

    const dialogRef = this.showCollectionDialog.open(ShowCollectionComponent, {
      width: '800px',
      data: {
        collection: collection,
      }
    });

    dialogRef.afterClosed().subscribe(() => {});
  }

  private mapCollectionsToTable(collections: CollectionModel[]) {
    collections.map( c => {
      let images;
      if (!c.Images) {
        images = [defaultTrashImage]
      } else {
        images = c.Images
      }
      this.showCollectionsTable.push(
        {
          Id: c.Id,
          Weight: c.Weight,
          CleanedTrash: c.CleanedTrash,
          TrashId: c.TrashId,
          EventId: c.EventId,
          Images: images,
          Users: [],
          CreatedAt: c.CreatedAt,
        }
      )
    })
  }
}


//DIALOG INFO


@Component({
  selector: 'app-edit-collection',
  templateUrl: './dialog-edit-collection/edit-collection-dialog.component.html',
  styleUrls: ['./dialog-edit-collection/edit-collection-dialog.component.css'],
})
export class EditCollectionComponent {
  images = []
  show: boolean = true;

  constructor(public dialogRef: MatDialogRef<EditCollectionComponent>,
              @Inject(MAT_DIALOG_DATA) public data: DialogDataEditCollection) {
    data.collection.Images.map(i =>
      this.images.push({
      Url: i.Url,
      inDeleteList: false,
    }))
  }

  onNoClick(): void {
    this.dialogRef.close();
  }

  // onDelete(url: string) {
  //   this.data.deleteImages.push(url)
  //   const index = this.images.findIndex(i => i.Url === url)
  //   this.images[index].inDeleteList = true
  //   this.reload()
  // }
  //
  // onRemoveFromList(url: string) {
  //   let index = this.data.deleteImages.findIndex(i => i === url)
  //   this.data.deleteImages = this.data.deleteImages.splice(index, 1)
  //
  //   index = this.images.findIndex(i => i.Url === url)
  //   this.images[index].inDeleteList = false
  //   this.reload()
  // }

  // reload() {
  //   this.show = false;
  //   setTimeout(() => this.show = true);
  // }

  onSave() {
    this.data.collection.Images.map(image => {
      if (image.inDeleteList) {
        this.data.deleteImages.push(image.Url)
      }
    })
    this.dialogRef.close(this.data);
  }

  onUpload(event) {
    this.data.uploadImages.delete('files')
    for (let i = 0; i < event.target.files.length; i++) {
      this.data.uploadImages.append("files", event.target.files[i], event.target.files[i].name);
    }
  }

  onSetDelete() {
    this.data.deleteCollection = true
    this.data.uploadImages = new FormData()
    this.dialogRef.close(this.data);
  }
}

@Component({
  selector: 'app-show-collection',
  templateUrl: './dialog-collection-detail/detail-dialog.component.html',
  styleUrls: ['./dialog-collection-detail/detail-dialog.component.css']
})
export class ShowCollectionComponent {

  constructor( public dialogRef: MatDialogRef<ShowCollectionComponent>,
               @Inject(MAT_DIALOG_DATA) public data: DialogDataShowCollection) {
  }

  onNoClick(): void {
    this.dialogRef.close();
  }

}
