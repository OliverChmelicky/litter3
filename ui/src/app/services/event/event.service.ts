import {Injectable} from '@angular/core';
import {HttpClient, HttpErrorResponse} from "@angular/common/http";
import {ApisModel} from "../../api/api-urls";
import {Observable, of, throwError} from "rxjs";
import {catchError} from "rxjs/operators";
import {
  AttendanceRequestModel, ChangePermisssionRequest,
  EventModel,
  EventPickerModel,
  EventRequestModel, EventWithCollectionsModel,
  EventWithPagingAnsw
} from "../../models/event.model";
import {SocietyWithPagingAnsw} from "../../models/society.model";
import {AttendantsModel, PagingModel} from "../../models/shared.models";

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

  setEventEditor(e: EventPickerModel) {
    localStorage.setItem('eventEditor', JSON.stringify(e));
  }

  getEventEditor(): EventPickerModel {
    const eventEditor = localStorage.getItem('eventEditor');
    return <EventPickerModel> JSON.parse(eventEditor);
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
    return this.http.get<EventWithCollectionsModel>(url).pipe(
      catchError(err => EventService.handleError<EventWithCollectionsModel>(err))
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

  notAttendEvent(eventId: string, eventPickerModel: EventPickerModel): Observable<any> {
    const url = `${this.apiUrl}/${ApisModel.event}/not-attend?event=${eventId}&picker=${eventPickerModel.Id}&asSociety=${eventPickerModel.AsSociety}`;
    return this.http.delete(url).pipe(
      catchError(err => EventService.handleError(err))
    );
  }

  createEvent(request: EventRequestModel) {
    console.log('date before: ', request.Date)
    request.Date.setSeconds(0)
    console.log('date after: ', request.Date)
    const url = `${this.apiUrl}/${ApisModel.event}`;
    return this.http.post<EventRequestModel>(url, request).pipe(
      catchError(err => EventService.handleError<EventRequestModel>(err))
    );
  }

  updateUserPermission(attendant: AttendantsModel, editor: EventPickerModel, eventId: string) {
    const request: ChangePermisssionRequest = {
      ChangingRightsTo: attendant.id,
      EventId: eventId,
      Permission: attendant.role,
      AsSociety: editor.AsSociety,  //userId is extracted from token
      SocietyId: editor.Id,
      ChangingToSociety: attendant.isSociety,
    }
    const url = `${this.apiUrl}/${ApisModel.event}/members/update`;
    return this.http.put<ChangePermisssionRequest>(url, request).pipe(
      catchError(err => EventService.handleError<ChangePermisssionRequest>(err))
    );
  }

  deleteEvent(eventEditor: EventPickerModel, eventId: string) {
    const url = `${this.apiUrl}/${ApisModel.event}/delete?event=${eventId}&picker=${eventEditor.Id}&asSociety=${eventEditor.AsSociety}`;
    return this.http.delete(url).pipe(
      catchError(err => EventService.handleError(err))
    );
  }

  updateEvent(request: EventRequestModel) {
    let newDate = new Date(request.Date)
    newDate.setSeconds(0)
    request.Date = newDate
    const url = `${this.apiUrl}/${ApisModel.event}/update`;
    return this.http.put(url, request).pipe(
      catchError(err => EventService.handleError(err))
    );
  }
}
