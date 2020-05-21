import { Injectable } from '@angular/core';
import {HttpClient, HttpErrorResponse} from "@angular/common/http";
import {catchError} from "rxjs/operators";
import {Observable, of, throwError} from "rxjs";

import {
  ApplicantModel,
  MemberModel,
  SocietyModel,
  SocietyWithPagingAnsw,
  UserSocietyRequestModel
} from "../../models/society.model";
import {ApisModel} from "../../api/api-urls";
import {PagingModel} from "../../models/shared.models";
import {UserGroupModel, UserModel} from "../../models/user.model";

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

  getSocietyMembers(societyId: string): Observable<MemberModel[]> {
    const url = `${this.apiUrl}/${ApisModel.society}/members/${societyId}`;
    return this.http.get<MemberModel[]>(url).pipe(
      catchError(err => SocietyService.handleError<MemberModel[]>(err,[]))
    );
  }

  askForMembership(societyId: string) {
    const request = {
      Id: societyId
    }
    const url = `${this.apiUrl}/membership`;
    return this.http.post(url, request).pipe(
      catchError(err => SocietyService.handleError(err))
    );
  }

  leaveSociety(societyId, userId: string) {
    const url = `${this.apiUrl}/${this.societyUrl}/${societyId}/${userId}`;
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


  updateSociety(society: SocietyModel): Observable<SocietyModel> {
    const url = `${this.apiUrl}/${ApisModel.society}/update`;
    return this.http.put<SocietyModel>(url, society).pipe(
      catchError(err => SocietyService.handleError<SocietyModel>(err))
    );
  }

  changePermissions(changeMemberPermission: MemberModel[]) {
      const url = `${this.apiUrl}/${ApisModel.society}/change/permission`;
      return this.http.put<SocietyModel>(url, changeMemberPermission).pipe(
        catchError(err => SocietyService.handleError<SocietyModel>(err))
      );
  }

  removeUser(userId, societyId: string) {
    const url = `${this.apiUrl}/${ApisModel.society}/${societyId}/${userId}`;
    return this.http.delete(url).pipe(
      catchError(err => SocietyService.handleError(err))
    );
  }

  deleteSociety(societyId: string) {
    const url = `${this.apiUrl}/${ApisModel.society}/delete/${societyId}`;
    return this.http.delete(url).pipe(
      catchError(err => SocietyService.handleError(err))
    );
  }

  removeApplication(societyId: string) {
    const url = `${this.apiUrl}/membership/${societyId}`;
    return this.http.delete(url).pipe(
      catchError(err => SocietyService.handleError(err))
    );
  }

  acceptApplicant(societyId: string, userId: string) {
    const request: UserGroupModel = {
      UserId: userId,
      SocietyId: societyId
    }
    const url = `${this.apiUrl}/membership/accept/${societyId}/${userId}`;
    return this.http.post(url, request).pipe(
      catchError(err => SocietyService.handleError(err))
    );
  }

  dismissApplicant(societyId, userId: string) {
    const url = `${this.apiUrl}/membership/deny/${societyId}/${userId}`;
    return this.http.delete(url).pipe(
      catchError(err => SocietyService.handleError(err))
    );
  }

}
