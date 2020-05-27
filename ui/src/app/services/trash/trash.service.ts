import {Injectable} from '@angular/core';
import {Observable, of, throwError} from "rxjs";
import {catchError, tap} from "rxjs/operators";
import {HttpClient, HttpErrorResponse, HttpHeaders} from '@angular/common/http';
import {ApisModel} from "../../api/api-urls";
import {UserModel} from "../../models/user.model";
import {
  AddPickersToCollectionRequest,
  CollectionImageModel,
  CollectionModel,
  CollectionUserModel,
  CommentModel,
  CreateCollectionRandomRequest,
  TrashModel,
  TrashTypeAutomotive,
  TrashTypeBooleanValues,
  TrashTypeCarcass,
  TrashTypeConstruction,
  TrashTypeDangerous,
  TrashTypeElectronic,
  TrashTypeGlass,
  TrashTypeHousehold,
  TrashTypeMask,
  TrashTypeMetal,
  TrashTypeOrganic,
  TrashTypeOther,
  TrashTypePlastics,
  UpdateCollectionModel
} from "../../models/trash.model";
import {MarkerModel} from "../../components/google-map/Marker.model";

@Injectable({
  providedIn: 'root'
})
export class TrashService {
  apiUrl: string;

  constructor(
    private http: HttpClient,
  ) {
    this.apiUrl = ApisModel.apiUrl
  }

  createTrash(trash: TrashModel): Observable<TrashModel> {
    const url = `${this.apiUrl}/${ApisModel.trash}/new`;
    return this.http.post<TrashModel>(url, trash).pipe(
      catchError(err => TrashService.handleError<TrashModel>(err))
    );
  }

  getTrashInRange(lat,lng: number,d: number): Observable<TrashModel[]> {
    const url = `${this.apiUrl}/${ApisModel.trash}/range?lat=${lat}&lng=${lng}&radius=${d}`;
    return this.http.get<TrashModel[]>(url).pipe(
      catchError(err => TrashService.handleError<TrashModel[]>(err))
    );
  }

  getTrashById(trashId: string): Observable<TrashModel> {
    const url = `${this.apiUrl}/${ApisModel.trash}/${trashId}`;
    return this.http.get<TrashModel>(url).pipe(
      catchError(err => TrashService.handleError<TrashModel>(err))
    );
  }

  getTrashByIds(trashIds: string[]): Observable<TrashModel[]> {
    const idsQueryParam = trashIds.join();
    const url = `${this.apiUrl}/${ApisModel.trash}?ids=${idsQueryParam}`;
    return this.http.get<TrashModel[]>(url).pipe(
      catchError(err => TrashService.handleError<TrashModel[]>(err,[]))
    );
  }


  private static handleError<T>(error: HttpErrorResponse, result?: T) {
    if (error.error instanceof ErrorEvent) {
      // A client-side or network error occurred. Handle it accordingly.
      console.error('An error occurred:', error.error.message);
    } else {
      // The backend returned an unsuccessful response code.
      // The response body may contain clues as to what went wrong,
      console.error(
        `Backend returned code ${error.status} \n` +
        `TITLE: ${error.error.errorMessage} \n` +
        `TYPE ${error.error.errorType} `);
    }
    // return an observable with a user-facing error message
    if (result == null) {
      return throwError(error);
    }
    return of(result as T);
  };

  updateTrash(trash: TrashModel) {
    const url = `${this.apiUrl}/${ApisModel.trash}/update`;
    return this.http.put<TrashModel>(url, trash).pipe(
      catchError(err => TrashService.handleError<TrashModel>(err))
    );
  }

  deleteTrash(trashId: string) {
    const url = `${this.apiUrl}/${ApisModel.trash}/delete/${trashId}`;
    return this.http.delete(url).pipe(
      catchError(err => TrashService.handleError(err))
    );
  }

  convertTrashTypeNumToBools(TrashType: number): TrashTypeBooleanValues {
    return  {
        TrashTypeHousehold: (TrashType & TrashTypeHousehold) > 0,
        TrashTypeAutomotive: (TrashType & TrashTypeAutomotive) > 0,
        TrashTypeConstruction: (TrashType & TrashTypeConstruction) > 0,
        TrashTypePlastics: (TrashType & TrashTypePlastics) > 0,
        TrashTypeElectronic: (TrashType & TrashTypeElectronic) > 0,
        TrashTypeGlass: (TrashType & TrashTypeGlass) > 0,
        TrashTypeMetal: (TrashType & TrashTypeMetal) > 0,
        TrashTypeDangerous: (TrashType & TrashTypeDangerous) > 0,
        TrashTypeCarcass: (TrashType & TrashTypeCarcass) > 0,
        TrashTypeOrganic: (TrashType & TrashTypeOrganic) > 0,
        TrashTypeOther: (TrashType & TrashTypeOther) > 0,
    }
  }

