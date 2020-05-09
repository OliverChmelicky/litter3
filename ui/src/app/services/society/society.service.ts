import { Injectable } from '@angular/core';
import {HttpClient, HttpErrorResponse} from "@angular/common/http";
import {catchError} from "rxjs/operators";
import {Observable, of, throwError} from "rxjs";

import {MemberModel, SocietyModel, SocietyWithPagingAnsw} from "../../models/society.model";
import {ApisModel} from "../../api/api-urls";
import {PagingModel} from "../../models/shared.models";

@Injectable({
  providedIn: 'root'
})
export class SocietyService {
  apiUrl: string
  societyUrl: string

  constructor(
    private readonly http: HttpClient,
  ) {
    this.apiUrl = ApisModel.apiUrl;
    this.societyUrl = ApisModel.society;
  }

  createSociety(society: SocietyModel) {
    const url = `${this.apiUrl}/${this.societyUrl}/new`;
    return this.http.post<SocietyModel>(url, society).pipe(
      catchError(err => SocietyService.handleError<SocietyModel>(err))
    );
  }

  getMySocietiesIds(): Observable<MemberModel[]> {
    const url = `${this.apiUrl}/${ApisModel.user}/societies`;
    return this.http.get<MemberModel[]>(url).pipe(
      catchError(err => SocietyService.handleError<MemberModel[]>(err, [])
      )
    );
  }

  getSociety(id: string): Observable<SocietyModel> {
    const url = `${this.apiUrl}/${this.societyUrl}/${id}`;
    return this.http.get<SocietyModel>(url).pipe(
      catchError(err => SocietyService.handleError<SocietyModel>(err))
    );
  }

  getSocieties(pagingRequest: PagingModel): Observable<SocietyWithPagingAnsw> {
    const url = `${this.apiUrl}/${this.societyUrl}?from=${pagingRequest.From}&to=${pagingRequest.To}`;
    return this.http.get<SocietyWithPagingAnsw>(url).pipe(
      catchError(err => SocietyService.handleError<SocietyWithPagingAnsw>(err,{
        Societies: [],
        Paging: {
          TotalCount: 0,
          From: 0,
          To: 10,
        }
      }))
    );
  }

  getSocietiesByIds(ids: string[]): Observable<SocietyModel[]> {
    const url = `${this.apiUrl}/${this.societyUrl}/${ids}`;
    return this.http.get<SocietyModel[]>(url).pipe(
      catchError(err => SocietyService.handleError<SocietyModel[]>(err, []))
    );
  }

  leaveSociety(societyId, userId: string) {
    const url = `${this.apiUrl}/${ApisModel.user}/${this.societyUrl}/${societyId}/${userId}`;
    return this.http.delete(url).pipe(
      catchError(err => SocietyService.handleError(err))
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
