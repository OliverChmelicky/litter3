import {Injectable} from '@angular/core';
import {Observable, of, throwError} from "rxjs";
import {catchError, tap} from "rxjs/operators";
import {HttpClient, HttpErrorResponse, HttpHeaders} from '@angular/common/http';
import {ApisModel} from "../../api/api-urls";
import {UserModel} from "../../models/user.model";
import {TrashModel} from "../../models/trash.model";
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

}
