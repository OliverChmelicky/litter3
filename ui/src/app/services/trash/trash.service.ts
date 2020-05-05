import {Injectable} from '@angular/core';
import {Observable} from "rxjs";
import {catchError, tap} from "rxjs/operators";
import {HttpClient, HttpHeaders} from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class TrashService {
  heroesUrl = "localhost:1323";

  constructor(
    private http: HttpClient,
  ) {
  }

  getTest(): Observable<string> {
    return this.http.get<string>(this.heroesUrl)
      .pipe(
        tap(_ => this.log('fetched heroes'))
      );
  }

  // getViaPromise(url: string): Promise<any> {
  //   return this.http.get(this.heroesUrl).toPromise().then(data => doSomethingWithData(data)).catch(err => console.error(err))
  // }
  //
  // doSomethingWithData()

  private log(fetchedHeroes: string) {
    console.log(fetchedHeroes)
  }

}
