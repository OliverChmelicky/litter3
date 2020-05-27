import {Injectable} from '@angular/core';
import {HttpClient, HttpErrorResponse} from "@angular/common/http";
import {ApisModel} from "../../api/api-urls";
import {Observable, of, throwError} from "rxjs";
import {catchError} from "rxjs/operators";
import {TrashModel} from "../../models/trash.model";
import {EventPickerModel} from "../../models/event.model";

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

  uploadUserImage(event) {
    const fd = new FormData();
    fd.append("file", event.target.files[0], event.target.files[0].name);
    const url = `${ApisModel.apiUrl}/${ApisModel.fileupload}/${ApisModel.trash}`;
    return this.http.post(url, fd).pipe(
      catchError(err => FileuploadService.handleError(err))
    );
  }

  uploadSocietyImage(fd: FormData, societyId: string): Observable<any> {
    const url = `${ApisModel.apiUrl}/${ApisModel.fileupload}/${ApisModel.society}/${societyId}`;
    return this.http.post(url, fd).pipe(
      catchError(err => FileuploadService.handleError(err))
    );
  }

  uploadTrashImages(fd: FormData, trashId: string): Observable<any> {
    const url = `${ApisModel.apiUrl}/${ApisModel.fileupload}/${ApisModel.trash}/${trashId}`;
    return this.http.post(url, fd).pipe(
      catchError(err => FileuploadService.handleError(err))
    );
  }

  uploadCollectionImages(fd: FormData, collectionId: string): Observable<any> {
    const url = `${ApisModel.apiUrl}/${ApisModel.fileupload}/${ApisModel.collection}/${collectionId}`;
    return this.http.post(url, fd).pipe(
      catchError(err => FileuploadService.handleError(err))
    );
  }

  deleteUserImage() {

  }

  deleteSocietyImage() {

  }

  deleteTrashImages() {

  }

  deleteCollectionImagesFromEvent(images: string[], eventId: string, eventPickerModel: EventPickerModel) {
    const imagesQueryParam = images.join();
    const url = `${this.apiUrl}/${ApisModel.fileupload}/${ApisModel.collection}/delete?ids=${imagesQueryParam}event=${eventId}&picker=${eventPickerModel.Id}&asSociety=${eventPickerModel.AsSociety}`;
    return this.http.delete(url).pipe(
      catchError(err => FileuploadService.handleError(err))
    );
  }

  deleteCollectionImagesFromRandom(images: string[], collectionId: string) {
    const imagesQueryParam = images.join();
    const url = `${this.apiUrl}/${ApisModel.fileupload}/${ApisModel.collection}/delete/${collectionId}?ids=${imagesQueryParam}`;
    return this.http.delete(url).pipe(
      catchError(err => FileuploadService.handleError(err))
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
