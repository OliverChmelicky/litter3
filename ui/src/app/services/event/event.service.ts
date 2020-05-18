import {Injectable} from '@angular/core';
import {HttpClient, HttpErrorResponse} from "@angular/common/http";
import {ApisModel} from "../../api/api-urls";
import {Observable, of, throwError} from "rxjs";
import {catchError} from "rxjs/operators";
import {AttendanceRequestModel, EventModel, EventPickerModel, EventWithPagingAnsw} from "../../models/event.model";
import {SocietyWithPagingAnsw} from "../../models/society.model";
import {PagingModel} from "../../models/shared.models";

@Injectable({
  providedIn: 'root'
})
export class EventService {
  apiUrl: string;

  constructor(
    private http: HttpClient,
  ) {
    this.apiUrl = ApisModel.apiUrl
  }

  getSocietyEvents(societyId: string): Observable<EventModel[]> {
    const url = `${this.apiUrl}/${ApisModel.event}/${ApisModel.society}/${societyId}`;
    return this.http.get<EventModel[]>(url).pipe(
      catchError(err => EventService.handleError<EventModel[]>(err, []))
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

  getEvent(eventId: string) {
    const url = `${this.apiUrl}/${ApisModel.event}/${eventId}`;
    return this.http.get<EventModel>(url).pipe(
      catchError(err => EventService.handleError<EventModel>(err))
    );
  }

  getEvents(pagingRequest: PagingModel): Observable<EventWithPagingAnsw> {
    const url = `${this.apiUrl}/${ApisModel.event}?from=${pagingRequest.From}&to=${pagingRequest.To}`;
    return this.http.get<EventWithPagingAnsw>(url).pipe(
      catchError(err => EventService.handleError<EventWithPagingAnsw>(err, {
        Events: [],
        Paging: {
          TotalCount: 0,
          From: 0,
          To: 10,
        }
      }))
    );
  }

  attendEvent(eventId: string, eventPickerModel: EventPickerModel) {
    const url = `${this.apiUrl}/${ApisModel.event}/attend`;
    const request: AttendanceRequestModel = {
      PickerId: eventPickerModel.Id,
      EventId: eventId,
      AsSociety: eventPickerModel.AsSociety,
    }
    return this.http.post<AttendanceRequestModel>(url, request).pipe(
      catchError(err => EventService.handleError<AttendanceRequestModel>(err))
    );
  }

  notAttendEvent(eventId: string, eventPickerModel: EventPickerModel) {
    const url = `${this.apiUrl}/${ApisModel.event}/not-attend?event=${eventId}&picker=${eventPickerModel.Id}&asSociety=${eventPickerModel.AsSociety}`;
    return this.http.delete<AttendanceRequestModel>(url).pipe(
      catchError(err => EventService.handleError<AttendanceRequestModel>(err))
    );
  }
}
