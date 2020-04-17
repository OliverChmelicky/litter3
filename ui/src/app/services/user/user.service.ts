import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";

import {UserModel} from "../../models/user.model";
import {catchError, tap} from "rxjs/operators";
import {Observable, of} from "rxjs";
import {ApisModel} from "../../api/api-urls";

@Injectable({
  providedIn: 'root'
})
export class UserService {
  apiUrl: string

  activeUser: UserModel;


  constructor(
    private readonly http: HttpClient,
  ) {
    this.apiUrl = ApisModel.apiUrl
  }

  getUser(id: string): Observable<UserModel> {
    const url = `${this.apiUrl}/${ApisModel.user}/${id}`;
    return this.http.get<UserModel>(url).pipe(
      catchError(this.handleError<UserModel>())
    );
  }

  private handleError<T>( result?: T) {
    return (error: any): Observable<T> => {
      console.error(error); // log to console instead
      // Let the app keep running by returning an empty result.
      return of(result as T);
    };
  }

  getRegistered() {
    return JSON.parse(localStorage.getItem('user'));
  }
}
