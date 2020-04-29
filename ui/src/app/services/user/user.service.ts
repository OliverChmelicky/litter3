import {Injectable} from '@angular/core';
import {HttpClient, HttpErrorResponse} from "@angular/common/http";

import {
  UserModel,
  FriendRequestModel,
  EmailMessageModel,
  MemberModel,
  SocietyModel,
  IdsMessageModel
} from "../../models/user.model";
import {catchError} from "rxjs/operators";
import {Observable, of, throwError} from "rxjs";
import {ApisModel} from "../../api/api-urls";
import * as firebase from "firebase";

@Injectable({
  providedIn: 'root'
})
export class UserService {
  apiUrl: string

  constructor(
    private readonly http: HttpClient,
  ) {
    this.apiUrl = ApisModel.apiUrl
  }

  createUser(userDetails: UserModel) {
    console.log('vytvaram usr')
    const url = `${this.apiUrl}/${ApisModel.user}/new`;
    return this.http.post<UserModel>(url, userDetails).pipe(
      catchError(err => UserService.handleError<UserModel>(err))
    );
  }

  getUser(id: string): Observable<UserModel> {
    const url = `${this.apiUrl}/${ApisModel.user}/${id}`;
    return this.http.get<UserModel>(url).pipe(
      catchError(err => UserService.handleError<UserModel>(err))
    );
  }

  getUserByEmail(email: string): Observable<UserModel> {
    const url = `${this.apiUrl}/${ApisModel.user}/${email}`;
    return this.http.get<UserModel>(url).pipe(
      catchError(err => UserService.handleError<UserModel>(err))
    );
  }

  getMe(): Observable<UserModel> {
    const url = `${this.apiUrl}/${ApisModel.user}/me`;
    return this.http.get<UserModel>(url).pipe(
      catchError(err => UserService.handleError<UserModel>(err)
      )
    );
  }

  getMyFriendRequests(): Observable<FriendRequestModel[]> {
    const url = `${this.apiUrl}/${ApisModel.user}/my/requests/friendship`;
    return this.http.get<FriendRequestModel[]>(url).pipe(
      catchError(err => UserService.handleError<FriendRequestModel[]>(err, [])
      )
    );
  }

  getUsersDetails(ids: string[]): Observable<UserModel[]> {
    const idsQueryParam = ids.join();
    const url = `${this.apiUrl}/${ApisModel.user}/details?Ids=${idsQueryParam}`;
    return this.http.get<UserModel[]>(url).pipe(
      catchError(err => UserService.handleError<UserModel[]>(err)
      )
    );
  }

  requestFriend(email: string): Observable<EmailMessageModel> {
    const url = `${this.apiUrl}/${ApisModel.user}/friend/add/email`;
    const request = <EmailMessageModel>{
      Email: email
    }
    return this.http.post<EmailMessageModel>(url, request).pipe(
      catchError(err => UserService.handleError<EmailMessageModel>(err))
    )
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
    if (result === null) {
      return throwError(error);
    }
    return of(result as T);
  };
}
