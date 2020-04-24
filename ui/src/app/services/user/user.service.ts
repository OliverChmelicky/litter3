import {Injectable} from '@angular/core';
import {HttpClient, HttpErrorResponse} from "@angular/common/http";

import {UserModel} from "../../models/user.model";
import {catchError} from "rxjs/operators";
import {Observable, throwError} from "rxjs";
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
    const user = JSON.parse(localStorage.getItem('user'));
    const url = `${this.apiUrl}/${ApisModel.user}/me`;
    return this.http.get<UserModel>(url).pipe(
      catchError(err => UserService.handleError<UserModel>(err))
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
        `Backend returned code ${error.status}, ` +
        `body was: ${error.error.message}`);
    }
    // return an observable with a user-facing error message
    return throwError(
      'Something bad happened; please try again later.');

    // Let the app keep running by returning an empty result.
    // return (error: any): Observable<T> => {
    //   console.error(error); // log to console instead
    //   return of(result as T);
    // };
  };
}
