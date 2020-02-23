import { Injectable } from '@angular/core';
import {Observable} from "rxjs";
import {catchError, tap} from "rxjs/operators";
import {HttpClient, HttpHeaders} from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class TrashService {
  heroesUrl = "localhost:";

  constructor(
    private http: HttpClient,
  ) { }

  getTest(): Observable<string> {
    return this.http.get<string>(this.heroesUrl)
      .pipe(
        tap(_ => this.log('fetched heroes')),
        catchError(this.handleError<string[]>('getHeroes', []))
      );
  }

  private log(fetchedHeroes: string) {
    console.log(fetchedHeroes)
  }


  private handleError<T>(heroes: string, anies: any[]) {
    return function (p1: any, p2: Observable<A>) {
      return undefined;
    };
  }
}


private handleError<T>(operation = 'operation', result?: T) {
  return (error: any): Observable<T> => {

    // TODO: send the error to remote logging infrastructure
    console.error(error); // log to console instead

    // TODO: better job of transforming error for user consumption
    this.log(`${operation} failed: ${error.message}`);

    // Let the app keep running by returning an empty result.
    return of(result as T);
  };
}