  convertTrashTypeBoolsToNums (TrashType: TrashTypeBooleanValues): number {
    const TrashTypeHouseholdMasked     = TrashType.TrashTypeHousehold ? TrashTypeHousehold : 0
    const TrashTypeAutomotiveMasked    = TrashType.TrashTypeAutomotive ? TrashTypeAutomotive : 0
    const TrashTypeConstructionMasked  = TrashType.TrashTypeConstruction ? TrashTypeConstruction : 0
    const TrashTypePlasticsMasked      = TrashType.TrashTypePlastics ? TrashTypePlastics : 0
    const TrashTypeElectronicMasked    = TrashType.TrashTypeElectronic ? TrashTypeElectronic : 0
    const TrashTypeGlassMasked         = TrashType.TrashTypeGlass ? TrashTypeGlass : 0
    const TrashTypeMetalMasked         = TrashType.TrashTypeMetal ? TrashTypeMetal : 0
    const TrashTypeDangerousMasked     = TrashType.TrashTypeDangerous ? TrashTypeDangerous : 0
    const TrashTypeCarcassMasked       = TrashType.TrashTypeCarcass ? TrashTypeCarcass : 0
    const TrashTypeOrganicMasked       = TrashType.TrashTypeOrganic ? TrashTypeOrganic : 0
    const TrashTypeOtherMasked         = TrashType.TrashTypeOther ? TrashTypeOther : 0

    return TrashTypeHouseholdMasked | TrashTypeAutomotiveMasked | TrashTypeConstructionMasked | TrashTypePlasticsMasked | TrashTypeElectronicMasked |
            TrashTypeGlassMasked | TrashTypeMetalMasked | TrashTypeDangerousMasked | TrashTypeCarcassMasked | TrashTypeOrganicMasked | TrashTypeOtherMasked

  }

  deleteTrashImage(image: string, trashId: string) {
    const url = `${this.apiUrl}/${ApisModel.fileupload}/${ApisModel.trash}/delete/${trashId}/${image}`;
    return this.http.delete(url).pipe(
      catchError(err => TrashService.handleError(err))
    );
  }

  deleteCollectionImage(image: string, collectionId: string) {
    const url = `${this.apiUrl}/${ApisModel.fileupload}/${ApisModel.collection}/delete/${collectionId}/${image}`;
    return this.http.delete(url).pipe(
      catchError(err => TrashService.handleError(err))
    );
  }

  getCollectionById(collectionId: string): Observable<CollectionModel> {
    const url = `${this.apiUrl}/${ApisModel.collection}/${collectionId}`;
    return this.http.get<CollectionModel>(url).pipe(
      catchError(err => TrashService.handleError<CollectionModel>(err))
    );
  }

  commentTrash(message: string, trashId: string): Observable<CommentModel> {
    const request = {
      message: message,
      id: trashId,
    }
    const url = `${this.apiUrl}/${ApisModel.trash}/comment`;
    return this.http.post<CommentModel>(url, request).pipe(
      catchError(err => TrashService.handleError<CommentModel>(err))
    );
  }

  //I want not cleaned and if I discover new places in map I don`t want the one in table
  filterCleanedAndSelected(markers: MarkerModel[], selectedTrash: MarkerModel[]): MarkerModel[]{
   let filteredMarkers: MarkerModel[] = []
    for (let i = 0; i < markers.length; i++) {
     let found = false
     for (let j = 0; j < selectedTrash.length; j++) {
       if (markers[i].id === selectedTrash[j].id) {
         found = true;
         break
       }
     }
     if (!found) {
       filteredMarkers.push(markers[i])
     }
   }

    filteredMarkers.filter( m => !m.cleaned)
    return filteredMarkers
  }

  deleteCollectionFromUser(collectionId: string) {
    const url = `${this.apiUrl}/${ApisModel.collection}/delete/${collectionId}`;
    return this.http.delete(url).pipe(
      catchError(err => TrashService.handleError(err))
    );
  }

  updateCollection(collection: CollectionModel) {
    const url = `${this.apiUrl}/${ApisModel.collection}/update/col-random`;
    return this.http.put(url, collection).pipe(
      catchError(err => TrashService.handleError(err))
    );
  }

  getIdsOfTrashOfUsers(): Observable<CollectionUserModel[]> {
    const url = `${this.apiUrl}/${ApisModel.collection}/personal`;
    return this.http.get<CollectionUserModel[]>(url).pipe(
      catchError(err => TrashService.handleError<CollectionUserModel[]>(err, []))
    );
  }


  createCollection(request: CreateCollectionRandomRequest): Observable<CollectionModel> {
    const url = `${this.apiUrl}/${ApisModel.collection}/random`;
    return this.http.post<CollectionModel>(url, request).pipe(
      catchError(err => TrashService.handleError<CollectionModel>(err))
    );
  }

  addFriendsToCollection(friends: string[], collectionId: string) {
    const request: AddPickersToCollectionRequest = {
      CollectionId: collectionId,
      UserId: friends,
    }
    const url = `${this.apiUrl}/${ApisModel.collection}/add-picker`;
    return this.http.post<CollectionModel>(url, request).pipe(
      catchError(err => TrashService.handleError<CollectionModel>(err))
    );
  }
}
