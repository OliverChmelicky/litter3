import {Injectable} from '@angular/core';
import {HttpClient, HttpErrorResponse} from "@angular/common/http";
import {ApisModel} from "../../api/api-urls";
import {Observable, of, throwError} from "rxjs";
import {catchError} from "rxjs/operators";
import {TrashModel} from "../../models/trash.model";

@Injectable({
  providedIn: 'root'
})
export class FileuploadService {
  apiUrl: string;

  constructor(
    private http: HttpClient,
  ) {
    this.apiUrl = ApisModel.apiUrl
  }

  getUserImage() {

  }

  getSocietyImage() {

  }

  getTrashImages() {

  }

  getCollectionImages() {

  }

  uploadUserImage(event) {
    const fd = new FormData();
    fd.append("file", event.target.files[0], event.target.files[0].name);
    const url = `${ApisModel.apiUrl}/${ApisModel.fileupload}/${ApisModel.trash}`;
    return this.http.post(url, fd).pipe(
      catchError(err => FileuploadService.handleError(err))
    );
  }

  uploadSocietyImage() {

  }

  uploadTrashImages(fd: FormData, trashId: string): Observable<any> {
    const url = `${ApisModel.apiUrl}/${ApisModel.fileupload}/${ApisModel.trash}/${trashId}`;
    return this.http.post(url, fd).pipe(
      catchError(err => FileuploadService.handleError(err))
    );
  }

  uploadCollectionImages() {

  }

  deleteUserImage() {

  }

  deleteSocietyImage() {

  }

  deleteTrashImages() {

  }

  deleteCollectionImages() {

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